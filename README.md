# collection

[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/collection.svg)](https://pkg.go.dev/github.com/gopi-frame/collection)
[![Go](https://github.com/gopi-frame/collection/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/gopi-frame/collection/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gopi-frame/collection/graph/badge.svg?token=UGVGP6QF5O)](https://codecov.io/gh/gopi-frame/collection)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/collection)](https://goreportcard.com/report/github.com/gopi-frame/collection)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
## Installation
```go
go get -u -v github.com/gopi-frame/collection
```

## Map

### Import

```go
import "github.com/gopi-frame/collection/kv"
```

### Hash Map

```go
package main

import (
	"fmt"
	"github.com/gopi-frame/collection/kv"
)

func main() {
	m := kv.NewMap[string, string]()
	// for multi-coroutines
    // m.Lock()
	// defer m.Unlock()
	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")
	m.Each(func(key string, value string) bool {
		fmt.Println(key, value)
		return true
	})
}
```

### Linked Hash Map

```go
package main

import (
	"fmt"
	"github.com/gopi-frame/collection/kv"
)

func main() {
	m := kv.NewLinkedMap[string, string]()
	// for multi-coroutines
	// m.Lock()
	// defer m.Unlock()
	m.Set("key1", "value1")
	m.Set("key2", "value2")
	m.Set("key3", "value3")
	m.Each(func(key string, value string) bool {
		fmt.Println(key, value)
		return true
	})
}
```

## List

### Import
```go
import "github.com/gopi-frame/collection/list"
```

### Array List

```go
package main

import (
	"fmt"
	"github.com/gopi-frame/collection/list"
)

func main() {
	l := list.NewList[int](1, 2, 3)
	// for multi-coroutines
	// l.Lock()
	// defer l.Unlock()
	l.Push(4, 5, 6)
	l.Set(0, 10)
	l.Remove(1)
	l.Each(func(index int, value int) bool {
		fmt.Println(index, value)
		return true
	})
}
```

### Linked List

```go
package main

import "github.com/gopi-frame/collection/list"

func main() {
	l := list.NewLinkedList[int]()
	// for multi-coroutines
	// l.Lock()
	// defer l.Unlock()
	l.Push(4, 5, 6)
	l.Set(0, 10)
	l.Remove(1)
	l.Each(func(index int, value int) bool {
		fmt.Println(index, value)
		return true
	})
}
```
## Set

### Import

```go
import "github.com/gopi-frame/collection/set"
```

### Hash Set
```go
package main

import (
	"fmt"
	"github.com/gopi-frame/collection/set"
)

func main() {
	s := set.NewSet[int](1, 2, 3)
	// for multi-coroutines
	// s.Lock()
	// defer s.Unlock()
	s.Push(4, 5, 6)
	s.Remove(1)
	s.Each(func(_ int, item string) bool {
		fmt.Println(item)
		return true
	})
}
```

### Linked Hash Set
```go
package main

import (
	"fmt"
	"github.com/gopi-frame/collection/set"
)

func main() {
	s := set.NewLinkedSet[int](1, 2, 3)
	// for multi-coroutines
	// s.Lock()
	// defer s.Unlock()
	s.Push(4, 5, 6)
	s.Remove(1)
	s.Each(func(_ int, item string) bool {
		fmt.Println(item)
		return true
	})
}
```

## Tree

### Import

```go
import "github.com/gopi-frame/collection/tree"
```

### AVL Tree

```go
package main

import (
	"fmt"
	"github.com/gopi-frame/collection/tree"
)

type Comparater struct{}

func (c Comparater) Compare(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func main() {
	t := tree.NewAVLTree[int](Comparater{})
	// for multi-coroutines
	// t.Lock()
	// defer t.Unlock()
	t.Push(2, 1, 3, 4, 0)
	t.Remove(1)
	t.Each(func(index int, value int) bool {
		fmt.Println(index, value)
		return true
	})
}
```

### Red-black Tree

```go
package main

import (
	"fmt"
	"github.com/gopi-frame/collection/tree"
)

type Comparater struct{}

func (c Comparater) Compare(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func main() {
	t := tree.NewRBTree[int](Comparater{})
	// for multi-coroutines
	// t.Lock()
	// defer t.Unlock()
	t.Push(2, 1, 3, 4, 0)
	t.Remove(1)
	t.Each(func(index int, value int) bool {
		fmt.Println(index, value)
		return true
	})
}
```

## Queue

### Import

```go
import "github.com/gopi-frame/collection/queue"
```

### Array Queue

```go
package main

import "github.com/gopi-frame/collection/queue"

func main() {
	q := queue.NewQueue[int]()
	// for multi-coroutines
	// q.Lock()
	// defer q.Unlock()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	value, ok := q.Dequeue() // 1, true
	value, ok = q.Dequeue() // 2, true
	value, ok = q.Dequeue() // 3, true
	value, ok = q.Dequeue() // 0, false
	
}
```

### Blocking Queue

```go
package main

import (
	"github.com/gopi-frame/collection/queue"
	"sync"
)

func main() {
	q := queue.NewBlockingQueue[int](10)
	// value, ok := q.Dequeue() // will block
	// value, ok := q.DequeueTimeout(time.Second) // will block for 1 sec
	value, ok := q.TryDequeue() // return 0, false
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Enqueue(1)
        }(i)
    }
	wg.Wait()
	// ok = q.Enqueue(11) // will block
	// ok = q.EnqueueTimeout(11, time.Second) // will block for 1 sec 
	ok = q.TryEnqueue(11) // return false 
}
```

### Linked Queue

```go
package main

import "github.com/gopi-frame/collection/queue"

func main() {
	q := queue.NewLinkedQueue[int]()
	// for multi-coroutines
	// q.Lock()
	// defer q.Unlock()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	value, ok := q.Dequeue() // 1, true
	value, ok = q.Dequeue() // 2, true
	value, ok = q.Dequeue() // 3, true
	value, ok = q.Dequeue() // 0, false
	
}
```

### Linked Blocking Queue

```go
package main

import (
	"github.com/gopi-frame/collection/queue"
	"sync"
)

func main() {
	q := queue.NewLinkedBlockingQueue[int](10)
	// value, ok := q.Dequeue() // will block
	// value, ok := q.DequeueTimeout(time.Second) // will block for 1 sec
	value, ok := q.TryDequeue() // return 0, false
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Enqueue(1)
        }(i)
    }
	wg.Wait()
	// ok = q.Enqueue(11) // will block
	// ok = q.EnqueueTimeout(11, time.Second) // will block for 1 sec 
	ok = q.TryEnqueue(11) // return false 
}
```

### Priority Queue

```go
package main

import "github.com/gopi-frame/collection/queue"

type Comparater struct{}

func (c Comparater) Compare(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func main() {
	q := queue.NewPriorityQueue[int](Comparater{})
	// for multi-coroutines
	// q.Lock()
	// defer q.Unlock()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	value, ok := q.Dequeue() // 1, true
	value, ok = q.Dequeue() // 2, true
	value, ok = q.Dequeue() // 3, true
	value, ok = q.Dequeue() // 0, false
}
```

### Priority Blocking Queue

```go
package main

import (
	"github.com/gopi-frame/collection/queue"
	"sync"
)

type Comparater struct{}

func (c Comparater) Compare(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func main() {
	q := queue.NewPriorityBlockingQueue[int](Comparater{}, 10)
	// value, ok := q.Dequeue() // will block
	// value, ok := q.DequeueTimeout(time.Second) // will block for 1 sec
	value, ok := q.TryDequeue() // return 0, false
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Enqueue(1)
        }(i)
    }
	wg.Wait()
	// ok = q.Enqueue(11) // will block
	// ok = q.EnqueueTimeout(11, time.Second) // will block for 1 sec 
	ok = q.TryEnqueue(11) // return false 
}
```

### Delayed Queue

```go
package main

import (
	"encoding/json"
	"github.com/gopi-frame/collection/queue"
	"sync"
	"time"
)

type DelayedItem struct {
	value int
	until time.Time
}

func NewDeleyedItem(value int, delay time.Duration) *DelayedItem {
	return &DelayedItem{
		value: value,
		until: time.Now().Add(delay),
	}
}

func (d *DelayedItem) Value() int {
	return d.value
}

func (d *DelayedItem) Until() time.Time {
	return d.until
}

func (d *DelayedItem) MarshalJSON() ([]byte, error) {
	type jsonObject struct {
		Value int       `json:"value"`
		Until time.Time `json:"until"`
	}
	return json.Marshal(jsonObject{d.value, d.until})
}

func (d *DelayedItem) UnmarshalJSON(data []byte) error {
	type jsonObject struct {
		Value int       `json:"value"`
		Until time.Time `json:"until"`
	}
	var obj jsonObject
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	d.value = obj.Value
	d.until = obj.Until
	return nil
}

func main() {
	q := queue.NewDelayedQueue[*DelayedItem]()
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Enqueue(NewDeleyedItem(i, time.Duration(10 - i) * time.Second))
        }(i)
    }
	wg.Wait()
}
```