package signals

import "sync"

type SignalConnector struct {
	connected map[string][]SignalHandler
	waitMap   map[string][]chan struct{}
	mu        sync.Mutex
}

type SignalConnection struct {
	Disconnect func()
}

func NewSignalConnector() *SignalConnector {
	return &SignalConnector{
		connected: make(map[string][]SignalHandler),
		waitMap:   make(map[string][]chan struct{}),
	}
}

func (sc *SignalConnector) Connect(signalType string, handler SignalHandler) SignalConnection {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.connected[signalType] = append(sc.connected[signalType], handler)

	return SignalConnection{
		Disconnect: func() {
			sc.Disconnect(signalType, &handler)
		},
	}
}

func (sc *SignalConnector) Disconnect(signalType string, handler *SignalHandler) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if handlers, ok := sc.connected[signalType]; ok {
		for i, h := range handlers {
			if &h == handler {
				sc.connected[signalType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

func (sc *SignalConnector) Fire(signalType string, signalData interface{}) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Call handler functions
	for _, handler := range sc.connected[signalType] {
		go handler(signalData)
	}

	// Signal wait channels
	for _, ch := range sc.waitMap[signalType] {
		close(ch)
		delete(sc.waitMap, signalType)
	}
}

func (sc *SignalConnector) Wait(signalType string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Create wait channel
	ch := make(chan struct{})
	sc.waitMap[signalType] = append(sc.waitMap[signalType], ch)

	// Unlock to allow fire to execute
	sc.mu.Unlock()
	// Wait until fired
	<-ch

	// Relock after waiting
	sc.mu.Lock()

}
