# Erlang

This is a Go package used to interface with Erlang. Currently only Ports are
supported.

I built this because I like using Elixir for large web applications, but in some
cases the performance can be an issue. I can circumvent this by using Go for
some more computationally-intense operations, and the `Port` system in Elixir
makes this pretty easy.

Although I haven't done so, it should also be possible to use this with an
Erlang system, as well.

## Usage

Get it:

```bash
$ go get github.com/mylanconnolly/erlang
```

Write some code:

```go
package main

import (
  "fmt"

  "github.com/mylanconnolly/erlang"
)

func main() {
	p := erlang.NewPort()

	for {
		// Receive a message from stdin, as well as any error that may be encountered.
		body, readErr := p.ReadMsg()

		// Check if there's an actual error (not EOF, we shouldn't handle that yet)
		if readErr != nil && readErr != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading string %s\n", readErr.Error())
			os.Exit(1)
		}
		// Do whatever you want to do with the message from stdin
		if _, err := p.Write([]byte("Responding to " + string(body))); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing string %s\n", err.Error())
			os.Exit(1)
		}
		// If we encountered an EOF, let's break the loop, since the program is
		// done.
		if readErr == io.EOF {
			break
		}
	}
}
```

Since you're passing in a byte slice, you could encode JSON, msgpack, or any
other format you want into it, and things should just work.

This is a very simple implementation, and it's likely that there are
improvements and optimizations to be made. I'd welcome any PRs you want to send!
