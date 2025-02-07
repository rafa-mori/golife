package golife

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func (lm *GWebLifeCycle) StartIPCServer() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor IPC:", err)
		return
	}
	defer ln.Close()
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

func (lm *GWebLifeCycle) handleIPCConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		switch message {
		case "START":
			lm.Start()
			conn.Write([]byte("Processos iniciados\n"))
		case "STOP":
			lm.Stop()
			conn.Write([]byte("Processos parados\n"))
		case "RESTART":
			lm.Restart()
			conn.Write([]byte("Processos reiniciados\n"))
		case "STATUS":
			status := lm.Status()
			conn.Write([]byte(fmt.Sprintf("Status: %s\n", status)))
		default:
			conn.Write([]byte("Comando desconhecido\n"))
		}
	}
}
