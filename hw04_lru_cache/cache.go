package hw04lrucache

import "sync"

// Key ключ для кэша.
type Key string

// Cache интерфейс кэша.
type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

// lruCache структура кэша.
type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       *sync.RWMutex
}

// NewCache создать новый кэш.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       &sync.RWMutex{},
	}
}

// Set установить значение кэша.
func (c *lruCache) Set(key Key, value interface{}) bool {
	defer c.mu.Unlock()
	c.mu.Lock()
	if el, ok := c.items[key]; ok {
		el.Value = value
		c.queue.MoveToFront(el)
		return true
	}
	el := c.queue.PushFront(value)
	c.items[key] = el

	if c.queue.Len() > c.capacity {
		lastElement := c.queue.Back()
		c.queue.Remove(lastElement)
		c.dropFromMap(lastElement)
	}

	return false
}

// Get получить значение кэша по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	defer c.mu.RUnlock()
	c.mu.RLock()
	if el, ok := c.items[key]; ok {
		c.queue.MoveToFront(el)
		return el.Value, true
	}

	return nil, false
}

// Clear очистить кэш.
func (c *lruCache) Clear() {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}

// dropFromMap удалить элемент из map по значению.
func (c *lruCache) dropFromMap(item *ListItem) {
	for key, v := range c.items {
		if v == item {
			delete(c.items, key)
		}
	}
}
