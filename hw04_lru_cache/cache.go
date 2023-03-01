package hw04lrucache

import (
	"sync"
)

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
	mu       *sync.Mutex
}

// cacheItem структура элемента кэша для хранения ключа.
type cacheItem struct {
	val interface{}
	key Key
}

// NewCache создать новый кэш.
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       &sync.Mutex{},
	}
}

// Set установить значение кэша.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		c.changeElement(el, value)
		return true
	}
	item := &cacheItem{
		key: key,
		val: value,
	}
	el := c.queue.PushFront(item)
	c.items[key] = el

	if c.queue.Len() > c.capacity {
		c.deleteLastElement()
	}

	return false
}

// Get получить значение кэша по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		c.queue.MoveToFront(el)
		item := convertListItemToCacheItem(el)
		return item.val, true
	}

	return nil, false
}

// Clear очистить кэш.
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}

// changeElement изменить существующий элемент.
func (c *lruCache) changeElement(el *ListItem, value interface{}) {
	item := convertListItemToCacheItem(el)
	item.val = value
	el.Value = item
	c.queue.MoveToFront(el)
}

// deleteLastElement удалить наименнее используемый элемент кэша.
func (c *lruCache) deleteLastElement() {
	lastElement := c.queue.Back()
	lastItem := convertListItemToCacheItem(lastElement)
	c.queue.Remove(lastElement)
	delete(c.items, lastItem.key)
}

// convertListItemToCacheItem конвертировать ListItem в cacheItem.
func convertListItemToCacheItem(listItem *ListItem) *cacheItem {
	if item, ok := listItem.Value.(*cacheItem); ok {
		return item
	}

	panic("не вышло конвертировать элемент кэша в структуру")
}
