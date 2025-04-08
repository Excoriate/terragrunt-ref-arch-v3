---
slug: /api/llm
---

# LLM Integration

Dagger's `LLM` core type includes API methods to attach objects to a Large Language Model (LLM), send prompts, and receive responses.

## Prompts

Use the `LLM.withPrompt()` API method to append prompts to the LLM context:

### System shell
```shell
dagger <<EOF
llm |
  with-prompt "What tools do you have available?"
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
llm |
  with-prompt "What tools do you have available?"
```

For longer or more complex prompts, use the `LLM.withPromptFile()` API method to read the prompt from a text file:

### System shell
```shell
dagger <<EOF
llm |
  with-prompt-file ./prompt.txt
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
llm |
  with-prompt-file ./prompt.txt
```


## Responses and Variables

Use the `LLM.lastReply()` API method to obtain the last reply from the LLM

Dagger supports the use of variables in prompts. This allows you to interpolate results of other operations into an LLM prompt:

### System shell
```shell
dagger <<EOF
source=\$(container |
  from alpine |
  with-directory /src https://github.com/dagger/dagger |
  directory /src)
environment=\$(env |
  with-directory-input 'source' \$source 'a directory with source code')
llm |
  with-env \$environment |
  with-prompt "The directory also has some tools available." |
  with-prompt "Use the tools in the directory to read the first paragraph of the README.md file in the directory." |
  with-prompt "Reply with only the selected text." |
  last-reply
EOF
```

### Dagger Shell
```shell title="First type 'dagger' for interactive mode."
source=$(container |
  from alpine |
  with-directory /src https://github.com/dagger/dagger |
  directory /src)
environment=$(env |
  with-directory-input 'source' $source 'a directory with source code')
llm |
  with-env $environment |
  with-prompt "The directory also has some tools available." |
  with-prompt "Use the tools in the directory to read the first paragraph of the README.md file in the directory." |
  with-prompt "Reply with only the selected text." |
  last-reply
```

> **Tip:**
> To get the complete message history, use the `LLM.History()` API method.

## Environments

Dagger [modules](../features/modules.md) are collections of Dagger Functions. When you give a Dagger module to the `LLM` core type, every Dagger Function is turned into a tool that the LLM can call.

Environments configure any number of inputs and outputs for the LLM. For example, an environment might provide a `Directory`, a `Container`, a custom module, and a `string` variable. The LLM can use the scalars and the functions of these objects to complete the assigned task.

The documentation for the modules are provided to the LLM, so make sure to provide helpful documentation in your Dagger Functions. The LLM should be able to figure out how to use the tools on its own. Don't worry about describing the objects too much in your prompts because it will be redundant with this automatic documentation.

Consider the following Dagger Function:

### Go
```go
package main

import (
	"context"

	"dagger.io/dagger"
)

type CodingAgent struct{}

// Write code to a file
func (m *CodingAgent) WriteCode(ctx context.Context, prompt string) (*dagger.Directory, error) {
	// Create a new directory for the toy workspace
	workspace := dag.Directory()

	// Create an environment with the toy workspace module
	env := dag.Env().
		WithInput("workspace", workspace, "A directory with source code").
		WithOutput("workspace", "The directory with the updated source code")

	// Create an LLM with the environment
	llm := dag.Llm().WithEnv(env)

	// Add prompts to the LLM
	llm = llm.
		WithPrompt("You are a helpful coding assistant.").
		WithPrompt("Use the tools available to you to write code that fulfills the user's request.").
		WithPrompt(prompt)

	// Get the last reply from the LLM
	reply := llm.LastReply()

	// Get the workspace directory from the reply
	return reply.Output("workspace").AsDirectory(), nil
}

```

### Python
```python
import dagger
from dagger import dag, function, object_type


@object_type
class CodingAgent:
    @function
    async def write_code(self, prompt: str) -> dagger.Directory:
        """Write code to a file"""
        # Create a new directory for the toy workspace
        workspace = dag.directory()

        # Create an environment with the toy workspace module
        env = (
            dag.env()
            .with_input("workspace", workspace, "A directory with source code")
            .with_output("workspace", "The directory with the updated source code")
        )

        # Create an LLM with the environment
        llm = dag.llm().with_env(env)

        # Add prompts to the LLM
        llm = (
            llm.with_prompt("You are a helpful coding assistant.")
            .with_prompt(
                "Use the tools available to you to write code that fulfills the user's request."
            )
            .with_prompt(prompt)
        )

        # Get the last reply from the LLM
        reply = llm.last_reply()

        # Get the workspace directory from the reply
        return reply.output("workspace").as_directory()

```

