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

package connection

import (
	"fmt"
	"net"
	"unsafe"

	log "github.com/ChrIgiSta/go-utils/logger"
)

const DefaultDelimiter = 0x00

type Protocol string

const (
	Tcp  Protocol = "tcp"
	Udp  Protocol = "udp"
	Tls  Protocol = "tls"
	Unix Protocol = "unix"
)

type Handler interface {
	Connected(id int)
	Received(id int, message []byte)
	Disconnected(id int)
}

type Message struct {
	Id      int
	Content []byte
}

type EventType int

const (
	DISCONNECTED EventType = 0
	CONNECTED    EventType = 1
	ERROR        EventType = -1
)

func GetIdFromConn(conn *net.Conn) int {
	return int(uintptr(unsafe.Pointer(conn)))
}

func GetIdFromAddr(addr *net.Addr) int {
	return int(uintptr(unsafe.Pointer(addr)))
}

func AppendDelimeter(raw []byte) []byte {
	return append(raw, DefaultDelimiter)
}

func TruncateDelimeter(raw []byte) []byte {
	return raw[:len(raw)-1]
}

func GetTcpAddress(host string,
	port uint16) (addr *net.TCPAddr, err error) {

	return net.ResolveTCPAddr(string(Tcp),
		Address(host, port))
}

func GetUdpAddress(host string,
	port uint16) (addr *net.UDPAddr, err error) {

	return net.ResolveUDPAddr(string(Udp),
		Address(host, port))
}

func GetUnixAddress(host string,
	port uint16) (addr *net.UnixAddr, err error) {

	return net.ResolveUnixAddr(string(Unix),
		Address(host, port))
}

func Address(host string, port uint16) string {
	return fmt.Sprintf("%s:%d", host, port)
}

type Event struct {
	Id int
	EventType
}

type EventsToChannel struct {
	messageChannel chan<- Message
	eventChannel   chan<- Event
}

func NewEventsToChannel(messageChannel chan<- Message,
	eventChannel chan<- Event) *EventsToChannel {

	return &EventsToChannel{
		messageChannel: messageChannel,
		eventChannel:   eventChannel,
	}
}

func (e2c *EventsToChannel) Connected(id int) {
	_ = log.Fine("evt2ch", "connected called with id %d", id)

	if e2c.eventChannel != nil {
		e2c.eventChannel <- Event{
			Id:        id,
			EventType: CONNECTED,
		}
	} else {
		_ = log.Warn("evt2ch", "event channel is nil")
	}
}

func (e2c *EventsToChannel) Received(id int, message []byte) {
	_ = log.Fine("evt2ch", "rx called with id %d, %v", id, string(message))

	if e2c.messageChannel != nil {
		e2c.messageChannel <- Message{
			Id:      id,
			Content: message,
		}
	} else {
		_ = log.Warn("evt2ch", "message channel is nil")
	}
}

func (e2c *EventsToChannel) Disconnected(id int) {
	_ = log.Fine("evt2ch", "disconnected called with id %d", id)

	if e2c.eventChannel != nil {
		e2c.eventChannel <- Event{
			Id:        id,
			EventType: DISCONNECTED,
		}
	} else {
		_ = log.Warn("evt2ch", "event channel is nil")
	}
}
