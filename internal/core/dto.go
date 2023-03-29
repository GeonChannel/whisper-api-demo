package core

import (
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"strings"
)

type ResultDTO struct {
	Segments []*SegmentDTO `json:"segments"`
	FullText string        `json:"fullText"`
}

type SegmentDTO struct {
	// Segment Number
	Num int `json:"num"`

	// Time beginning and end timestamps for the segment.
	Start int64 `json:"start"`
	End   int64 `json:"end"`

	// The text of the segment.
	Text string `json:"text"`
}
type ReqS3DTO struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func segmentToDTO(segment whisper.Segment) *SegmentDTO {
	return &SegmentDTO{
		Num:   segment.Num,
		Start: segment.Start.Milliseconds(),
		End:   segment.End.Milliseconds(),
		Text:  segment.Text,
	}
}

func segmentsToDTO(segments []whisper.Segment) []*SegmentDTO {
	res := make([]*SegmentDTO, 0)
	for _, segment := range segments {
		res = append(res, segmentToDTO(segment))
	}
	return res
}

func segmentsToFullText(segments []whisper.Segment) string {
	var builder strings.Builder
	for _, segmentText := range segments {
		builder.WriteString(" " + segmentText.Text)
	}
	return builder.String()
}
