package myaction

import (
	maa "github.com/MaaXYZ/maa-framework-go/v4"
	"github.com/rs/zerolog/log"
)

// MyCustomAction is a sample custom action that mirrors the Python agent's my_action_111.
type MyCustomAction struct{}

// Compile-time interface check
var _ maa.CustomActionRunner = &MyCustomAction{}

func (a *MyCustomAction) Run(ctx *maa.Context, arg *maa.CustomActionArg) bool {
	log.Info().
		Str("task", arg.CurrentTaskName).
		Str("action", arg.CustomActionName).
		Msg("Custom action is running")
	return true
}
