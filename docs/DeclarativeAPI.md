# Declarative API

The Declarative API feature of GoLife provides an intuitive way to manage processes using a declarative approach. This allows you to define the desired state of your processes and let GoLife handle the rest.

## Examples

### Defining a Process

```go
package main

import (
	"fmt"
	"lifecycle"
)

func main() {
	manager := lifecycle.NewLifecycleManager()

	manager.DefineStage("start").
		OnEnter(func() { fmt.Println("Server started") }).
		OnExit(func() { fmt.Println("Server stopped") })

	manager.DefineStage("processing").
		AutoScale(3).
		OnEvent("request", func(data interface{}) {
			fmt.Println("Processing:", data)
		})

	manager.Trigger("processing", "request", "Request 1")
	manager.Trigger("processing", "request", "Request 2")
}
```

### Using the CLI

You can also use the CLI to manage processes declaratively.

```sh
golife start --name myProcess --cmd "myCommand" --args "arg1,arg2" --stages "start,processing" --triggers "request"
```

This command will start a process named `myProcess` with the command `myCommand` and arguments `arg1` and `arg2`. It will define two stages: `start` and `processing`, and set up a trigger for the `request` event.

### Responding to Events

The Declarative API allows you to respond to events in real-time.

```go
package main

import (
	"fmt"
	"lifecycle"
)

func main() {
	manager := lifecycle.NewLifecycleManager()

	manager.DefineStage("start").
		OnEnter(func() { fmt.Println("Server started") }).
		OnExit(func() { fmt.Println("Server stopped") })

	manager.DefineStage("processing").
		AutoScale(3).
		OnEvent("request", func(data interface{}) {
			fmt.Println("Processing:", data)
		})

	manager.Trigger("processing", "request", "Request 1")
	manager.Trigger("processing", "request", "Request 2")
}
```

In this example, the `processing` stage is set up to handle `request` events. When a `request` event is triggered, the provided function will be executed with the event data.

## Conclusion

The Declarative API feature of GoLife simplifies process management by allowing you to define the desired state of your processes and respond to events in real-time. Whether you are using the CLI or the embedded module, the Declarative API provides an intuitive and powerful way to manage your processes.
