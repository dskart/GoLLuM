package openai

import "errors"

var ErrOpenAiKeyNotSet = errors.New("openai key not set")
var ErrGptModelNotSet = errors.New("gpt model not set")
