package types

import (
	"fmt"
	ci "github.com/faelmori/golife/internal/components/interfaces"
	t "github.com/faelmori/golife/internal/components/types"
	l "github.com/faelmori/logz"
	"sync"
	"testing"
	"time"
)

var (
	start   time.Time
	elapsed time.Duration
)

func init() {
	start = time.Now()
}

func elapsedTimeLog() {
	elapsed = time.Since(start)
	fmt.Println("#####################################################################################")
	fmt.Println(fmt.Sprintf("Tempo total da execução: %bms\n", elapsed.Milliseconds()))
	fmt.Println(fmt.Sprintf("Tempo total da execução: %bs\n", elapsed.Seconds()))
	fmt.Println(fmt.Sprintf("Tempo total da execução: %bm\n", elapsed.Minutes()))
	fmt.Println("#####################################################################################")
}

func TestInitializeLifecycle(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")

	lifecycle := t.NewLifeCycle[string, ci.IProperty[ci.IProcessInput[string]]](nil, logger)

	if err := lifecycle.Initialize(); err == nil {
		ts.Fatalf("Lifecycle deveria falhar ao inicializar sem configuração válida")
	} else {
		ts.Logf("Lifecycle falhou corretamente ao inicializar sem configuração: %v", err)
	}
}

func TestAddProcess(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")
	lifecycle := t.NewLifeCycle[string, ci.IProperty[ci.IProcessInput[string]]](nil, logger)

	if rawProcManager, ok := lifecycle.GetComponent("process"); ok {
		if procManager, ok := rawProcManager.(*t.ProcessManager[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]]); ok {
			if err := procManager.AddProcess("test-process", nil); err == nil {
				ts.Fatalf("Lifecycle deveria falhar ao adicionar um processo inválido")
			} else {
				ts.Logf("Falha correta ao adicionar um processo inválido: %v", err)
			}
		} else {
			ts.Fatalf("Falha ao converter para ProcessManager: %v", ok)
		}
	}
}

func TestSetGetConfig(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")

	props := t.NewProcessInput[string]("test-process", "tail", []string{"-f", "/dev/null"}, false, false, nil, logger, false)

	processInput := t.NewProperty[ci.IProcessInput[string]](
		"test-config",
		&props,
		false,
		func(d any) (bool, error) {
			if d == nil {
				return false, fmt.Errorf("process input não pode ser nulo")
			}
			fmt.Println(fmt.Sprintf("Callback payload: %v", d))
			return true, nil
		},
	)

	lifecycle := t.NewLifeCycle[string, ci.IProperty[ci.IProcessInput[string]]](&processInput, logger)

	if err := lifecycle.Initialize(); err != nil {
		ts.Fatalf("Falha ao inicializar o lifecycle: %v", err)
	}

	lifecycle.SetConfig(processInput)
	rawResult := lifecycle.GetConfig()
	if rawResult == nil {
		ts.Fatalf("Falha ao recuperar configuração, esperava %v, recebeu nil", processInput)
	}
	resultProp := *rawResult

	if resultProp.GetValue() != processInput.GetValue() {
		ts.Fatalf("Falha ao recuperar configuração, esperava %v, recebeu %v", processInput, resultProp)
	}

	ts.Logf("Configuração definida e recuperada com sucesso!")
}

func TestSetInitializeWithoutInput(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")
	lifecycle := t.NewLifeCycle[string, ci.IProperty[ci.IProcessInput[string]]](nil, logger)

	if err := lifecycle.Initialize(); err != nil {
		ts.Logf("Falha ao inicializar o lifecycle: %v\n%s", err, "Que é o esperado, pois não foi inicializado corretamente. Iniciou com o input nil.")
	} else {
		ts.Fatalf("Falha ao inicializar o lifecycle: %v", err)
	}
}

func TestRestartLifecycle(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")
	lifecycle := t.NewLifeCycle[string, ci.IProperty[ci.IProcessInput[string]]](nil, logger)

	if err := lifecycle.RestartLifecycle(); err != nil {
		ts.Logf("Erro ao reiniciar lifecycle: %v\n%s", err, "Que é o esperado, pois não foi inicializado corretamente")
	} else {
		ts.Fatalf("Lifecycle reiniciado com sucesso, o que não deveria acontecer por não ter inicializado corretamente")
	}
}

