package streaming

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/bac-unified/agent/internal/memory"
	"github.com/bac-unified/agent/internal/ocr"
	"github.com/bac-unified/agent/internal/solver"
)

var plogger = slog.Default()

type PipelineConfig struct {
	EnableOCR     bool
	EnableMemory  bool
	EnableSolver  bool
	EnableStorage bool
}

func DefaultPipelineConfig() PipelineConfig {
	return PipelineConfig{
		EnableOCR:     true,
		EnableMemory:  true,
		EnableSolver:  true,
		EnableStorage: true,
	}
}

type Pipeline struct {
	stream *StreamProcessor
	config PipelineConfig
	logger *slog.Logger
}

func NewPipeline(cfg PipelineConfig) *Pipeline {
	return &Pipeline{
		stream: NewStreamProcessor(DefaultStreamConfig()),
		config: cfg,
		logger: plogger,
	}
}

func (p *Pipeline) Start(ctx context.Context) {
	p.stream.StartWorker(TopicOCRInput, func(msg Message) Message {
		return p.processOCR(ctx, msg)
	})

	p.stream.StartWorker(TopicOCROutput, func(msg Message) Message {
		return p.processMemory(ctx, msg)
	})

	p.stream.StartWorker(TopicSolverInput, func(msg Message) Message {
		return p.processSolver(ctx, msg)
	})

	p.stream.StartWorker(TopicSolverOutput, func(msg Message) Message {
		return p.processStorage(ctx, msg)
	})

	p.logger.Info("pipeline started")
}

func (p *Pipeline) processOCR(ctx context.Context, msg Message) Message {
	var pmsg PipelineMessage
	if err := json.Unmarshal(msg.Value, &pmsg); err != nil {
		p.logger.Error("unmarshal error", "error", err)
		return msg
	}

	if !p.config.EnableOCR {
		msg.Topic = TopicOCROutput
		return msg
	}

	result := ocr.ProcessImage(ctx, pmsg.Data)

	response := PipelineMessage{
		ID:         pmsg.ID,
		InputType:  "ocr-result",
		Data:       result.Text,
		Result:     result.Text,
		Confidence: result.Confidence,
		Timestamp:  time.Now(),
	}

	data, _ := json.Marshal(response)
	msg.Value = data
	msg.Topic = TopicOCROutput

	return msg
}

func (p *Pipeline) processMemory(ctx context.Context, msg Message) Message {
	var pmsg PipelineMessage
	if err := json.Unmarshal(msg.Value, &pmsg); err != nil {
		p.logger.Error("unmarshal error", "error", err)
		return msg
	}

	if !p.config.EnableMemory {
		msg.Topic = TopicSolverInput
		return msg
	}

	result := memory.Lookup(ctx, pmsg.Data, 3)

	response := PipelineMessage{
		ID:        pmsg.ID,
		InputType: "memory-result",
		Data:      pmsg.Data,
		Result:    result.Context,
		Timestamp: time.Now(),
	}

	data, _ := json.Marshal(response)
	msg.Value = data
	msg.Topic = TopicSolverInput

	return msg
}

func (p *Pipeline) processSolver(ctx context.Context, msg Message) Message {
	var pmsg PipelineMessage
	if err := json.Unmarshal(msg.Value, &pmsg); err != nil {
		p.logger.Error("unmarshal error", "error", err)
		return msg
	}

	if !p.config.EnableSolver {
		msg.Topic = TopicSolverOutput
		return msg
	}

	result, err := solver.Solve(ctx, pmsg.Data, pmsg.Result)
	if err != nil {
		p.logger.Error("solver error", "error", err)
		msg.Topic = TopicSolverOutput
		return msg
	}

	response := PipelineMessage{
		ID:         pmsg.ID,
		InputType:  "solver-result",
		Data:       pmsg.Data,
		Result:     result.Solution,
		Subject:    result.Subject,
		Concepts:   result.Concepts,
		Confidence: result.Confidence,
		Model:      result.Model,
		Timestamp:  time.Now(),
	}

	data, _ := json.Marshal(response)
	msg.Value = data
	msg.Topic = TopicSolverOutput

	return msg
}

func (p *Pipeline) processStorage(ctx context.Context, msg Message) Message {
	var pmsg PipelineMessage
	if err := json.Unmarshal(msg.Value, &pmsg); err != nil {
		p.logger.Error("unmarshal error", "error", err)
		return msg
	}

	if !p.config.EnableStorage {
		msg.Topic = ""
		return msg
	}

	err := memory.Store(ctx, pmsg.Data, pmsg.Result, pmsg.Subject, "", pmsg.Concepts)
	if err != nil {
		p.logger.Error("storage error", "error", err)
	}

	msg.Topic = ""

	return msg
}

func (p *Pipeline) Submit(ctx context.Context, data, language string) (string, error) {
	msg := NewOCRInputMessage(data, language)
	if err := p.stream.Publish(TopicOCRInput, Message{
		ID:        msg.ID,
		Topic:     TopicOCRInput,
		Value:     msg.Value,
		Timestamp: time.Now(),
	}); err != nil {
		return "", err
	}

	return msg.ID, nil
}

func (p *Pipeline) Stop() {
	p.stream.Stop()
	p.logger.Info("pipeline stopped")
}
