package core

import (
	"fmt"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/spf13/viper"
	"sync"
)

type Model struct {
	whisper.Model
}

func NewModel() *Model {
	modelName := viper.GetString("model") // "ggml-large.bin"
	model, err := whisper.New(viper.GetString("modelPath") + modelName)
	if err != nil {
		panic(err)
	}

	return &Model{
		Model: model,
	}
}

var mu sync.Mutex

func (m *Model) Transcribe(samples []float32) (*ResultDTO, error) {
	//mu.Lock()
	//defer mu.Unlock()
	context, err := m.NewContext()
	if err != nil {
		return nil, err
	}

	_ = context.SetLanguage("auto") // auto detecting : "auto"; korean : "ko"
	context.SetTranslate(false)
	context.SetThreads(8)

	if err := context.Process(samples, nil); err != nil {
		return nil, err
	}

	//todo: process 가 너무 오래걸리면?
	segments := make([]whisper.Segment, 0)
	// Print out the results
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		fmt.Printf("[%6s->%6s] %s\n", segment.Start, segment.End, segment.Text)
		segments = append(segments, segment)
	}

	return &ResultDTO{
		Segments: segmentsToDTO(segments),
		FullText: segmentsToFullText(segments),
	}, nil
}
