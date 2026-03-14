package streaming

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewStreamProcessor(t *testing.T) {
	cfg := DefaultStreamConfig()
	sp := NewStreamProcessor(cfg)
	defer sp.Stop()

	topics := sp.Topics()

	expected := []string{
		TopicOCRInput,
		TopicOCROutput,
		TopicSolverInput,
		TopicSolverOutput,
		TopicStorage,
	}

	if len(topics) != len(expected) {
		t.Errorf("expected %d topics, got %d", len(expected), len(topics))
	}
}

func TestPublishAndSubscribe(t *testing.T) {
	cfg := DefaultStreamConfig()
	sp := NewStreamProcessor(cfg)
	defer sp.Stop()

	msg := Message{
		ID:    "test-1",
		Topic: TopicOCRInput,
		Key:   "key1",
		Value: json.RawMessage(`{"test":"data"}`),
	}

	err := sp.Publish(TopicOCRInput, msg)
	if err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	select {
	case received := <-sp.Subscribe(TopicOCRInput):
		if received.ID != msg.ID {
			t.Errorf("expected ID %s, got %s", msg.ID, received.ID)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for message")
	}
}

func TestWorkerProcessing(t *testing.T) {
	cfg := DefaultStreamConfig()
	sp := NewStreamProcessor(cfg)
	defer sp.Stop()

	sp.StartWorker(TopicOCRInput, func(m Message) Message {
		m.Key = "processed"
		m.Topic = TopicOCROutput
		return m
	})

	msg := Message{
		ID:    "test-worker",
		Topic: TopicOCRInput,
		Key:   "original",
	}
	sp.Publish(TopicOCRInput, msg)

	select {
	case received := <-sp.Subscribe(TopicOCROutput):
		if received.Key != "processed" {
			t.Errorf("expected key 'processed', got '%s'", received.Key)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timeout waiting for processed message")
	}
}

func TestBatchWorker(t *testing.T) {
	cfg := StreamConfig{
		BufferSize:    100,
		Workers:       1,
		BatchSize:     3,
		FlushInterval: 100 * time.Millisecond,
	}
	sp := NewStreamProcessor(cfg)
	defer sp.Stop()

	batch := make([]Message, 0)
	sp.StartBatchWorker(TopicOCRInput, func(msgs []Message) []Message {
		batch = append(batch, msgs...)
		for i := range msgs {
			msgs[i].Topic = TopicOCROutput
		}
		return msgs
	})

	for i := 0; i < 3; i++ {
		sp.Publish(TopicOCRInput, Message{ID: "batch-test"})
	}

	<-time.After(200 * time.Millisecond)

	if len(batch) != 3 {
		t.Errorf("expected batch size 3, got %d", len(batch))
	}
}

func TestStreamError(t *testing.T) {
	err := NewStreamError("test-topic", "buffer full")
	if err.Error() != "stream error: buffer full on topic test-topic" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestStopDrainsChannels(t *testing.T) {
	cfg := DefaultStreamConfig()
	sp := NewStreamProcessor(cfg)

	sp.StartWorker(TopicOCRInput, func(m Message) Message {
		return m
	})

	sp.Publish(TopicOCRInput, Message{ID: "before-stop"})

	<-time.After(50 * time.Millisecond)
	sp.Stop()
}

func TestPublishToNonexistentTopic(t *testing.T) {
	cfg := DefaultStreamConfig()
	sp := NewStreamProcessor(cfg)
	defer sp.Stop()

	err := sp.Publish("nonexistent-topic", Message{ID: "test"})
	if err == nil {
		t.Error("expected error for nonexistent topic")
	}
}
