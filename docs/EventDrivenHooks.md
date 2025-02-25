# Event-Driven Hooks

Event-Driven Hooks in GoLife allow you to respond to real-time events within your application. This feature is essential for building reactive systems that can handle various events and trigger corresponding actions.

## Usage

### Registering an Event

To register an event, use the `RegisterEvent` method. This method associates an event with a specific stage.

```go
package main

import (
	"fmt"
	"lifecycle"
)

func main() {
	manager := lifecycle.NewLifecycleManager()

	manager.RegisterEvent("dataReceived", "processing")

	manager.DefineStage("processing").
		OnEvent("dataReceived", func(data interface{}) {
			fmt.Println("Data received:", data)
		})

	manager.Trigger("processing", "dataReceived", "Sample Data")
}
```

### Triggering an Event

To trigger an event, use the `Trigger` method. This method triggers the specified event in the given stage.

```go
package main

import (
	"fmt"
	"lifecycle"
)

func main() {
	manager := lifecycle.NewLifecycleManager()

	manager.RegisterEvent("dataReceived", "processing")

	manager.DefineStage("processing").
		OnEvent("dataReceived", func(data interface{}) {
			fmt.Println("Data received:", data)
		})

	manager.Trigger("processing", "dataReceived", "Sample Data")
}
```

### Removing an Event

To remove an event, use the `RemoveEvent` method. This method removes the specified event from the given stage.

```go
package main

import (
	"fmt"
	"lifecycle"
)

func main() {
	manager := lifecycle.NewLifecycleManager()

	manager.RegisterEvent("dataReceived", "processing")

	manager.DefineStage("processing").
		OnEvent("dataReceived", func(data interface{}) {
			fmt.Println("Data received:", data)
		})

	manager.RemoveEvent("dataReceived", "processing")
}
```

### Stopping All Events

To stop all events, use the `StopEvents` method. This method stops all registered events.

```go
package main

import (
	"lifecycle"
)

func main() {
	manager := lifecycle.NewLifecycleManager()

	manager.RegisterEvent("dataReceived", "processing")

	manager.DefineStage("processing").
		OnEvent("dataReceived", func(data interface{}) {
			fmt.Println("Data received:", data)
		})

	manager.StopEvents()
}
```

## Conclusion

Event-Driven Hooks in GoLife provide a powerful way to build reactive systems that can handle real-time events efficiently. By registering, triggering, removing, and stopping events, you can create a flexible and responsive application.

