/**
 * Copyright © 2023, Staufi Tech - Switzerland
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

package containers

import "testing"

func TestContainersUniqueIdList(t *testing.T) {
	l := NewList()
	ids := l.GetIds()
	if len(ids) != 0 {
		t.Error("new initialized list len != 0: ", len(ids))
	}

	l.AddOrUpdate(1234, 34.2345)
	ids = l.GetIds()
	if len(ids) != 1 {
		t.Error("len != 1 after adding one item: ", len(ids))
	}
	if ids[0] != 1234 {
		t.Error("id not match: ", ids[0])
	}
	idx, item := l.Get(1234)
	if idx != 0 {
		t.Error("unexpected index of first item: ", idx)
	}
	if item.(float64) != 34.2345 {
		t.Error("item value don't match: ", item.(float64))
	}

	l.AddOrUpdate(3312, int(3412))
	l.AddOrUpdate(2212, "hi there")

	ids = l.GetIds()
	if len(ids) != 3 {
		t.Error("list unexpected len !=3 : ", len(ids))
	}

	_, item1 := l.Get(1234)
	_, item2 := l.Get(3312)
	_, item3 := l.Get(2212)
	if item1.(float64) != 34.2345 || item2.(int) != 3412 || item3 != "hi there" {
		t.Error("some of the items don't contain the initalized values: ", item1, item2, item3)
	}

	l.AddOrUpdate(2212, "HelloWorld")
	ids = l.GetIds()
	if len(ids) != 3 {
		t.Error("list unexpected len !=3 : ", len(ids))
	}

	_, item1 = l.Get(1234)
	_, item2 = l.Get(3312)
	_, item3 = l.Get(2212)
	if item1.(float64) != 34.2345 || item2.(int) != 3412 || item3 != "HelloWorld" {
		t.Error("some of the items don't contain the updated values: ", item1, item2, item3)
	}

	l.Delete(3312)
	ids = l.GetIds()
	if len(ids) != 2 {
		t.Error("list unexpected len !=2 : ", len(ids))
	}

	_, item1 = l.Get(1234)
	_, item2 = l.Get(3312)
	_, item3 = l.Get(2212)
	if item1.(float64) != 34.2345 || item2 != nil || item3 != "HelloWorld" {
		t.Error("some of the items don't contain the updated values: ", item1, item2, item3)
	}

	l.Reset()
	ids = l.GetIds()
	if len(ids) != 0 {
		t.Error("list unexpected len after reset: ", len(ids))
	}
}
