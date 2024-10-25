package openai

type Config struct {
	OpenAiKey string `yaml:"OpenAiKey"`
	GptModel  string `yaml:"GptModel"`
}

func (c Config) Validate() error {
	if c.OpenAiKey == "" {
		return ErrOpenAiKeyNotSet
	}

	if c.GptModel == "" {
		return ErrGptModelNotSet
	}

	return nil
}
