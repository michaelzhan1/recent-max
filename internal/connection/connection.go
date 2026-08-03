package connection

import (
	"encoding/json"
	"net"
)

// TCPListener is a wrapper around net.Listener that provides convenience methods
type TCPListener struct {
	ln net.Listener
}

// NewTCPListener creates a new TCPListener on the given port
func NewTCPListener(port string) (*TCPListener, error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	return &TCPListener{ln: ln}, nil
}

// Accept waits for and returns the next connection to the listener
func (c *TCPListener) Accept() (net.Conn, error) {
	return c.ln.Accept()
}

// Handle runs a given function on a decoded connection
func (c *TCPListener) Handle(conn net.Conn, handler func(*json.Decoder) error) error {
	dec := json.NewDecoder(conn)
	return handler(dec)
}

// Close closes the TCPListener
func (c *TCPListener) Close() error {
	return c.ln.Close()
}

// TCPDialer is a wrapper around net.Conn that provides convenience methods
type TCPDialer struct {
	conn net.Conn
}

// NewTCPDialer creates a new TCPDialer that connects to the given address
func NewTCPDialer(address string) (*TCPDialer, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &TCPDialer{conn: conn}, nil
}

// Receive reads from the connection and decodes the JSON into the provided variable
func (d *TCPDialer) Send(v any) error {
	enc := json.NewEncoder(d.conn)
	return enc.Encode(v)
}

// Close closes the TCPDialer connection
func (d *TCPDialer) Close() error {
	return d.conn.Close()
}
