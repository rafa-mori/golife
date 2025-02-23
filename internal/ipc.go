package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type IPCSettings struct {
	AuthToken string `json:"auth_token"`
}

var authToken string

func (lm *gWebLifeCycle) StartIPCServer() error {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		fmt.Println("Erro ao obter diretório home:", homeDirErr)
		return homeDirErr
	}

	cfgDir := filepath.Join(homeDir, ".kubex", ".golife")
	cfgPath := filepath.Join(cfgDir, "golife.conf")

	ipcSettings := IPCSettings{}
	if _, statErr := os.Stat(cfgPath); statErr != nil {
		if os.IsNotExist(statErr) && os.IsPermission(statErr) && os.IsExist(statErr) {
			return statErr
		} else {
			fmt.Println("Erro ao verificar existência do arquivo de configuração:", statErr)
			return statErr
		}
	} else {
		readToken, readTokenErr := os.ReadFile("config.json")
		if readTokenErr != nil {
			fmt.Println("Erro ao ler token de configuração:", readTokenErr)
			return readTokenErr
		}
		ipcSettingsErr := json.Unmarshal(readToken, &ipcSettings)
		if ipcSettingsErr != nil {
			fmt.Println("Erro ao decodificar token de configuração:", ipcSettingsErr)
			return ipcSettingsErr
		}
	}

	authToken = ipcSettings.AuthToken

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor IPC:", err)
		return err
	}
	defer func(ln net.Listener) {
		_ = ln.Close()
	}(ln)
	fmt.Println("Servidor IPC escutando na porta 8081")

	for {
		conn, connErr := ln.Accept()
		if connErr != nil {
			fmt.Println("Erro ao aceitar conexão:", connErr)
			continue
		}
		go lm.handleIPCConnection(conn)
	}
}
func (lm *gWebLifeCycle) handleIPCConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	reader := bufio.NewReader(conn)

	// Primeira mensagem deve ser um token válido
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	if token != authToken {
		_, wrtErr := conn.Write([]byte("Erro: Acesso negado! Token inválido.\n"))
		if wrtErr != nil {
			return
		}
		fmt.Println("Conexão recusada: Token inválido.")
		return
	}

	fmt.Println("Conexão autenticada com sucesso.")

	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		switch message {
		case "START":
			startErr := lm.Start()
			if startErr != nil {
				return
			}
			_, wrtErr := conn.Write([]byte("Processos iniciados\n"))
			if wrtErr != nil {
				return
			}
		case "STOP":
			stopErr := lm.Stop()
			if stopErr != nil {
				return
			}
			_, wrtErr := conn.Write([]byte("Processos parados\n"))
			if wrtErr != nil {
				return
			}
		case "RESTART":
			restartErr := lm.Restart()
			if restartErr != nil {
				return
			}
			_, wrtErr := conn.Write([]byte("Processos reiniciados\n"))
			if wrtErr != nil {
				return
			}
		case "STATUS":
			status := lm.Status()
			_, wrtErr := conn.Write([]byte(fmt.Sprintf("Status: %s\n", status)))
			if wrtErr != nil {
				return
			}
		default:
			_, wrtErr := conn.Write([]byte("Comando desconhecido\n"))
			if wrtErr != nil {
				return
			}
		}
	}
}
