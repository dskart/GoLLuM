![release](https://img.shields.io/github/v/release/dskart/gollum)
[![Go Report Card](https://goreportcard.com/badge/github.com/dskart/gollum)](https://goreportcard.com/report/github.com/dskart/gollum)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# 🧌 oLLuM

GoLLuM is a powerful Go library for building LLM-powered applications and agent workflows. With a focus on simplicity and extensibility, GoLLuM enables developers to easily integrate LLM capabilities into their Go applications.

## Features

- **Simple API**: Clean, idiomatic Go interfaces make it easy to integrate LLMs into your applications
- **Modular Architecture**: Use the modules independently or together to build complex workflows
- **Concurrent Execution**: Execute agent workflows efficiently with built-in concurrency support
- **Extensible Design**: Easy to extend with custom tools and components

## Installation

```bash
go get github.com/dskart/gollum
```

## Modules

GoLLuM is constructed of 3 main modules that can be used together to build complex LLM agent workflows:

### OpenAI

The [openai](./openai) module provides a clean, type-safe client for interacting with the OpenAI API. It handles authentication, request formatting, and response parsing, making it easy to use OpenAI's models in your Go applications.

### Scrolls

The [scrolls](./scrolls) module is a templating engine for building and managing prompts. Scrolls makes it easy to create, maintain, and execute complex prompt templates with the OpenAI module.

### Ringchain

The [ringchain](./ringchain) module is a graph-based framework for building and running agent workflows concurrently. It allows you to define complex, multi-step workflows as directed acyclic graphs (DAGs) and execute them efficiently.

## Examples

You can find examples of how to use the modules in the [examples](./examples) directory:

- [Basic OpenAI integration](./examples/openai)
- [Multi-agent workflows with Ringchain](./examples/ringchain)
- [Prompt templating with Scrolls](./examples/scrolls)

## Getting Started

Check out the individual module READMEs for detailed usage instructions:

- [OpenAI Module](./openai/README.md)
- [Scrolls Module](./scrolls/README.md)
- [Ringchain Module](./ringchain/README.md)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
