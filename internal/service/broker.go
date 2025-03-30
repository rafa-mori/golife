package service

import (
	"errors"
	"fmt"
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

// BrokerClient represents a client for the broker.
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

// NewBrokerClient creates a new broker client.
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

// handleIncomingReplies handles incoming replies from the broker.
func (c *BrokerClient) handleIncomingReplies() {
	for {
		select {
		case <-c.quitCh:
			lep, _ := c.client.GetLastEndpoint()
			log.GetLogger().Debug("Shutting down the reply reception goroutine...", map[string]interface{}{
				"context":        "handleIncomingReplies",
				"brokerEndpoint": lep,
				"timeout":        c.timeout,
			})
			return
		case reply := <-c.receiveCh:
			if len(reply) < 1 {
				log.GetLogger().Warn("Reply received from broker (not deserializable)", map[string]interface{}{
					"context": "handleIncomingReplies",
					"reply":   reply,
					"timeout": c.timeout,
					"retries": c.retries,
				})
				continue
			}

			var payload interface{}
			if err := json.Unmarshal([]byte(reply[0]), &payload); err != nil {
				log.GetLogger().Error("Error decoding payload", map[string]interface{}{
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
				log.GetLogger().Error("Payload type not found", map[string]interface{}{
					"context": "handleIncomingReplies",
					"payload": payload,
					"timeout": c.timeout,
					"retries": c.retries,
				})
				continue
			}
			tp := p["type"].(string)
			if _, ok := p["data"]; !ok {
				log.GetLogger().Error("Payload data not found", map[string]interface{}{
					"context": "handleIncomingReplies",
					"payload": payload,
					"timeout": c.timeout,
					"retries": c.retries,
				})
				continue
			}
			dt := p["data"]

			log.GetLogger().Debug("Payload decoded, sending to data channel...", map[string]interface{}{
				"context": "handleIncomingReplies",
				"type":    tp,
				"data":    dt,
			})

			c.dataCh <- payload

			return
		case <-c.timeCh:
			if c.retries == c.retryLimit {
				log.GetLogger().Warn("Retry limit exceeded. Shutting down the reply reception goroutine...", map[string]interface{}{
					"context": "handleIncomingReplies",
					"timeout": c.timeout,
					"retries": c.retries,
				})
				close(c.quitCh)
				return
			} else {
				log.GetLogger().Warn("Waiting before resending message...", map[string]interface{}{
					"context": "handleIncomingReplies",
					"timeout": c.timeout,
					"retries": c.retries,
				})

				time.Sleep(c.tryInterval)

				log.GetLogger().Warn("Resending message...", map[string]interface{}{
					"context": "handleIncomingReplies",
					"timeout": c.timeout,
					"retries": c.retries,
				})

				c.sendCh <- c.sendCh
			}
			continue
		case msg := <-c.sendCh:
			log.GetLogger().Debug("Message received for sending...", map[string]interface{}{
				"context": "handleIncomingReplies",
				"msg":     msg,
				"timeout": c.timeout,
				"retries": c.retries,
			})

			if err := c.trySendReceive(msg); err != nil {
				log.GetLogger().Error("Error sending message", map[string]interface{}{
					"context": "handleIncomingReplies",
					"msg":     msg,
					"timeout": c.timeout,
					"retries": c.retries,
					"error":   err,
				})
				return
			} else {
				log.GetLogger().Debug("Message sent successfully", map[string]interface{}{
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

// SendMessage sends a message to the broker.
func (c *BrokerClient) SendMessage(service string, requestPayload interface{}) error {
	payload, err := json.Marshal(requestPayload)
	if err != nil {
		log.GetLogger().Error("Error serializing payload", map[string]interface{}{
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
			log.GetLogger().Warn(fmt.Sprintf("Empty frame detected: position %d", i), map[string]interface{}{
				"context":  "SendMessage",
				"frame":    frame,
				"position": i,
				"service":  service,
				"payload":  requestPayload,
				"timeout":  c.timeout,
				"retries":  c.retries,
			})
			return fmt.Errorf("empty frame detected")
		}
	}

	c.sendCh <- req

	return nil
}

// trySendReceive tries to send and receive a message from the broker.
func (c *BrokerClient) trySendReceive(req interface{}) error {
	if _, sendErr := c.client.SendMessage(req); sendErr != nil {
		log.GetLogger().Error("Error sending message", map[string]interface{}{
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
			log.GetLogger().Warn("Reinitializing socket due to state machine error (EFSM)...", map[string]interface{}{
				"context": "SendMessage",
				"reply":   reply,
				"req":     req,
				"timeout": c.timeout,
				"error":   recvErr,
			})
			return recvErr
		} else {
			log.GetLogger().Error("Error receiving reply", map[string]interface{}{
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
		log.GetLogger().Debug("Reply received from broker", map[string]interface{}{
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

// handleEFSMError handles EFSM errors by reinitializing the socket.
func (c *BrokerClient) handleEFSMError() {
	log.GetLogger().Warn("Reinitializing socket due to state machine error (EFSM)...", nil)
	closeErr := c.client.Close()
	if closeErr != nil {
		log.GetLogger().Error("Error closing socket", map[string]interface{}{
			"context": "handleEFSMError",
			"error":   closeErr,
		})
		return
	}
	newSocket, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		log.GetLogger().Error("Error recreating socket", map[string]interface{}{
			"context": "handleEFSMError",
			"error":   err,
		})
		return
	}
	c.client = newSocket
	connectErr := c.client.Connect(defaultBrokerEndpoint)
	if connectErr != nil {
		log.GetLogger().Error("Error reconnecting to broker", map[string]interface{}{
			"context": "handleEFSMError",
			"error":   connectErr,
		})
		return
	}
}

// Start starts the broker client.
func (c *BrokerClient) Start() error {
	if c.client != nil {
		if c.poller != nil {
			remPollerErr := c.poller.RemoveBySocket(c.client)
			if remPollerErr != nil {
				log.GetLogger().Warn("Error removing socket from poller (existing).", map[string]interface{}{
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
		log.GetLogger().Error(fmt.Sprintf("Error creating socket: %v", err), nil)
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
		log.GetLogger().Error(fmt.Sprintf("Error connecting to broker: %v", err), nil)
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

// Status returns the status of the broker client.
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
	log.GetLogger().Debug("Client status", status)
	return status
}

// Stop stops the broker client.
func (c *BrokerClient) Stop() {
	close(c.quitCh)
	_ = c.client.Close()
}

// Reset resets the broker client channels.
func (c *BrokerClient) Reset() {
	log.GetLogger().Warn("Reinitializing channels...", map[string]interface{}{
		"context": "handleIncomingReplies",
		"timeout": c.timeout,
		"retries": c.retries,
	})

	c.retries = 0
	c.sendCh = make(chan interface{}, 100)
	c.receiveCh = make(chan []string, 100)
	c.timeCh = make(chan time.Time, 1)
	c.quitCh = make(chan bool, 1)

	log.GetLogger().Warn("Channels reinitialized", map[string]interface{}{
		"context": "handleIncomingReplies",
		"timeout": c.timeout,
		"retries": c.retries,
	})
}

// Restart restarts the broker client.
func (c *BrokerClient) Restart() {
	c.Stop()
	c.Reset()
	startErr := c.Start()
	if startErr != nil {
		log.GetLogger().Error("Error restarting client", map[string]interface{}{
			"context": "Restart",
			"error":   startErr,
		})
		return
	}
}
