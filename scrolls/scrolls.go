package scrolls

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/dskart/gollum/openai"
)

type Scroll struct {
	openAi openai.OpenAi

	text    string
	funcMap template.FuncMap
	mu      sync.RWMutex
}

type Options struct {
	funcMap template.FuncMap
}

func WithFuncMap(funcMap template.FuncMap) func(*Options) {
	return func(opts *Options) {
		opts.funcMap = funcMap
	}
}

func New(text string, openAi openai.OpenAi, opts ...func(*Options)) *Scroll {
	options := Options{}
	for _, option := range opts {
		option(&options)
	}

	return &Scroll{
		text:    text,
		openAi:  openAi,
		funcMap: options.funcMap,
	}
}

var re = regexp.MustCompile(`(?s)\[\[#(system|user|assistant)~\]\](.*?)\[\[~/(system|user|assistant)\]\]`)

func (s *Scroll) ParseBlocks(args map[string]any) ([]openai.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tmpl, err := template.New("scroll").Funcs(s.funcMap).Parse(s.text)
	if err != nil {
		return nil, fmt.Errorf("could not parse template text: %w", err)
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, args)
	if err != nil {
		return nil, fmt.Errorf("could not execute template: %w", err)
	}

	parsedText := buf.String()
	matches := re.FindAllStringSubmatch(parsedText, -1)
	if len(matches) < 1 {
		return nil, fmt.Errorf("could not find any message container matches")
	}

	msgs := make([]openai.Message, 0, len(matches))
	for _, match := range matches {
		role := match[1]    // The matched role (system, user, or assistant)
		content := match[2] // The content between {{#...~}} and {{~/...}}
		content = strings.Trim(content, "\n")

		msgs = append(msgs, openai.Message{
			Role:    openai.RoleType(role),
			Content: &content,
		})
	}

	return msgs, nil
}

// Execute executes the scroll template and returns all the execute openai.Messages
func (s *Scroll) Execute(ctx context.Context, args map[string]any) ([]openai.Message, map[string]string, error) {
	blocks, err := s.ParseBlocks(args)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse scroll: %w", err)
	}

	genArgs := maps.Clone(args)
	genOutputs := make(map[string]string)
	outputMsgs := make([]openai.Message, 0, len(blocks))
	for _, msg := range blocks {
		newHistoryMsg := openai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}

		if msg.Role == "assistant" {
			assistantAction := true
			var assistantBody AssistantBody
			if err := json.Unmarshal([]byte(*msg.Content), &assistantBody); err != nil {
				// if we can't unmarshal the assistant body, then it's not an assistant action
				assistantAction = false
			}

			if assistantAction {
				resp, err := promptOpenAi(ctx, s.openAi, outputMsgs, assistantBody.OpenAiPromptOptions()...)
				if err != nil {
					return nil, genOutputs, fmt.Errorf("failed to prompt openai: %w", err)
				}
				genArgs[assistantBody.OutputName] = resp
				genOutputs[assistantBody.OutputName] = resp
				newHistoryMsg.Content = &resp
			}
		}
		outputMsgs = append(outputMsgs, newHistoryMsg)
	}

	return outputMsgs, genOutputs, nil
}
