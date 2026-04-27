package myreco

import (
	maa "github.com/MaaXYZ/maa-framework-go/v4"
	"github.com/rs/zerolog/log"
)

// MyCustomRecognition is a sample custom recognition that mirrors the Python agent's my_reco_222.
type MyCustomRecognition struct{}

// Compile-time interface check
var _ maa.CustomRecognitionRunner = &MyCustomRecognition{}

func (r *MyCustomRecognition) Run(ctx *maa.Context, arg *maa.CustomRecognitionArg) (*maa.CustomRecognitionResult, bool) {
	log.Info().
		Str("task", arg.CurrentTaskName).
		Str("recognition", arg.CustomRecognitionName).
		Msg("Custom recognition is running")

	// Example: override pipeline and run a sub-recognition
	ctx.OverridePipeline(map[string]any{
		"MyCustomOCR": map[string]any{
			"roi": []int{100, 100, 200, 300},
		},
	})

	return &maa.CustomRecognitionResult{
		Detail: "Hello from Go!",
	}, true
}
