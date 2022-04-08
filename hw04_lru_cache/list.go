package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	first *ListItem
	last  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{Value: v, Next: l.first}

	if l.first == nil && l.last == nil {
		l.first = &newItem
		l.last = &newItem
	} else {
		l.first.Prev = &newItem
		l.first = &newItem
	}

	l.len++
	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{Value: v, Prev: l.last}

	if l.first == nil && l.last == nil {
		l.first = &newItem
		l.last = &newItem
	} else {
		l.last.Next = &newItem
		l.last = &newItem
	}

	l.len++
	return l.last
}

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.last = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.first = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}
	l.Remove(i)
	i.Prev = nil
	i.Next = l.first
	l.first.Prev = i
	l.first = i
	l.len++
}

func NewList() List {
	return new(list)
}
