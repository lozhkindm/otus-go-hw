package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mx       sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	if item, ok := c.items[key]; ok {
		item.Value = cacheItem{key: key, value: value}
		c.queue.MoveToFront(item)
		return true
	}

	c.items[key] = c.queue.PushFront(cacheItem{key: key, value: value})
	if c.queue.Len() > c.capacity {
		if ci, ok := c.queue.Back().Value.(cacheItem); ok {
			delete(c.items, ci.key)
		}
		c.queue.Remove(c.queue.Back())
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if item, ok := c.items[key]; ok {
		if ci, ok := item.Value.(cacheItem); ok {
			c.queue.MoveToFront(item)
			return ci.value, true
		}
		return nil, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
