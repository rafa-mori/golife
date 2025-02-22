package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var authToken = "secreto-token-123" // Token de autenticação para controle de acesso

func (lm *gWebLifeCycle) StartIPCServer() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor IPC:", err)
		return
	}
	defer func(ln net.Listener) {
		_ = ln.Close()
	}(ln)
	fmt.Println("Servidor IPC escutando na porta 8081")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
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
