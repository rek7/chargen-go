package chargen

import (
	"bufio"
	"fmt"
	"net"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

// ServeTCP serves chargen TCP requests.
func (s *Server) ServeTCP(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) error {
			bw := bufio.NewWriter(conn)
			for {
				// Generate a stream of characters starting with ASCII character 32 (' ')
				// and ending with ASCII character 126 ('~').
				for i := 32; i <= 126; i++ {
					if _, err := fmt.Fprintf(bw, "%c", i); err != nil {
						return err
					}
				}
				bw.Flush()
			}
		}(conn)
	}
}

// Serve serves chargen UDP requests.
func (s *Server) ServeUDP(ln *net.UDPConn) error {
	p := make([]byte, 2048)
	for {
		_, remoteaddr, err := ln.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go func(conn *net.UDPConn, addr *net.UDPAddr) error {
			_, err := conn.WriteToUDP(genData(0), addr)
			if err != nil {
				return err
			}
			return nil
		}(ln, remoteaddr)
	}
}
