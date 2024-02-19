# Signals

## Introduction
Application signals allow interaction between the core application and its extensions.

## Components
The signal system consists of these omponents:
- `SignalConnector`: Manages connections and fire signals
- `SignalConnection`: Represents a connection between a signal type and a handler
- `SignalHandler`: The connected function that is fired when an event happens.

## Usage

### Connecting signals
You can connect handlers to respond to specific signal types using the `Connect` method of `SignalConnector`:
```go
connection := sc.Connect("onSignal",func(data interface{})) {
    // Handle the signal here
}
```

### Disconnecting signals
To disconnect the signal handler, you can call the `Disconnect` function provided by `SignalConnection`:
```go
connection.Disconnect()
```

## Firing signals
To fire a signal along with optional data, you can use the `Fire` method of `SignalConnector`:
```go
sc.Fire("onSignal")
```

### Waiting for signals
You can wait for a signal to be fired using the `Wait` method of `SignalConnector`. This method is currently unstable at the moment due to inconsistent unit testing results.
```go
sc.Wait("onSignal")
```

### Creating custom signals
To create a custom event when a signal is fired, you can define your own signal type and connect handlers to them as shown above.
