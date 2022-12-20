package chargen

import (
	"bufio"
	"fmt"
	"net"
)

// Server represents a chargen server.
type Server struct {
	ln net.Listener
}

// NewServer creates a new chargen server.
func NewServer(ln net.Listener) *Server {
	return &Server{ln}
}

// Serve serves chargen requests.
func (s *Server) Serve() error {
	if s.ln.Addr().Network() != "tcp" || s.ln.Addr().Network() != "udp" {
		return fmt.Errorf("server protocol must be tcp or udp")
	}
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}

		go s.handleConnectionTCP(conn, s.ln.Addr().Network())
	}
}

func (s *Server) handleConnectionTCP(conn net.Conn, proto string) error {
	defer conn.Close()
	if proto == "tcp" {
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
	} else {
		buf := make([]byte, 1)
		if _, err := conn.Read(buf); err != nil {
			return err
		}
		// gen random data and send it over
		if _, err := conn.Write(genData(0)); err != nil {
			return err
		}
	}
	return nil
}