### TypeScript
```typescript
import { dag, Directory, func, object } from "@dagger.io/dagger"

@object()
class CodingAgent {
  /**
   * Write code to a file
   */
  @func()
  async writeCode(prompt: string): Promise<Directory> {
    // Create a new directory for the toy workspace
    const workspace = dag.directory()

    // Create an environment with the toy workspace module
    const env = dag
      .env()
      .withInput("workspace", workspace, "A directory with source code")
      .withOutput("workspace", "The directory with the updated source code")

    // Create an LLM with the environment
    let llm = dag.llm().withEnv(env)

    // Add prompts to the LLM
    llm = llm
      .withPrompt("You are a helpful coding assistant.")
      .withPrompt(
        "Use the tools available to you to write code that fulfills the user's request.",
      )
      .withPrompt(prompt)

    // Get the last reply from the LLM
    const reply = llm.lastReply()

    // Get the workspace directory from the reply
    return reply.output("workspace").asDirectory()
  }
}

```

### PHP
```php
<?php

declare(strict_types=1);

namespace DaggerModule;

use Dagger\Attribute\DaggerFunction;
use Dagger\Attribute\DaggerObject;
use Dagger\Client\Directory;

use function Dagger\dag;

#[DaggerObject]
class CodingAgent
{
    /**
     * Write code to a file
     */
    #[DaggerFunction]
    public function writeCode(string $prompt): Directory
    {
        // Create a new directory for the toy workspace
        $workspace = dag()->directory();

        // Create an environment with the toy workspace module
        $env = dag()
            ->env()
            ->withInput('workspace', $workspace, 'A directory with source code')
            ->withOutput('workspace', 'The directory with the updated source code');

        // Create an LLM with the environment
        $llm = dag()->llm()->withEnv($env);

        // Add prompts to the LLM
        $llm = $llm
            ->withPrompt('You are a helpful coding assistant.')
            ->withPrompt('Use the tools available to you to write code that fulfills the user\'s request.')
            ->withPrompt($prompt);

        // Get the last reply from the LLM
        $reply = $llm->lastReply();

        // Get the workspace directory from the reply
        return $reply->output('workspace')->asDirectory();
    }
}

```

### Java
```java
package io.dagger.modules.codingagent;

import io.dagger.client.Client;
import io.dagger.client.Dagger;
import io.dagger.client.Directory;
import io.dagger.client.Env;
import io.dagger.client.Llm;
import io.dagger.client.LlmReply;
import io.dagger.module.annotation.Module;
import io.dagger.module.annotation.Object;
import io.dagger.module.annotation.Function;

@Module
@Object
public class CodingAgent {

  /**
   * Write code to a file
   */
  @Function
  public Directory writeCode(String prompt) throws Exception {
    try (Client client = Dagger.connect()) {
      // Create a new directory for the toy workspace
      Directory workspace = client.directory();

      // Create an environment with the toy workspace module
      Env env = client
          .env()
          .withInput("workspace", workspace, "A directory with source code")
          .withOutput("workspace", "The directory with the updated source code");

      // Create an LLM with the environment
      Llm llm = client.llm().withEnv(env);

      // Add prompts to the LLM
      llm = llm
          .withPrompt("You are a helpful coding assistant.")
          .withPrompt(
              "Use the tools available to you to write code that fulfills the user's request.")
          .withPrompt(prompt);

      // Get the last reply from the LLM
      LlmReply reply = llm.lastReply();

      // Get the workspace directory from the reply
      return reply.output("workspace").asDirectory();
    }
  }
}

```

Here, an instance of the `ToyWorkspace` module is attached as an input to the `Env` environment. The `ToyWorkspace` module contains a number of Dagger Functions for developing code: `Read()`, `Write()`, and `Build()`. When this environment is attached to an `LLM`, the LLM can call any of these Dagger Functions to change the state of the `ToyWorkspace` and complete the assigned task.

In the `Env`, a `ToyWorkspace` instance called `after` is specified as a desired output of the LLM. This means that the LLM should return the `ToyWorkspace` module instance as a result of completing its task. The resulting `ToyWorkspace` object is then available for further processing or for use in other Dagger Functions.