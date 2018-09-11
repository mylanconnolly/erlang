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
    body, readErr := p.ReadMsg()

    if readErr != nil && readErr != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading string %s\n", readErr.Error())
			os.Exit(1)
    }
    if err := p.WriteMsg([]byte("Responding to " + string(body))); err != nil {
			fmt.Fprintf(os.Stderr, "Error wririting strign %s\n", err.Error())
			os.Exit(1)
    }
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
