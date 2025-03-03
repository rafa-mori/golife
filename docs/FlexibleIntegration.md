# Flexible Integration

## Introduction

Flexible Integration in GoLife allows you to use the system via CLI or as an embedded module. This feature provides versatility in how you can integrate GoLife into your existing workflows and applications.

## Using the CLI

The CLI provides a straightforward way to interact with GoLife. You can use various commands to manage the lifecycle of your processes, trigger events, and more.

### Example

```sh
# Start the application
golife start --name myApp --cmd "myAppCommand"

# Trigger an event
golife trigger --stage processing --event request --data "Request 1"

# Check the status of the application
golife status
```

In this example, we use the CLI to start an application, trigger an event, and check the status of the application.

## Using as an Embedded Module

You can also use GoLife as an embedded module in your Go applications. This allows you to programmatically manage the lifecycle of your processes and trigger events.

### Example

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
		OnEvent("request", func(data interface{}) {
			fmt.Println("Processing:", data)
		})

	manager.Trigger("processing", "request", "Request 1")
	manager.Trigger("processing", "request", "Request 2")
}
```

In this example, we use GoLife as an embedded module to define stages, trigger events, and manage the lifecycle of our application.

## Usage Snippets

### Using the CLI

```sh
# Start the application
golife start --name myApp --cmd "myAppCommand"

# Trigger an event
golife trigger --stage processing --event request --data "Request 1"

# Check the status of the application
golife status
```

### Using as an Embedded Module

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
		OnEvent("request", func(data interface{}) {
			fmt.Println("Processing:", data)
		})

	manager.Trigger("processing", "request", "Request 1")
	manager.Trigger("processing", "request", "Request 2")
}
```

## Conclusion

Flexible Integration in GoLife provides versatility in how you can integrate the system into your existing workflows and applications. Whether you prefer using the CLI or embedding GoLife as a module, you can easily manage the lifecycle of your processes and trigger events.

For more information, refer to the [GoLife documentation](../README.md).
