# eslgo

[![Go](https://github.com/GoingBytes/eslgo/actions/workflows/go.yml/badge.svg)](https://github.com/GoingBytes/eslgo/actions/workflows/go.yml)

`eslgo` is an idiomatic [FreeSWITCH™](https://freeswitch.com/) Event Socket Library for Go. It powers high-volume production systems and provides the primitives needed to originate calls, react to FreeSWITCH events, and build resilient telephony workflows from Go applications.

## Features
- Inbound ESL client with automatic reconnect hooks and context-aware commands
- Outbound ESL server for driving FreeSWITCH call-control scripts
- Event listener registration scoped by `UUID`, `Application-UUID`, `Job-UUID`, or catch-all handlers
- Typed command abstractions with support for custom command implementations via the `Command` interface (`BuildMessage() string`)
- Convenience helpers for DTMF, call origination, answer/hangup, and audio playback
- Context propagation for canceling long-running requests

## Installation

```bash
go get github.com/GoingBytes/eslgo
```

The module follows standard Go module semantics, so importing it in your project and executing `go mod tidy` is usually enough.

## Quick Start

### Outbound ESL server
Expose an outbound ESL handler that FreeSWITCH can call into:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/GoingBytes/eslgo"
)

func main() {
	// Listen blocks while FreeSWITCH connections are handled.
	log.Fatalln(eslgo.ListenAndServe(":8084", handleConnection))
}

func handleConnection(ctx context.Context, conn *eslgo.Conn, response *eslgo.RawResponse) {
	fmt.Printf("Got connection! %#v\n", response)

	// Originate a foreground (api) call to user 100 and play audio on the b-leg.
	response, err := conn.OriginateCall(
		ctx,
		false,
		eslgo.Leg{CallURL: "user/100"},
		eslgo.Leg{CallURL: "&playback(misc/ivr-to_hear_screaming_monkeys.wav)"},
		map[string]string{},
	)
	fmt.Println("Call originated:", response, err)
}
```

### Inbound ESL client
Connect to an existing FreeSWITCH instance and issue background jobs:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/GoingBytes/eslgo"
)

func main() {
	conn, err := eslgo.Dial("127.0.0.1:8021", "ClueCon", func() {
		fmt.Println("Inbound connection disconnected")
	})
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	response, err := conn.OriginateCall(
		ctx,
		true,
		eslgo.Leg{CallURL: "user/100"},
		eslgo.Leg{CallURL: "&playback(misc/ivr-to_hear_screaming_monkeys.wav)"},
		map[string]string{},
	)
	fmt.Println("Call originated:", response, err)

	time.Sleep(60 * time.Second)
	conn.ExitAndClose()
}
```

## Examples
- `example/outbound/`: runnable outbound server sample
- `example/inbound/`: minimal inbound client

Each example is self-contained and can be executed with `go run ./example/<name>`.

## Development
- Run the unit tests with `go test ./...`
- File issues or pull requests if you run into bugs or have feature suggestions

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).
