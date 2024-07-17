package connection

import (
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCPool struct {
	address 	string
	mu      	sync.Mutex
	client 		*grpc.ClientConn
	lastUsed 	time.Time
	maxIdleTime time.Duration
}

func NewGRPCPool(address string) *GRPCPool {
	pool := &GRPCPool{
		address: 		address,
		maxIdleTime:	30 * time.Second,
	}

	go pool.cleanupLoop()

	return pool
}

func (p *GRPCPool) Get() (*grpc.ClientConn, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.client == nil {
        client, err := grpc.NewClient(p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
            return nil, err
        }
        p.client = client
    }
    p.lastUsed = time.Now()
    return p.client, nil
}

func (p *GRPCPool) cleanupLoop() {
    for {
        time.Sleep(p.maxIdleTime)
        p.mu.Lock()
        if p.client != nil && time.Since(p.lastUsed) > p.maxIdleTime {
            log.Printf("remove grpc connection to %s", p.address)
            _ = p.client.Close()
            p.client = nil
        }
        p.mu.Unlock()
    }
}
