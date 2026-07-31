package connection

import (
	"encoding/json"
	"net"
)

type TCPListener struct {
	ln net.Listener
}

func NewTCPListener(port string) (*TCPListener, error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	return &TCPListener{ln: ln}, nil
}

func (c *TCPListener) AcceptAndHandle(handler func(*json.Decoder) error) error {
	conn, err := c.ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	dec := json.NewDecoder(conn)
	return handler(dec)
}

func (c *TCPListener) Close() error {
	return c.ln.Close()
}

type TCPDialer struct {
	conn net.Conn
}

func NewTCPDialer(address string) (*TCPDialer, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TCPDialer{conn: conn}, nil
}

func (d *TCPDialer) Send(v any) error {
	enc := json.NewEncoder(d.conn)
	return enc.Encode(v)
}

func (d *TCPDialer) Close() error {
	return d.conn.Close()
}