func TestConcurrentProcessHandlingWith_5(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")

	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ { // Simulamos múltiplos processos
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			//processManaged := t.
			validationFunc := t.NewValidationFunc[ci.IManagedProcess[any]](0, func(value *ci.IManagedProcess[any], args ...any) ci.IValidationResult {
				if value == nil {
					return t.NewValidationResult(false, "Processo não pode ser nulo", fmt.Errorf("processo não pode ser nulo"))
				}
				return t.NewValidationResult(true, "Processo válido", nil)
			})

			propertyI := t.NewProcessInput[ci.IManagedProcess[any]](
				"TESTE_COM_INPUT",
				"tail",
				[]string{"-f", "/dev/null"},
				true,
				false,
				validationFunc,
				logger,
				false,
			)

			propertyT := t.NewProperty[ci.IProcessInput[ci.IManagedProcess[any]]]("TESTE_COM_INPUT", &propertyI, false, nil)

			lifecycle := t.NewLifeCycle[ci.IManagedProcess[any], ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](&propertyT, logger)

			rawProcessManager, ok := lifecycle.GetComponent("process")
			if !ok {
				ts.Errorf("Falha ao obter o componente de processo: %v", ok)
				return
			}

			if procManager, ok := rawProcessManager.(*t.ProcessManager[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]]); ok {
				if err := procManager.AddProcess(fmt.Sprintf("test-process-%d", i), propertyT); err != nil {
					ts.Errorf("Erro ao adicionar processo concorrente %d: %v", i, err)
				}
				if err := procManager.StartProcess(fmt.Sprintf("test-process-%d", i)); err != nil {
					ts.Errorf("Erro ao iniciar processo concorrente %d: %v", i, err)
				}
			}
		}(i)
	}

	wg.Wait()

	ts.Logf("Todos os 5 processos adicionados corretamente em concorrência!")
}

func TestConcurrentProcessHandlingWith_10(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ { // Simulamos múltiplos processos
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			//processManaged := t.
			validationFunc := t.NewValidationFunc[ci.IManagedProcess[any]](0, func(value *ci.IManagedProcess[any], args ...any) ci.IValidationResult {
				if value == nil {
					return t.NewValidationResult(false, "Processo não pode ser nulo", fmt.Errorf("processo não pode ser nulo"))
				}
				return t.NewValidationResult(true, "Processo válido", nil)
			})

			propertyI := t.NewProcessInput[ci.IManagedProcess[any]](
				"TESTE_COM_INPUT",
				"tail",
				[]string{"-f", "/dev/null"},
				true,
				false,
				validationFunc,
				logger,
				false,
			)

			propertyT := t.NewProperty[ci.IProcessInput[ci.IManagedProcess[any]]]("TESTE_COM_INPUT", &propertyI, false, nil)

			lifecycle := t.NewLifeCycle[ci.IManagedProcess[any], ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](&propertyT, logger)

			rawProcessManager, ok := lifecycle.GetComponent("process")
			if !ok {
				ts.Errorf("Falha ao obter o componente de processo: %v", ok)
				return
			}

			if procManager, ok := rawProcessManager.(*t.ProcessManager[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]]); ok {
				if err := procManager.AddProcess(fmt.Sprintf("test-process-%d", i), propertyT); err != nil {
					ts.Errorf("Erro ao adicionar processo concorrente %d: %v", i, err)
				}
				if err := procManager.StartProcess(fmt.Sprintf("test-process-%d", i)); err != nil {
					ts.Errorf("Erro ao iniciar processo concorrente %d: %v", i, err)
				}
			}
		}(i)
	}

	wg.Wait()

	ts.Logf("Todos os 10 processos adicionados corretamente em concorrência!")
}

func TestConcurrentProcessHandlingWith_30(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")

	wg := sync.WaitGroup{}
	for i := 0; i < 30; i++ { // Simulamos múltiplos processos
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			//processManaged := t.
			validationFunc := t.NewValidationFunc[ci.IManagedProcess[any]](0, func(value *ci.IManagedProcess[any], args ...any) ci.IValidationResult {
				if value == nil {
					return t.NewValidationResult(false, "Processo não pode ser nulo", fmt.Errorf("processo não pode ser nulo"))
				}
				return t.NewValidationResult(true, "Processo válido", nil)
			})

			propertyI := t.NewProcessInput[ci.IManagedProcess[any]](
				"TESTE_COM_INPUT",
				"tail",
				[]string{"-f", "/dev/null"},
				true,
				false,
				validationFunc,
				logger,
				false,
			)

			propertyT := t.NewProperty[ci.IProcessInput[ci.IManagedProcess[any]]]("TESTE_COM_INPUT", &propertyI, false, nil)

			lifecycle := t.NewLifeCycle[ci.IManagedProcess[any], ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](&propertyT, logger)

			rawProcessManager, ok := lifecycle.GetComponent("process")
			if !ok {
				ts.Errorf("Falha ao obter o componente de processo: %v", ok)
				return
			}

			if procManager, ok := rawProcessManager.(*t.ProcessManager[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]]); ok {
				if err := procManager.AddProcess(fmt.Sprintf("test-process-%d", i), propertyT); err != nil {
					ts.Errorf("Erro ao adicionar processo concorrente %d: %v", i, err)
				}
				if err := procManager.StartProcess(fmt.Sprintf("test-process-%d", i)); err != nil {
					ts.Errorf("Erro ao iniciar processo concorrente %d: %v", i, err)
				}
			}
		}(i)
	}

	wg.Wait()

	ts.Logf("Todos os 30 processos adicionados corretamente em concorrência!")
}

