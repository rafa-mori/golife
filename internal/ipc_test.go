package internal

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

type mockConn struct {
	net.Conn
	readBuffer  *bufio.Reader
	writeBuffer *bufio.Writer
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func newMockConn(input string) *mockConn {
	return &mockConn{
		readBuffer:  bufio.NewReader(strings.NewReader(input)),
		writeBuffer: bufio.NewWriter(&strings.Builder{}),
	}
}

func TestHandleIPCConnection_ValidToken(t *testing.T) {
	lm := &gWebLifeCycle{}
	conn := newMockConn("secreto-token-123\nSTATUS\n")

	lm.handleIPCConnection(conn)

	output := string(conn.writeBuffer.AvailableBuffer())
	if !strings.Contains(output, "Status:") {
		t.Errorf("Expected status response, got %s", output)
	}
}

func TestHandleIPCConnection_InvalidToken(t *testing.T) {
	lm := &gWebLifeCycle{}
	conn := newMockConn("invalid-token\n")

	lm.handleIPCConnection(conn)

	output := string(conn.writeBuffer.AvailableBuffer())
	if !strings.Contains(output, "Erro: Acesso negado! Token inválido.") {
		t.Errorf("Expected access denied response, got %s", output)
	}
}

func TestHandleIPCConnection_UnknownCommand(t *testing.T) {
	lm := &gWebLifeCycle{}
	conn := newMockConn("secreto-token-123\nUNKNOWN\n")

	lm.handleIPCConnection(conn)

	output := string(conn.writeBuffer.AvailableBuffer())
	if !strings.Contains(output, "Comando desconhecido") {
		t.Errorf("Expected unknown command response, got %s", output)
	}
}

func TestHandleIPCConnection_StartCommand(t *testing.T) {
	lm := &gWebLifeCycle{}
	conn := newMockConn("secreto-token-123\nSTART\n")

	lm.handleIPCConnection(conn)

	output := string(conn.writeBuffer.AvailableBuffer())
	if !strings.Contains(output, "Processos iniciados") {
		t.Errorf("Expected start response, got %s", output)
	}
}

func TestHandleIPCConnection_StopCommand(t *testing.T) {
	lm := &gWebLifeCycle{}
	conn := newMockConn("secreto-token-123\nSTOP\n")

	lm.handleIPCConnection(conn)

	output := string(conn.writeBuffer.AvailableBuffer())
	if !strings.Contains(output, "Processos parados") {
		t.Errorf("Expected stop response, got %s", output)
	}
}

func TestHandleIPCConnection_RestartCommand(t *testing.T) {
	lm := &gWebLifeCycle{}
	conn := newMockConn("secreto-token-123\nRESTART\n")

	lm.handleIPCConnection(conn)

	output := string(conn.writeBuffer.AvailableBuffer())
	if !strings.Contains(output, "Processos reiniciados") {
		t.Errorf("Expected restart response, got %s", output)
	}
}
