package service

import (
	"errors"
	"fmt"
	logz "github.com/faelmori/golife/internal/log"
	"github.com/goccy/go-json"
	"github.com/pebbe/zmq4"
	"time"
)

var (
	quitCh    = make(chan bool, 1)
	sendCh    = make(chan interface{}, 100)
	receiveCh = make(chan []string, 100)
	timeCh    = make(chan time.Time, 1)
)

type BrokerClient struct {
	endpoint    string
	client      *zmq4.Socket
	poller      *zmq4.Poller
	receiveCh   chan []string
	sendCh      chan interface{}
	timeCh      chan time.Time
	dataCh      chan interface{}
	quitCh      chan bool
	timeout     time.Duration
	tryInterval time.Duration
	retryLimit  int
	retries     int
}

const defaultBrokerEndpoint = "tcp://localhost:5555"

func NewBrokerClient(brokerEndpoint string, dataCh chan interface{}) (*BrokerClient, error) {
	client := &BrokerClient{
		endpoint:    brokerEndpoint,
		client:      nil,
		poller:      nil,
		quitCh:      quitCh,
		sendCh:      sendCh,
		receiveCh:   receiveCh,
		timeCh:      timeCh,
		dataCh:      dataCh,
		timeout:     0 * time.Second,
		tryInterval: 0 * time.Second,
		retryLimit:  3,
		retries:     0,
	}

	if err := client.Start(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *BrokerClient) handleIncomingReplies() {
	for {
		select {
		case <-c.quitCh:
			lep, _ := c.client.GetLastEndpoint()
			logz.Logger().Debug("Encerrando a goroutine de recepção de respostas...", map[string]interface{}{
				"context":        "handleIncomingReplies",
				"brokerEndpoint": lep,
				"timeout":        c.timeout,
			})
			return
		case reply := <-c.receiveCh:
			if len(reply) < 1 {
				logz.Logger().Warn("Resposta recebida do broker (não desserializável)", map[string]interface{}{
					"context": "handleIncomingReplies",
					"reply":   reply,
					"timeout": c.timeout,
					"retries": c.retries,
				})
				continue
			}

			var payload interface{}
			if err := json.Unmarshal([]byte(reply[0]), &payload); err != nil {
				logz.Logger().Error("Erro ao decodificar payload", map[string]interface{}{
					"context": "handleIncomingReplies",
					"reply":   reply,
					"timeout": c.timeout,
					"retries": c.retries,
					"error":   err,
				})
				continue
			}

			c.retries = 0
			p := payload.(map[string]interface{})
			if _, ok := p["type"]; !ok {
				logz.Logger().Error("Tipo de payload não encontrado", map[string]interface{}{
					"context": "handleIncomingReplies",
					"payload": payload,
					"timeout": c.timeout,
					"retries": c.retries,
				})
				continue
			}
			tp := p["type"].(string)
			if _, ok := p["data"]; !ok {
				logz.Logger().Error("Dados de payload não encontrados", map[string]interface{}{
					"context": "handleIncomingReplies",
					"payload": payload,
					"timeout": c.timeout,
					"retries": c.retries,
				})
				continue
			}
			dt := p["data"]

			logz.Logger().Debug("Payload decodificado, enviando para o canal de dados...", map[string]interface{}{
				"context": "handleIncomingReplies",
				"type":    tp,
				"data":    dt,
			})

			c.dataCh <- payload

			return
		case <-c.timeCh:
			if c.retries == c.retryLimit {
				logz.Logger().Warn("Número de tentativas excedido. Encerrando a goroutine de recepção de respostas...", map[string]interface{}{
					"context": "handleIncomingReplies",
					"timeout": c.timeout,
					"retries": c.retries,
				})
				close(c.quitCh)
				return
			} else {
				logz.Logger().Warn("Aguardando antes de reenviar mensagem...", map[string]interface{}{
					"context": "handleIncomingReplies",
					"timeout": c.timeout,
					"retries": c.retries,
				})

				time.Sleep(c.tryInterval)

				logz.Logger().Warn("Reenviando mensagem...", map[string]interface{}{
					"context": "handleIncomingReplies",
					"timeout": c.timeout,
					"retries": c.retries,
				})

				c.sendCh <- c.sendCh
			}
			continue
		case msg := <-c.sendCh:
			logz.Logger().Debug("Mensagem recebida para envio...", map[string]interface{}{
				"context": "handleIncomingReplies",
				"msg":     msg,
				"timeout": c.timeout,
				"retries": c.retries,
			})

			if err := c.trySendReceive(msg); err != nil {
				logz.Logger().Error("Erro ao enviar mensagem", map[string]interface{}{
					"context": "handleIncomingReplies",
					"msg":     msg,
					"timeout": c.timeout,
					"retries": c.retries,
					"error":   err,
				})
				return
			} else {
				logz.Logger().Debug("Mensagem enviada com sucesso", map[string]interface{}{
					"context": "handleIncomingReplies",
					"msg":     msg,
					"timeout": c.timeout,
					"retries": c.retries,
				})
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (c *BrokerClient) SendMessage(service string, requestPayload interface{}) error {
	payload, err := json.Marshal(requestPayload)
	if err != nil {
		logz.Logger().Error("Erro ao serializar payload", map[string]interface{}{
			"context":        "SendMessage",
			"service":        service,
			"requestPayload": requestPayload,
			"timeout":        c.timeout,
			"retries":        c.retries,
			"payload":        payload,
			"error":          err,
		})
		return err
	}

	req := []string{service, string(payload)}

	for i, frame := range req {
		if frame == "" {
			logz.Logger().Warn(fmt.Sprintf("Frame vazio detectado: posição %d", i), map[string]interface{}{
				"context":  "SendMessage",
				"frame":    frame,
				"position": i,
				"service":  service,
				"payload":  requestPayload,
				"timeout":  c.timeout,
				"retries":  c.retries,
			})
			return fmt.Errorf("frame vazio detectado")
		}
	}

	c.sendCh <- req

	return nil
}

func (c *BrokerClient) trySendReceive(req interface{}) error {
	if _, sendErr := c.client.SendMessage(req); sendErr != nil {
		logz.Logger().Error("Erro ao enviar mensagem", map[string]interface{}{
			"context": "SendMessage",
			"req":     req,
			"timeout": c.timeout,
			"retries": c.retries,
			"error":   sendErr,
		})
		return sendErr
	}

	if c.tryInterval > 0 {
		go func() {
			select {
			case <-time.After(c.timeout):
				c.timeCh <- time.Now()
				return
			}
		}()
	}

	c.retries++

	if reply, recvErr := c.client.RecvMessage(0); recvErr != nil {
		if errors.Is(recvErr, zmq4.EFSM) {
			c.handleEFSMError()
			logz.Logger().Warn("Reinicializando o socket devido a erro de máquina de estados (EFSM)...", map[string]interface{}{
				"context": "SendMessage",
				"reply":   reply,
				"req":     req,
				"timeout": c.timeout,
				"error":   recvErr,
			})
			return recvErr
		} else {
			logz.Logger().Error("Erro ao receber resposta", map[string]interface{}{
				"context": "SendMessage",
				"reply":   reply,
				"req":     req,
				"timeout": c.timeout,
				"retries": c.retries,
				"error":   recvErr,
			})
		}
		return recvErr
	} else {
		logz.Logger().Debug("Resposta recebida do broker", map[string]interface{}{
			"context": "SendMessage",
			"reply":   reply,
			"req":     req,
			"timeout": c.timeout,
			"retries": c.retries,
		})
		c.receiveCh <- reply
	}

	return nil
}

func (c *BrokerClient) handleEFSMError() {
	logz.Logger().Warn("Reinicializando o socket devido a erro de máquina de estados (EFSM)...", nil)
	closeErr := c.client.Close()
	if closeErr != nil {
		logz.Logger().Error("Erro ao fechar o socket", map[string]interface{}{
			"context": "handleEFSMError",
			"error":   closeErr,
		})
		return
	}
	newSocket, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		logz.Logger().Error("Erro ao recriar socket", map[string]interface{}{
			"context": "handleEFSMError",
			"error":   err,
		})
		return
	}
	c.client = newSocket
	connectErr := c.client.Connect(defaultBrokerEndpoint)
	if connectErr != nil {
		logz.Logger().Error("Erro ao reconectar ao broker", map[string]interface{}{
			"context": "handleEFSMError",
			"error":   connectErr,
		})
		return
	}
}

func (c *BrokerClient) Start() error {
	if c.client != nil {
		if c.poller != nil {
			remPollerErr := c.poller.RemoveBySocket(c.client)
			if remPollerErr != nil {
				logz.Logger().Warn("Erro ao remover socket do poller (existentes).", map[string]interface{}{
					"context": "Start",
					"error":   remPollerErr,
				})
			}
		}

		if c.quitCh != nil {
			c.quitCh <- true
			close(c.quitCh)
		}
		if c.sendCh != nil {
			close(c.sendCh)
		}
		if c.timeCh != nil {
			close(c.timeCh)
		}
		if c.receiveCh != nil {
			close(c.receiveCh)
		}
		if c.dataCh != nil {
			close(c.dataCh)
		}
	}

	c.retries = 0

	var err error
	c.client, err = zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		logz.Logger().Error(fmt.Sprintf("Erro ao criar socket: %v", err), nil)
		return err
	}
	c.poller = zmq4.NewPoller()

	if err := c.client.SetIdentity("gospyder-client"); err != nil {
		return err
	}

	if c.endpoint == "" {
		c.endpoint = defaultBrokerEndpoint
	}
	if err := c.client.Connect(c.endpoint); err != nil {
		logz.Logger().Error(fmt.Sprintf("Erro ao conectar ao broker: %v", err), nil)
		return err
	}

	if c.timeout > 0 {
		if err = c.client.SetRcvtimeo(c.timeout * time.Second); err != nil {
			return err
		}
		if err = c.client.SetSndtimeo(c.timeout * time.Second); err != nil {
			return err
		}
	}

	c.poller.Add(c.client, zmq4.POLLIN)

	go c.handleIncomingReplies()

	return nil
}

func (c *BrokerClient) Status() map[string]interface{} {
	var sktStatus string
	var sktState zmq4.State
	var sktStateErr error
	if c.client != nil {
		sktState, sktStateErr = c.client.GetEvents()
		if sktStateErr != nil {
			sktStatus = sktStateErr.Error()
		} else {
			sktStatus = sktState.String()
		}
	} else {
		sktStatus = "not initialized"
	}
	status := map[string]interface{}{
		"context":     "Status",
		"timeout":     c.timeout,
		"tryInterval": c.tryInterval,
		"retryLimit":  c.retryLimit,
		"endpoint":    c.endpoint,
		"retries":     c.retries,
		"status":      sktStatus,
	}
	logz.Logger().Debug("Status do cliente", status)
	return status
}

func (c *BrokerClient) Stop() {
	close(c.quitCh)
	_ = c.client.Close()
}

func (c *BrokerClient) Reset() {
	logz.Logger().Warn("Reinicializando canais...", map[string]interface{}{
		"context": "handleIncomingReplies",
		"timeout": c.timeout,
		"retries": c.retries,
	})

	c.retries = 0
	c.sendCh = make(chan interface{}, 100)
	c.receiveCh = make(chan []string, 100)
	c.timeCh = make(chan time.Time, 1)
	c.quitCh = make(chan bool, 1)

	logz.Logger().Warn("Canais reinicializados", map[string]interface{}{
		"context": "handleIncomingReplies",
		"timeout": c.timeout,
		"retries": c.retries,
	})
}

func (c *BrokerClient) Restart() {
	c.Stop()
	c.Reset()
	startErr := c.Start()
	if startErr != nil {
		logz.Logger().Error("Erro ao reiniciar o cliente", map[string]interface{}{
			"context": "Restart",
			"error":   startErr,
		})
		return
	}
}
