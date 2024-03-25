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

package containers

import (
	"sync"
)

type Item struct {
	id   int
	item any
}

type List struct {
	items []Item
	lock  sync.Mutex
}

func NewList() *List {
	return &List{
		items: make([]Item, 0),
		lock:  sync.Mutex{},
	}
}

func (l *List) AddOrUpdate(id int, item any) {

	index, _ := l.Get(id)

	l.lock.Lock()
	defer l.lock.Unlock()

	if index >= 0 {
		l.items[index].item = item
	} else {
		l.items = append(l.items, Item{
			id:   id,
			item: item,
		})
	}
}

func (l *List) Get(id int) (index int, item any) {
	index = -1
	item = nil

	l.lock.Lock()
	defer l.lock.Unlock()

	for i, content := range l.items {
		if content.id == id {
			item = content.item
			index = i
			return
		}
	}

	return
}

func (l *List) GetIds() (ids []int) {
	l.lock.Lock()
	defer l.lock.Unlock()

	for _, item := range l.items {
		ids = append(ids, item.id)
	}

	return
}

func (l *List) Exist(id int) bool {
	if i, _ := l.Get(id); i < 0 {
		return false
	}
	return true
}

func (l *List) Delete(id int) any {
	index, oldItem := l.Get(id)

	if index < 0 {
		return oldItem
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	l.items[index] = l.items[len(l.items)-1]
	l.items = l.items[:len(l.items)-1]

	return oldItem
}

func (l *List) Reset() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.items = make([]Item, 0)
}
