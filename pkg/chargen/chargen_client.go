package chargen

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/bxcodec/faker"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Client represents a chargen client.
type Client struct {
	conn   net.Conn
	port   int
	layers []gopacket.SerializableLayer
}

// NewClient creates a new chargen client.
func NewClient(target, protocol string) (*Client, error) {
	// 0 -> ip;1->port
	ip, port, isFound := strings.Cut(target, ":")
	if !isFound {
		return nil, fmt.Errorf("target needs to be ip:port")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("issue converting port %v", err)
	}

	if protocol != "udp" && protocol != "tcp" {
		return nil, fmt.Errorf("protocol needs to be either tcp or udp")
	}

	targetip, err := net.LookupIP(ip)
	if err != nil {
		return nil, err
	}

	prefix := "ip4"
	if strings.Contains(targetip[0].String(), ":") {
		prefix = "ip6"
	}

	conn, err := net.Dial(fmt.Sprintf("%v:%v", prefix, protocol), targetip[0].String())
	if err != nil {
		return nil, err
	}

	l := make([]gopacket.SerializableLayer, 0)
	l = append(l, &layers.Ethernet{})

	if protocol == "tcp" {
		l = append(l, &layers.TCP{
			SrcPort: layers.TCPPort(rand.Intn(65535)),
			DstPort: layers.TCPPort(portInt),
			Seq:     0,
			Window:  65535,
		})
	} else if protocol == "udp" {
		l = append(l, &layers.UDP{
			SrcPort: layers.UDPPort(rand.Intn(65535)),
			DstPort: layers.UDPPort(portInt),
		})
	}

	return &Client{
		conn:   conn,
		port:   portInt,
		layers: l,
	}, nil
}

func (c *Client) UpdateSrcIP(newSrcInfo net.IP) error {
	// generate a random public IP address
	type IPs struct {
		IPV4 string `faker:"ipv4"`
		IPV6 string `faker:"ipv6"`
	}
	a := IPs{}
	if newSrcInfo.String() == "" {
		err := faker.FakeData(&a)
		if err != nil {
			return err
		}
	}

	ip, _, _ := strings.Cut(c.conn.RemoteAddr().String(), ":")
	// Ipv6 Address Detected
	if strings.Contains(newSrcInfo.String(), ":") {
		l := &layers.IPv6{
			DstIP: net.IP(ip),
			SrcIP: newSrcInfo,
		}

		if newSrcInfo.String() == "" {
			l.SrcIP = net.IP(a.IPV6)
		}

		c.layers = append(c.layers, l)
	} else {
		l := &layers.IPv4{
			DstIP: net.IP(ip),
			SrcIP: newSrcInfo,
		}

		if newSrcInfo.String() == "" {
			l.SrcIP = net.IP(a.IPV4)
		}

		c.layers = append(c.layers, l)
	}
	return nil
}

func (c *Client) genData(num int) []byte {
	if num == 0 {
		num = rand.Intn(512-1) + 1
	}
	b := new(bytes.Buffer)
	for i := 0; num >= i; i++ {
		b.Write([]byte(fmt.Sprintf("%c", rand.Intn(126-33)+3)))
	}
	return b.Bytes()
}

func (c *Client) Write(numBytes int) error {
	payload := c.genData(numBytes)
	buf := gopacket.NewSerializeBuffer()

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	n := append(make([]gopacket.SerializableLayer, 0), c.layers...)
	fmt.Println(n, c.layers)
	//n = append(n, c.layers...)
	n = append(n, gopacket.Payload(payload))
	fmt.Println(n, c.layers)
	gopacket.SerializeLayers(buf, opts,
		n...,
	)

	_, err := c.conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// Read reads characters from the chargen server.
func (c *Client) Read() (string, error) {
	br := bufio.NewReader(c.conn)
	line, _, err := br.ReadLine()
	if err != nil {
		return "", err
	}
	return string(line), nil
}

// Close closes the connection to the chargen server.
func (c *Client) Close() error {
	return c.conn.Close()
}
