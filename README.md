# (Previous content unchanged)

## Usage
Once you have it all set up, here's how to use it.

```bash
fabric -h
```

```bash
usage: fabric -h
Usage:
  fabric [OPTIONS]

Application Options:
  -p, --pattern=                    Choose a pattern
  -v, --variable=                   Values for pattern variables, e.g. -v=$name:John -v=$age:30
  -C, --context=                    Choose a context
      --session=                    Choose a session
  -S, --setup                       Run setup
      --setup-skip-update-patterns  Skip update patterns at setup
  -t, --temperature=                Set temperature (default: 0.7)
  -T, --topp=                       Set top P (default: 0.9)
  -s, --stream                      Stream
  -P, --presencepenalty=            Set presence penalty (default: 0.0)
  -F, --frequencypenalty=           Set frequency penalty (default: 0.0)
  -l, --listpatterns                List all patterns
  -L, --listmodels                  List all available models
  -x, --listcontexts                List all contexts
  -X, --listsessions                List all sessions
  -U, --updatepatterns              Update patterns
  -c, --copy                        Copy to clipboard
  -m, --model=                      Choose model
  -o, --output=                     Output to file
  -n, --latest=                     Number of latest patterns to list (default: 0)
  -d, --changeDefaultModel          Change default pattern
  -y, --youtube=                    YouTube video url to grab transcript, comments from it and send to chat
      --transcript                  Grab transcript from YouTube video and send to chat
      --comments                    Grab comments from YouTube video and send to chat
      --dry-run                     Show what would be sent to the model without actually sending it

Help Options:
  -h, --help                        Show this help message

```

### Supported Models and Integrations

Fabric supports various AI models and integrations, including:

- OpenAI models (GPT-3.5, GPT-4, etc.)
- Anthropic models (Claude, etc.)
- Local models through LM Studio integration

#### LM Studio Integration

Fabric now supports integration with LM Studio, allowing you to use local language models. This integration provides the following features:

- Listing available models
- Sending chat completions
- Generating text completions
- Creating embeddings

To use LM Studio with Fabric, ensure you have LM Studio installed and running on your local machine. Then, use the `-m lmstudio` flag when running Fabric commands:

```bash
fabric -m lmstudio -p summarize "Your text to summarize here"
```

For detailed information on setting up and using the LM Studio integration, please refer to the [LM Studio Integration Guide](vendors/lmstudio/README.md).

For detailed information on using specific models or integrations, please refer to the documentation in the respective vendor folders.

# (Rest of the content unchanged)
# fabricc
