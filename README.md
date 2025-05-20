![release](https://img.shields.io/github/v/release/dskart/gollum)

# 🧌 GoLLuM

GoLLuM is a simple LLM library that can be used to build LLM applications.

GoLLuM is constructed of 3 main modules. Using the modules together allow you to build complex LLM Agent workflows.

You can find examples of how to use the modules in the [examples](./examples) directory.

## OpenAi

The [openai](./openai) module is a simple openai client that can be used to interact with the openai API. 


## Scrolls

The [scrolls](./scrolls) module is a simple templating engine that can be used to build prompts. You can use Scrolls with the OpenAi module run your prompts against the OpenAi API.


## Ringchain

The [ringchain](./ringchain) module is a simple graph library that can be used to build and run Agent Graphs concurrently. You can use the Scrolls module to manage your Agent and tool prompts.
