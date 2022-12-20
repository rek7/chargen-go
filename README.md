# Chargen-Go
Chargen protocol implementation as per [RFC 864](https://www.rfc-editor.org/rfc/rfc864), [RFC 865](https://www.rfc-editor.org/rfc/rfc865), and [RFC 866](https://www.rfc-editor.org/rfc/rfc866). Ipv6 friendly.

Client also has ability to spoof source IPs for chargen [DDoS amplication attacks](https://www.link11.com/en/blog/threat-landscape/chargen-flood-attacks-explained/).

A Cli tool [main.go](./main.go), provides the ability to stand up a server and use a client.


## Installation
`go get -u github.com/rek7/chargen-go`

Install cli tool:

`go install github.com/rek7/chargen-go`



## Cli Command
```
$ ~/go/bin/chargen-go
chargen - Wrapper for github.com/rek7/chargen-go library to stand up a server or run a client.

Usage:
  chargen [command]

Available Commands:
  client      Chargen UDP/TCP client
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  server      Chargen TCP/UDP client

Flags:
  -h, --help   help for chargen

Use "chargen [command] --help" for more information about a command.
```


## Server Example
### TCP
```go
package main

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
	server.ServeTCP()
}
```
### UDP
```go
package main

import (
    "net"

    "github.com/rek7/chargen-go/pkg/chargen"
)

func main() {
    ln, err := net.ListenUDP("udp", &net.UDPAddr{
        Port: 19,
        IP:   net.ParseIP("0.0.0.0"),
    })
    if err != nil {
        // handle error
    }
    server := chargen.NewServer(ln)
    server.ServeUDP()
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
    client, err := chargen.NewClient("tcp", "127.0.0.1:19")
    if err != nil {
        // handle error
    }
    defer client.Close()

    // Used to spoof Source IP for amplification attacks. Leave a blank IP if you want it to randomly
    // generate a public ip.
    client.UpdateSrcIP(net.IP("192.168.1.1"))

    // Specify 0 to generate a packet random within the size 1-512 bytes.
    client.Write(0)

    for {
        line, err := client.Read()
        if err != nil {
            // handle error
        }
        fmt.Println(string(line))
    }
}
```