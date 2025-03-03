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

func TestTrigger(t *testing.T) {
	manager := NewLifecycleManager(nil, nil, nil, nil, nil, nil)
	manager.RegisterEvent("testEvent", "testStage")

	triggered := false
	manager.Trigger("testStage", "testEvent", nil)
	if !triggered {
		t.Errorf("Expected event to be triggered")
	}
}

func TestRegisterEvent(t *testing.T) {
	manager := NewLifecycleManager(nil, nil, nil, nil, nil, nil)
	err := manager.RegisterEvent("testEvent", "testStage")
	if err != nil {
		t.Errorf("Expected event to be registered, got error: %v", err)
	}
}

func TestRemoveEvent(t *testing.T) {
	manager := NewLifecycleManager(nil, nil, nil, nil, nil, nil)
	manager.RegisterEvent("testEvent", "testStage")
	err := manager.RemoveEvent("testEvent", "testStage")
	if err != nil {
		t.Errorf("Expected event to be removed, got error: %v", err)
	}
}

func TestStopEvents(t *testing.T) {
	manager := NewLifecycleManager(nil, nil, nil, nil, nil, nil)
	manager.RegisterEvent("testEvent", "testStage")
	err := manager.StopEvents()
	if err != nil {
		t.Errorf("Expected all events to be stopped, got error: %v", err)
	}
}
