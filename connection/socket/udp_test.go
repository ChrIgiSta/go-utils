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
	"fmt"
	"testing"
	"time"

	"github.com/ChrIgiSta/go-utils/connection"
)

func TestUdpClientServer(t *testing.T) {
	cMsgCh := make(chan connection.Message, 1)
	cEvtCh := make(chan connection.Event, 1)

	sMsgCh := make(chan connection.Message, 1)
	sEvtCh := make(chan connection.Event, 1)

	cCh := connection.NewEventsToChannel(cMsgCh, cEvtCh)
	c := NewUdpClient("localhost", 22333, cCh)

	sCh := connection.NewEventsToChannel(sMsgCh, sEvtCh)
	s := NewUdpServer("localhost", 22333, sCh)

	err := s.ListenAndServe()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = c.Connect()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	cEvt := <-cEvtCh
	if cEvt.EventType != connection.CONNECTED {
		t.Error("client event is not connected")
		t.FailNow()
	}

	err = c.Send([]byte("hello server"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(1 * time.Second)
	select {
	case sMsg := <-sMsgCh:
		if string(sMsg.Content) != "hello server" {
			t.Error("unexpected rx on server: ", string(sMsg.Content))
		}
	default:
		t.Error("server no rx")
	}

	sEvt := <-sEvtCh
	if sEvt.EventType != connection.CONNECTED {
		t.Error("server event is not connected")
		t.FailNow()
	}
	clientId := sEvt.Id
	fmt.Printf("client id is %d\r\n", clientId)

	err = s.Send(clientId, []byte("hello client"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	time.Sleep(2 * time.Second)
	select {
	case cMsg := <-cMsgCh:
		if string(cMsg.Content) != "hello client" {
			t.Error("unexpected rx on client: ", string(cMsg.Content))
		}
	default:
		t.Error("client no rx")
	}

	if err = c.Disconnect(); err != nil {
		t.Error(err)
	}
	s.Stop()
}
