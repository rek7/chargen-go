package chargen

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sort"
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
func NewClient(protocol, target string) (*Client, error) {
	var conn net.Conn
	if protocol != "udp" && protocol != "tcp" {
		return nil, fmt.Errorf("protocol needs to be either tcp or udp")
	}

	// 0 -> ip;1->port
	ip, portStr, err := net.SplitHostPort(target)
	if err != nil {
		return nil, err
	}

	if ip == "" || portStr == "" {
		return nil, fmt.Errorf("port or ip is empty, both need to be present")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("issue converting port %v", err)
	}

	targetip, err := net.LookupIP(ip)
	if err != nil {
		return nil, err
	}

	log.Printf("dialing host/port: %v:%v proto: %v\n", targetip[0].String(), port, protocol)

	l := make([]gopacket.SerializableLayer, 0)
	l = append(l, &layers.Ethernet{})

	if protocol == "tcp" {
		tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", targetip[0].String(), port))
		if err != nil {
			return nil, err
		}

		tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			return nil, err
		}

		conn = tcpConn
		l = append(l, &layers.TCP{
			SrcPort: layers.TCPPort(rand.Intn(65535)),
			DstPort: layers.TCPPort(port),
		})
	} else if protocol == "udp" {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", targetip[0].String(), port))
		if err != nil {
			return nil, err
		}

		udpConn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return nil, err
		}

		conn = udpConn
		l = append(l, &layers.UDP{
			SrcPort: layers.UDPPort(rand.Intn(65535)),
			DstPort: layers.UDPPort(port),
		})
	}

	return &Client{
		conn:   conn,
		port:   port,
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
	if newSrcInfo.String() == "<nil>" {
		if err := faker.FakeData(&a); err != nil {
			return err
		}
	}

	// check if field already exists
	getLayerIndex := func(layerNum int64) float32 {
		for i, layer := range c.layers {
			if int64(layer.LayerType()) == layerNum {
				return float32(i)
			}
		}

		return -1
	}

	ip, _, err := net.SplitHostPort(c.conn.RemoteAddr().String())
	if err != nil {
		return err
	}

	fmt.Println(ip)
	if strings.Contains(ip, ":") {
		l := &layers.IPv6{
			DstIP: net.IP(ip),
			SrcIP: newSrcInfo,
		}

		if newSrcInfo.String() == "<nil>" {
			l.SrcIP = net.IP(a.IPV6)
		}
		log.Printf("spoofing ip6 as: %v %v\n", a.IPV6, newSrcInfo.To16().String())

		if index := getLayerIndex(int64(l.LayerType())); index != -1 {
			c.layers[int(index)] = l
		} else {
			c.layers = append(c.layers, l)
		}
	} else {
		l := &layers.IPv4{
			DstIP: net.IP(ip),
			SrcIP: newSrcInfo,
		}

		if newSrcInfo.String() == "<nil>" {
			l.SrcIP = net.IP(a.IPV4)
		}
		log.Printf("spoofing ip4 as: %v %v\n", a.IPV4, newSrcInfo.String())

		if index := getLayerIndex(int64(l.LayerType())); index != -1 {
			c.layers[int(index)] = l
		} else {
			c.layers = append(c.layers, l)
		}
	}
	return nil
}

func (c *Client) order() {
	sort.Slice(c.layers[:], func(i, j int) bool {
		return int64(c.layers[i].LayerType()) < int64(c.layers[j].LayerType())
	})
}

func (c *Client) Write(numBytes int) error {
	c.order()
	payload := genData(numBytes)
	buf := gopacket.NewSerializeBuffer()

	n := make([]gopacket.SerializableLayer, 0)
	n = append(n, c.layers...)
	n = append(n, gopacket.Payload(payload))
	gopacket.SerializeLayers(buf,
		gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
		n...,
	)

	log.Printf("sending packet size: %v\n", len(payload))
	if _, err := c.conn.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// Read reads characters from the chargen server.
func (c *Client) Read() ([]byte, error) {
	br := bufio.NewReader(c.conn)
	line, _, err := br.ReadLine()
	if err != nil {
		return nil, err
	}
	return line, nil
}

// Close closes the connection to the chargen server.
func (c *Client) Close() error {
	return c.conn.Close()
}
