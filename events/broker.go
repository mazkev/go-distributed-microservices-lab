package events

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type EventHandler func(ctx context.Context, event Event) error

type EventBroker interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler)
	Close()
}

// AsyncChannelBroker adalah event broker in-memory berbasis Goroutines & Channels
// Bekerja mandiri tanpa ketergantungan broker eksternal
type AsyncChannelBroker struct {
	mu        sync.RWMutex
	handlers  map[string][]EventHandler
	eventChan chan Event
	workers   int
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

func NewAsyncChannelBroker(bufferSize int, workers int) *AsyncChannelBroker {
	b := &AsyncChannelBroker{
		handlers:  make(map[string][]EventHandler),
		eventChan: make(chan Event, bufferSize),
		workers:   workers,
		stopChan:  make(chan struct{}),
	}

	// Start Background Worker Pool
	for i := 1; i <= workers; i++ {
		b.wg.Add(1)
		go b.startWorker(i)
	}

	return b
}

func (b *AsyncChannelBroker) startWorker(workerID int) {
	defer b.wg.Done()

	for {
		select {
		case <-b.stopChan:
			return
		case event, ok := <-b.eventChan:
			if !ok {
				return
			}
			b.dispatch(event)
		}
	}
}

func (b *AsyncChannelBroker) dispatch(event Event) {
	b.mu.RLock()
	handlers := b.handlers[event.EventType()]
	b.mu.RUnlock()

	for _, handler := range handlers {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := handler(ctx, event); err != nil {
			fmt.Printf("❌ [EVENT ERROR] Gagal memproses event %s: %v\n", event.EventType(), err)
		}
		cancel()
	}
}

func (b *AsyncChannelBroker) Publish(ctx context.Context, event Event) error {
	select {
	case b.eventChan <- event:
		return nil
	default:
		// Jika buffer penuh, proses di goroutine terpisah agar tidak memblok caller
		go func() {
			b.eventChan <- event
		}()
		return nil
	}
}

func (b *AsyncChannelBroker) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *AsyncChannelBroker) Close() {
	close(b.stopChan)
	close(b.eventChan)
	b.wg.Wait()
}
