package signals

import (
	"sync"
	"testing"
	"time"
)

func TestSignalConnector(t *testing.T) {
  connector := NewSignalConnector()

  // Waitgroup for syncing the test
  var wg sync.WaitGroup
  wg.Add(1)

  mockMod := SignalHandler(func(interface{}) {
    t.Log("EVENT FIRED")
    defer wg.Done()
  })

  connection := connector.Connect("test_signal",mockMod)
  connector.Fire("test_singal","test_data")

  connection.Disconnect()
  connector.Fire("test_signal","test_data")

  wg.Wait()

  // Test wait function
  go func() {
    t.Log("WAITING EVENT TO FIRE")
    time.Sleep(time.Second*3)
    connector.Fire("test_signal","test_data")
  }()

  connector.Wait("test_signal")
  t.Log("EVENT YIELDED")

}
