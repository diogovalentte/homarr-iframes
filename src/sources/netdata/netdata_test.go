package netdata

import (
	"sync"
	"testing"

	"github.com/diogovalentte/homarr-iframes/src/config"
)

// TestNewConcurrent must be run with -race. Each HTTP request is handled in its
// own goroutine, so New() can be called concurrently.
func TestNewConcurrent(t *testing.T) {
	config.GlobalConfigs = &config.Configs{}
	config.GlobalConfigs.NetdataConfigs.Address = "http://localhost:19999"
	config.GlobalConfigs.NetdataConfigs.Token = "token"
	n = nil

	instances := make([]*Netdata, 8)
	var wg sync.WaitGroup
	for i := range instances {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			instance, err := New()
			if err != nil {
				t.Error(err)
				return
			}
			instances[i] = instance
		}(i)
	}
	wg.Wait()

	for _, instance := range instances {
		if instance != instances[0] {
			t.Fatal("New() returned different instances")
		}
	}
}
