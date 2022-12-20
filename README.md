# Chargen-Go
Chargen protocol implementation as per [RFC 864](https://www.rfc-editor.org/rfc/rfc864), [RFC 865](https://www.rfc-editor.org/rfc/rfc865), and [RFC 866](https://www.rfc-editor.org/rfc/rfc866). Ipv6 friendly.

Client also has ability to change dst ip, for chargen [DDoS amplication attacks](https://www.link11.com/en/blog/threat-landscape/chargen-flood-attacks-explained/).

A Cli tool [main.go](./main.go), provides the ability to stand up a server and use a client.


## Installation
`go get -u github.com/rek7/chargen-go`


## Server Example

```go
import (
	"net"

	"github.com/rek7/chargen-go/pkg/chargen"
)

func main() {
	ln, err := net.Listen("tcp", ":19")
	if err != nil {
		// handle error
	}
	server := chargen.NewServer(ln)
	server.Serve()
}
```

## Client Example

```go
package main

import (
	"fmt"
    "net"

	"github.com/rek7/chargen-go/pkg/chargen"
)

func main() {
	client, err := chargen.NewClient("127.0.0.1:19", "tcp")
	if err != nil {
		// handle error
	}
	defer client.Close()

    // Used to spoof Source IP for amplification attacks. Leave a blank IP if you want it to randomly
    // generate a public ip.
    client.UpdateSrcIP(net.IP("192.168.1.1"))

	for {
		line, err := client.Read()
		if err != nil {
			// handle error
		}
		fmt.Println(line)
	}
}
```