func TestConcurrentProcessHandlingWith_50(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")

	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ { // Simulamos múltiplos processos
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			//processManaged := t.
			validationFunc := t.NewValidationFunc[ci.IManagedProcess[any]](0, func(value *ci.IManagedProcess[any], args ...any) ci.IValidationResult {
				if value == nil {
					return t.NewValidationResult(false, "Processo não pode ser nulo", fmt.Errorf("processo não pode ser nulo"))
				}
				return t.NewValidationResult(true, "Processo válido", nil)
			})

			propertyI := t.NewProcessInput[ci.IManagedProcess[any]](
				"TESTE_COM_INPUT",
				"tail",
				[]string{"-f", "/dev/null"},
				true,
				false,
				validationFunc,
				logger,
				false,
			)

			propertyT := t.NewProperty[ci.IProcessInput[ci.IManagedProcess[any]]]("TESTE_COM_INPUT", &propertyI, false, nil)

			lifecycle := t.NewLifeCycle[ci.IManagedProcess[any], ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](&propertyT, logger)

			rawProcessManager, ok := lifecycle.GetComponent("process")
			if !ok {
				ts.Errorf("Falha ao obter o componente de processo: %v", ok)
				return
			}

			if procManager, ok := rawProcessManager.(*t.ProcessManager[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]]); ok {
				if err := procManager.AddProcess(fmt.Sprintf("test-process-%d", i), propertyT); err != nil {
					ts.Errorf("Erro ao adicionar processo concorrente %d: %v", i, err)
				}
				if err := procManager.StartProcess(fmt.Sprintf("test-process-%d", i)); err != nil {
					ts.Errorf("Erro ao iniciar processo concorrente %d: %v", i, err)
				}
			}
		}(i)
	}

	wg.Wait()

	ts.Logf("Todos os 50 processos adicionados corretamente em concorrência!")
}

func TestConcurrentProcessHandlingWith_100(ts *testing.T) {
	defer elapsedTimeLog()
	logger := l.GetLogger("TestLogger")

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ { // Simulamos múltiplos processos
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			//processManaged := t.
			validationFunc := t.NewValidationFunc[ci.IManagedProcess[any]](0, func(value *ci.IManagedProcess[any], args ...any) ci.IValidationResult {
				if value == nil {
					return t.NewValidationResult(false, "Processo não pode ser nulo", fmt.Errorf("processo não pode ser nulo"))
				}
				return t.NewValidationResult(true, "Processo válido", nil)
			})

			propertyI := t.NewProcessInput[ci.IManagedProcess[any]](
				"TESTE_COM_INPUT",
				"tail",
				[]string{"-f", "/dev/null"},
				true,
				false,
				validationFunc,
				logger,
				false,
			)

			propertyT := t.NewProperty[ci.IProcessInput[ci.IManagedProcess[any]]]("TESTE_COM_INPUT", &propertyI, false, nil)

			lifecycle := t.NewLifeCycle[ci.IManagedProcess[any], ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]](&propertyT, logger)

			rawProcessManager, ok := lifecycle.GetComponent("process")
			if !ok {
				ts.Errorf("Falha ao obter o componente de processo: %v", ok)
				return
			}

			if procManager, ok := rawProcessManager.(*t.ProcessManager[ci.IProperty[ci.IProcessInput[ci.IManagedProcess[any]]]]); ok {
				if err := procManager.AddProcess(fmt.Sprintf("test-process-%d", i), propertyT); err != nil {
					ts.Errorf("Erro ao adicionar processo concorrente %d: %v", i, err)
				}
				if err := procManager.StartProcess(fmt.Sprintf("test-process-%d", i)); err != nil {
					ts.Errorf("Erro ao iniciar processo concorrente %d: %v", i, err)
				}
			}
		}(i)
	}

	wg.Wait()

	ts.Logf("Todos os 100 processos adicionados corretamente em concorrência!")
}
