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
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/ChrIgiSta/go-utils/connection"
	"github.com/ChrIgiSta/go-utils/containers"
	log "github.com/ChrIgiSta/go-utils/logger"
)

type Server struct {
	wg          sync.WaitGroup
	host        string
	port        uint16
	interrupted bool
	listener    net.Listener
	udpListener net.PacketConn
	handler     connection.Handler
	clients     *containers.List
	proto       connection.Protocol
	tlsConfig   *tls.Config
}

func NewServer(host string, port uint16, handler connection.Handler, protocol connection.Protocol) *Server {

	return &Server{
		host:        host,
		port:        port,
		wg:          sync.WaitGroup{},
		interrupted: false,
		listener:    nil,
		udpListener: nil,
		handler:     handler,
		clients:     containers.NewList(),
		proto:       protocol,
	}
}

func NewUdpServer(host string, port uint16, handler connection.Handler) *Server {
	return NewServer(host, port, handler, connection.Udp)
}

func NewTcpServer(host string, port uint16, handler connection.Handler) *Server {
	return NewServer(host, port, handler, connection.Tcp)
}

func NewUnixServer(host string, port uint16, handler connection.Handler) *Server {
	return NewServer(host, port, handler, connection.Unix)
}

func NewTlsServer(host string, port uint16, handle connection.Handler, certificate []byte, privateKey []byte) *Server {

	server := NewServer(host, port, handle, connection.Tls)

	if err := server.TlsConfig(certificate, privateKey); err != nil {
		_ = log.Error("socket", "setup tls listener: %v", err)
	}

	return server
}

func (s *Server) TlsConfig(certificate []byte, privateKey []byte) (err error) {
	var keyPair tls.Certificate

	keyPair, err = tls.X509KeyPair(certificate, privateKey)
	s.tlsConfig = &tls.Config{Certificates: []tls.Certificate{keyPair}}

	return
}

func (s *Server) ListenAndServe() (err error) {

	s.interrupted = false

	if s.listener != nil {
		return errors.New("listener already up")
	}

	sAddr := connection.Address(s.host, s.port)

	switch s.proto {
	case connection.Tcp:
		s.listener, err = net.Listen(string(s.proto),
			sAddr)
	case connection.Tls:
		s.listener, err = tls.Listen(string(connection.Tcp), sAddr, s.tlsConfig)
	case connection.Udp:
		var udpAddr *net.UDPAddr

		udpAddr, err = connection.GetUdpAddress(s.host, s.port)
		if err == nil {
			s.udpListener, err = net.ListenUDP(string(s.proto), udpAddr)
		}
	case connection.Unix:
		var lAddr *net.UnixAddr

		lAddr, err = connection.GetUnixAddress(s.host, s.port)
		if err != nil {
			return err
		}
		s.listener, err = net.ListenUnix(string(s.proto), lAddr)
	default:
		err = fmt.Errorf("unknown protocol: %s", s.proto)
	}

	if err != nil {
		return err
	}

	_ = log.Debug("socket", "server listen on %s:%d", s.host, s.port)

	s.wg.Add(1)
	switch s.proto {
	case connection.Tcp, connection.Tls, connection.Unix:
		go s.listenTcp(&s.wg)
	case connection.Udp:
		go s.listenUdp(&s.wg)
	}

	return
}

func (s *Server) Send(id int, msg []byte) (err error) {

	_, item := s.clients.Get(id)
	if item == nil {
		return errors.New("no connection with given id")
	}

	switch s.proto {
	case connection.Tcp, connection.Tls, connection.Unix:
		conn := item.(net.Conn)

		_, err = conn.Write(connection.AppendDelimeter(msg))

	case connection.Udp:
		addr := item.(net.Addr)

		_, err = s.udpListener.WriteTo(connection.AppendDelimeter(msg), addr)
		if err != nil {
			s.clients.Delete(id)
			s.handler.Disconnected(id)
		}
	}

	return err
}

func (s *Server) Broadcast(msg []byte) (err error) {
	cIds := s.clients.GetIds()

	for _, id := range cIds {
		err = s.Send(id, msg)
		if err != nil {
			return
		}
	}

	return
}

func (s *Server) ClientIp(id int) (ip string, err error) {

	_, item := s.clients.Get(id)
	if item == nil {
		return "", errors.New("no connection with given id")
	}

	switch s.proto {
	case connection.Tcp, connection.Tls, connection.Unix:
		conn := item.(net.Conn)
		ip = conn.RemoteAddr().String()

	case connection.Udp:
		addr := item.(net.Addr)
		ip = addr.String()
	}
	return
}

func (s *Server) Stop() {
	defer s.wg.Wait()

	s.interrupted = true

	if s.proto == connection.Tcp || s.proto == connection.Tls || s.proto == connection.Unix {
		s.listener.Close()

		cIds := s.clients.GetIds()
		for _, id := range cIds {
			conn, err := s.getConnFromId(id)
			if err != nil {
				conn.Close()
			}
		}
	}

	if s.proto == connection.Udp {
		s.udpListener.Close()
	}

	s.clients.Reset()
}

func (s *Server) listenTcp(wg *sync.WaitGroup) {

	// defer s.listener.Close()
	defer wg.Done()

	for !s.interrupted {
		conn, err := s.listener.Accept()
		if err != nil {
			_ = log.Warn("socket", "accept client: %v", err)
			continue
		}

		_ = log.Debug("socket", "server accept client %v", conn.RemoteAddr())

		if !s.interrupted {
			id := connection.GetIdFromConn(&conn)
			s.clients.AddOrUpdate(id, conn)
			wg.Add(1)
			go s.clientHandler(wg, conn, id)
		}
	}
	_ = log.Debug("socket", "listener exited")
}

func (s *Server) listenUdp(wg *sync.WaitGroup) {

	defer wg.Done()
	// defer s.udpListener.Close()

	var buffer []byte = make([]byte, 1024)

	for !s.interrupted {
		n, addr, err := s.udpListener.ReadFrom(buffer)
		if err != nil {
			_ = log.Error("socket", "read udp: %v", err)
			return
		}

		id := connection.GetIdFromAddr(&addr)
		s.clients.AddOrUpdate(id, addr)
		s.handler.Connected(id)
		s.handler.Received(id, buffer[:n-1]) // delim ?
		_ = log.Fine("socket", "udp pck from %v", addr.String())
	}
	_ = log.Debug("socket", "listener exited")
}

func (s *Server) clientHandler(wg *sync.WaitGroup,
	client net.Conn, id int) {

	defer wg.Done()
	defer client.Close()

	bufReader := bufio.NewReader(client)

	s.handler.Connected(id)
	defer s.handler.Disconnected(id)

	for !s.interrupted {
		msg, err := bufReader.ReadBytes(connection.DefaultDelimiter)

		if err != nil {
			_ = log.Error("socket", "read from client: %v", err)
			return
		} else {
			msg = connection.TruncateDelimeter(msg)
			_ = log.Fine("socket", "server: rx %v", string(msg[:len(msg)-1]))
			s.handler.Received(id, msg)
		}
	}
}

func (s *Server) getConnFromId(id int) (conn net.Conn, err error) {
	_, connIf := s.clients.Get(id)
	if connIf == nil {
		return nil, errors.New("not found")
	}
	return connIf.(net.Conn), err
}
