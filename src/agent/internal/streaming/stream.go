package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

const (
	TopicOCRInput     = "ocr-input"
	TopicOCROutput    = "ocr-output"
	TopicSolverInput  = "solver-input"
	TopicSolverOutput = "solver-output"
	TopicStorage      = "storage-input"
)

type Message struct {
	ID        string          `json:"id"`
	Topic     string          `json:"topic"`
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	Timestamp time.Time       `json:"timestamp"`
}

type PipelineMessage struct {
	ID         string    `json:"id"`
	InputType  string    `json:"input_type"`
	Data       string    `json:"data"`
	Result     string    `json:"result"`
	Confidence float64   `json:"confidence"`
	Subject    string    `json:"subject"`
	Chapter    string    `json:"chapter"`
	Concepts   []string  `json:"concepts"`
	Model      string    `json:"model"`
	Timestamp  time.Time `json:"timestamp"`
	Error      string    `json:"error,omitempty"`
}

func NewOCRInputMessage(data, language string) Message {
	pmsg := PipelineMessage{
		ID:        generateID(),
		InputType: "ocr-input",
		Data:      data,
		Timestamp: time.Now(),
	}
	value, _ := json.Marshal(pmsg)
	return Message{
		ID:        pmsg.ID,
		Topic:     TopicOCRInput,
		Value:     value,
		Timestamp: time.Now(),
	}
}

func generateID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

type StreamConfig struct {
	BufferSize    int
	Workers       int
	BatchSize     int
	FlushInterval time.Duration
}

func DefaultStreamConfig() StreamConfig {
	return StreamConfig{
		BufferSize:    1000,
		Workers:       4,
		BatchSize:     50,
		FlushInterval: 500 * time.Millisecond,
	}
}

type StreamProcessor struct {
	config StreamConfig
	topics map[string]chan Message
	mu     sync.RWMutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	logger *slog.Logger
}

func NewStreamProcessor(cfg StreamConfig) *StreamProcessor {
	ctx, cancel := context.WithCancel(context.Background())

	topics := make(map[string]chan Message)
	for _, topic := range []string{
		TopicOCRInput, TopicOCROutput,
		TopicSolverInput, TopicSolverOutput,
		TopicStorage,
	} {
		topics[topic] = make(chan Message, cfg.BufferSize)
	}

	return &StreamProcessor{
		config: cfg,
		topics: topics,
		ctx:    ctx,
		cancel: cancel,
		logger: slog.Default(),
	}
}

func (sp *StreamProcessor) Topic(name string) chan<- Message {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return sp.topics[name]
}

func (sp *StreamProcessor) Subscribe(topic string) <-chan Message {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return sp.topics[topic]
}

func (sp *StreamProcessor) Publish(topic string, msg Message) error {
	sp.mu.RLock()
	ch, ok := sp.topics[topic]
	sp.mu.RUnlock()

	if !ok {
		return &StreamError{Topic: topic, Reason: "topic not found"}
	}

	select {
	case ch <- msg:
		return nil
	case <-sp.ctx.Done():
		return sp.ctx.Err()
	default:
		return &StreamError{Topic: topic, Reason: "buffer full"}
	}
}

func (sp *StreamProcessor) StartWorker(name string, handler func(Message) Message) {
	sp.wg.Add(1)
	go func() {
		defer sp.wg.Done()
		for {
			select {
			case <-sp.ctx.Done():
				return
			case msg, ok := <-sp.topics[name]:
				if !ok {
					return
				}
				result := handler(msg)
				if result.Topic != "" {
					sp.Publish(result.Topic, result)
				}
			}
		}
	}()
}

func (sp *StreamProcessor) StartBatchWorker(name string, handler func([]Message) []Message) {
	sp.wg.Add(1)
	go func() {
		defer sp.wg.Done()

		batch := make([]Message, 0, sp.config.BatchSize)
		ticker := time.NewTicker(sp.config.FlushInterval)
		defer ticker.Stop()

		flush := func() {
			if len(batch) == 0 {
				return
			}
			results := handler(batch)
			for _, result := range results {
				if result.Topic != "" {
					sp.Publish(result.Topic, result)
				}
			}
			batch = batch[:0]
		}

		for {
			select {
			case <-sp.ctx.Done():
				flush()
				return
			case msg, ok := <-sp.topics[name]:
				if !ok {
					flush()
					return
				}
				batch = append(batch, msg)
				if len(batch) >= sp.config.BatchSize {
					flush()
				}
			case <-ticker.C:
				flush()
			}
		}
	}()
}

func (sp *StreamProcessor) Stop() {
	sp.cancel()
	sp.wg.Wait()

	sp.mu.RLock()
	defer sp.mu.RUnlock()
	for _, ch := range sp.topics {
		close(ch)
	}
}

func (sp *StreamProcessor) Topics() []string {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	topics := make([]string, 0, len(sp.topics))
	for name := range sp.topics {
		topics = append(topics, name)
	}
	return topics
}

type StreamError struct {
	Topic  string
	Reason string
}

func (e *StreamError) Error() string {
	return "stream error: " + e.Reason + " on topic " + e.Topic
}

func NewStreamError(topic, reason string) *StreamError {
	return &StreamError{Topic: topic, Reason: reason}
}
