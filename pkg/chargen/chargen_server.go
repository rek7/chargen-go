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
		if s.ln.Addr().Network() == "tcp" {
			go s.handleConnectionTCP(conn)
		} else if s.ln.Addr().Network() == "udp" {

		}
	}
}

func (s *Server) handleConnectionTCP(conn net.Conn) error {
	defer conn.Close()
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
	return nil
}
