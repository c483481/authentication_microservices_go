package connection

import (
	"log"
	"net/rpc"
	"sync"
	"time"
)

type RPCPool struct {
	address     string
	mu          sync.Mutex
	client      *rpc.Client
	lastUsed    time.Time
	maxIdleTime time.Duration
}

func NewRPCPool(address string) *RPCPool {
    pool := &RPCPool{
        address:     address,
		maxIdleTime: 30 * time.Second,
    }

	go pool.cleanupLoop()

    return pool
}

func (p *RPCPool) Get() (*rpc.Client, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.client == nil {
        client, err := rpc.Dial("tcp", p.address)
        if err != nil {
            return nil, err
        }
        p.client = client
    }
    p.lastUsed = time.Now()
    return p.client, nil
}

func (p *RPCPool) cleanupLoop() {
    for {
        time.Sleep(p.maxIdleTime)
        p.mu.Lock()
        if p.client != nil && time.Since(p.lastUsed) > p.maxIdleTime {
            log.Printf("remove rpc connection to %s", p.address)
            _ = p.client.Close()
            p.client = nil
        }
        p.mu.Unlock()
    }
}

