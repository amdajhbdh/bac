package streaming

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(nil, nil))

type Topic struct {
	Name        string
	Partitions  int
	Retention   time.Duration
	Compression string
}

type Message struct {
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	Timestamp time.Time       `json:"timestamp"`
	Topic     string          `json:"topic"`
	Partition int             `json:"partition"`
}

type Producer struct {
	topic      string
	partitions int
	channels   []chan Message
	mu         sync.RWMutex
}

func NewProducer(topic string, partitions int) *Producer {
	channels := make([]chan Message, partitions)
	for i := 0; i < partitions; i++ {
		channels[i] = make(chan Message, 1000)
	}
	return &Producer{
		topic:      topic,
		partitions: partitions,
		channels:   channels,
	}
}

func (p *Producer) Send(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	partition := p.partition(key)
	msg := Message{
		Key:       key,
		Value:     data,
		Timestamp: time.Now(),
		Topic:     p.topic,
		Partition: partition,
	}

	select {
	case p.channels[partition] <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) partition(key string) int {
	if key == "" {
		return int(time.Now().UnixNano() % int64(p.partitions))
	}
	h := 0
	for _, c := range key {
		h = h*31 + int(c)
	}
	return h % p.partitions
}

func (p *Producer) Channel(partition int) <-chan Message {
	return p.channels[partition]
}

type Consumer struct {
	topic     string
	partition int
	ch        <-chan Message
	position  int64
}

func NewConsumer(topic string, partition int, ch <-chan Message) *Consumer {
	return &Consumer{
		topic:     topic,
		partition: partition,
		ch:        ch,
	}
}

func (c *Consumer) Poll(ctx context.Context, timeout time.Duration) ([]Message, error) {
	var messages []Message

	select {
	case msg := <-c.ch:
		messages = append(messages, msg)
	case <-ctx.Done():
		return messages, ctx.Err()
	case <-time.After(timeout):
	}

	return messages, nil
}

func (c *Consumer) Position() int64 {
	return c.position
}

type TopicManager struct {
	topics map[string]*Topic
	mu     sync.RWMutex
}

func NewTopicManager() *TopicManager {
	return &TopicManager{
		topics: make(map[string]*Topic),
	}
}

func (tm *TopicManager) Create(ctx context.Context, name string, partitions int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.topics[name]; exists {
		return nil
	}

	tm.topics[name] = &Topic{
		Name:        name,
		Partitions:  partitions,
		Retention:   24 * time.Hour,
		Compression: "gzip",
	}

	logger.Info("topic created", "name", name, "partitions", partitions)
	return nil
}

func (tm *TopicManager) Delete(ctx context.Context, name string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	delete(tm.topics, name)
	logger.Info("topic deleted", "name", name)
	return nil
}

func (tm *TopicManager) List(ctx context.Context) ([]string, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	names := make([]string, 0, len(tm.topics))
	for name := range tm.topics {
		names = append(names, name)
	}
	return names, nil
}

type StreamProcessor struct {
	inputTopic  string
	outputTopic string
	producer    *Producer
	consumers   []*Consumer
	transforms  []TransformFunc
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

type TransformFunc func(Message) (Message, error)

func NewStreamProcessor(inputTopic, outputTopic string, partitions int) *StreamProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &StreamProcessor{
		inputTopic:  inputTopic,
		outputTopic: outputTopic,
		producer:    NewProducer(outputTopic, partitions),
		consumers:   make([]*Consumer, partitions),
		transforms:  make([]TransformFunc, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (sp *StreamProcessor) AddTransform(fn TransformFunc) *StreamProcessor {
	sp.transforms = append(sp.transforms, fn)
	return sp
}

func (sp *StreamProcessor) RegisterInput(partition int, ch <-chan Message) {
	sp.consumers[partition] = NewConsumer(sp.inputTopic, partition, ch)
}

func (sp *StreamProcessor) Start(ctx context.Context) error {
	for i, consumer := range sp.consumers {
		if consumer == nil {
			continue
		}
		sp.wg.Add(1)
		go func(c *Consumer, partition int) {
			defer sp.wg.Done()
			for {
				select {
				case <-sp.ctx.Done():
					return
				case <-ctx.Done():
					return
				default:
					msgs, err := c.Poll(sp.ctx, 100*time.Millisecond)
					if err != nil {
						if sp.ctx.Err() != nil {
							return
						}
						logger.Error("poll error", "error", err)
						continue
					}

					for _, msg := range msgs {
						processed := msg
						var err error
						for _, fn := range sp.transforms {
							processed, err = fn(processed)
							if err != nil {
								logger.Error("transform error", "error", err)
								break
							}
						}
						if err == nil {
							sp.producer.Send(sp.ctx, processed.Key, processed.Value)
						}
					}
				}
			}
		}(consumer, i)
	}
	logger.Info("stream processor started", "input", sp.inputTopic, "output", sp.outputTopic)
	return nil
}

func (sp *StreamProcessor) Stop() {
	sp.cancel()
	sp.wg.Wait()
	logger.Info("stream processor stopped")
}

type PipelineConfig struct {
	Brokers       []string
	TopicPrefix   string
	Partitions    int
	ConsumerGroup string
}

func NewPipelineConfig() PipelineConfig {
	return PipelineConfig{
		Brokers:       []string{"localhost:9092"},
		TopicPrefix:   "bac-",
		Partitions:    8,
		ConsumerGroup: "bac-pipeline",
	}
}
