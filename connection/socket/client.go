/**
 * Copyright © 2024, Staufi Tech - Switzerland
 * All rights reserved.
 *
 *   ________________________   ___ _     ________________  _  ____
 *  / _____  _  ____________/  / __|_|   /_______________  | | ___/
 * ( (____ _| |_ _____ _   _ _| |__ _      | |_____  ____| |_|_
 *  \____ (_   _|____ | | | (_   __) |     | | ___ |/ ___)  _  \
 *  _____) )| |_/ ___ | |_| | | |  | |     | | ____( (___| | | |
 * (______/  \__)_____|____/  |_|  |_|     |_|_____)\____)_| |_|
 *
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

package socket

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"sync"

	"github.com/ChrIgiSta/go-utils/connection"
	log "github.com/ChrIgiSta/go-utils/logger"
)

type Client struct {
	host        string
	port        uint16
	conn        net.Conn
	connected   bool
	wg          sync.WaitGroup
	interrupted bool
	handler     connection.Handler
	proto       connection.Protocol
	tlsConfig   *tls.Config
}

func NewClient(host string, port uint16, handler connection.Handler, protocol connection.Protocol) *Client {

	return &Client{
		host:        host,
		port:        port,
		connected:   false,
		wg:          sync.WaitGroup{},
		interrupted: false,
		handler:     handler,
		proto:       protocol,
	}
}

func NewTcpClient(host string, port uint16, handler connection.Handler) *Client {
	return NewClient(host, port, handler, connection.Tcp)
}

func NewUdpClient(host string, port uint16, handler connection.Handler) *Client {
	return NewClient(host, port, handler, connection.Udp)
}

func NewTlsClient(host string, port uint16, handle connection.Handler, caCert []byte, verifyServer bool) *Client {
	client := NewClient(host, port, handle, connection.Tls)
	if err := client.TlsConfig(caCert, verifyServer); err != nil {
		_ = log.Error("socket", "setup tls client: %v", err)
	}

	return client
}

func NewUnixClient(path string, port uint16, handler connection.Handler) *Client {
	return NewClient(path, port, handler, connection.Unix)
}

func (c *Client) TlsConfig(caCertificates []byte, verifyCert bool) error {
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCertificates)
	if !ok {
		_ = log.Error("socket", "load ca cert: \r\n%v\r\n", string(caCertificates))
		return errors.New("cannot load ca certs")
	}
	c.tlsConfig = &tls.Config{RootCAs: caCertPool, InsecureSkipVerify: !verifyCert}

	return nil
}

func (c *Client) Connect() (err error) {
	c.interrupted = false

	switch c.proto {
	case connection.Tcp:
		c.conn, err = c.dailTcp()
	case connection.Udp:
		c.conn, err = c.dailUdp()
	case connection.Tls:
		c.conn, err = c.dailTls()
	case connection.Unix:
		c.conn, err = c.dailUnix()
	default:
		err = errors.New("unknown protocol")
	}

	if err != nil || c.conn == nil {
		if err == nil {
			err = errors.New("connection nil")
		}
		return err
	}

	c.wg.Add(1)
	c.connected = true
	go c.reader(&c.wg)

	return
}

func (c *Client) Send(msg []byte) error {
	if c.conn == nil {
		return errors.New("not connected")
	}

	_, err := c.conn.Write(connection.AppendDelimeter(msg))
	return err
}

func (c *Client) Disconnect() (err error) {
	defer c.wg.Wait()

	c.interrupted = true

	err = c.conn.Close()
	return
}

func (c *Client) reader(wg *sync.WaitGroup) {
	defer wg.Done()
	defer c.conn.Close()
	defer c.handler.Disconnected(1)

	c.handler.Connected(1)

	bufReader := bufio.NewReader(c.conn)

	_ = log.Debug("socket", "client: connected")

	for !c.interrupted {
		msg, err := bufReader.ReadBytes(connection.DefaultDelimiter)

		if err != nil {
			_ = log.Warn("socket", "client read: %v", err)
			break
		}
		msg = connection.TruncateDelimeter(msg)
		_ = log.Fine("socket", "client rx: %v", string(msg[:len(msg)-1]))
		c.handler.Received(1, msg)
	}

	_ = log.Debug("socket", "client: disconnected")

	c.connected = false
}

func (c *Client) dailTcp() (conn net.Conn, err error) {
	remoteAddr, err := connection.GetTcpAddress(c.host, c.port)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialTCP(string(c.proto), nil, remoteAddr)

	return
}

func (c *Client) dailTls() (conn net.Conn, err error) {
	remoteAddr, err := connection.GetTcpAddress(c.host, c.port)
	if err != nil {
		return nil, err
	}

	conn, err = tls.Dial(string(connection.Tcp), remoteAddr.String(), c.tlsConfig)

	return
}

func (c *Client) dailUnix() (conn net.Conn, err error) {
	remoteAddr, err := connection.GetUnixAddress(c.host, c.port)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialUnix(string(connection.Unix), nil, remoteAddr)

	return
}

func (c *Client) dailUdp() (conn net.Conn, err error) {
	remoteAddr, err := connection.GetUdpAddress(c.host, c.port)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialUDP(string(c.proto), nil, remoteAddr)
	return
}
