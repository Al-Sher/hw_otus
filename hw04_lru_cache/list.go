package hw04lrucache

// List интерфейс двусвязного списка.
type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// ListItem элемент двусвязного списка.
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

// list структура двусвязного списка.
type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

// NewList функция для создания списка.
func NewList() List {
	return new(list)
}

// Len функция для получения количества элементов в списке.
func (l *list) Len() int {
	return l.len
}

// Front функция для получения первого элемента списка.
func (l *list) Front() *ListItem {
	return l.head
}

// Back функция для получения последнего элемента списка.
func (l *list) Back() *ListItem {
	return l.tail
}

// PushFront функция для вставки нового элемента в начало списка.
func (l *list) PushFront(v interface{}) *ListItem {
	val := &ListItem{Value: v, Next: l.head}
	if l.head != nil {
		l.head.Prev = val
	}
	l.head = val

	if l.tail == nil {
		l.tail = val
	}

	l.len++

	return val
}

// PushBack функция для вставки нового элемента в конец списка.
func (l *list) PushBack(v interface{}) *ListItem {
	val := &ListItem{Value: v, Prev: l.tail}
	if l.tail != nil {
		l.tail.Next = val
	}
	l.tail = val

	if l.head == nil {
		l.head = val
	}

	l.len++

	return val
}

// Remove функция для удаления элемента из списка.
func (l *list) Remove(i *ListItem) {
	switch {
	case i.Next != nil && i.Prev != nil:
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	case i.Next != nil:
		i.Next.Prev = nil
		l.head = i.Next
	case i.Prev != nil:
		i.Prev.Next = nil
		l.tail = i.Prev
	default:
		l.tail = nil
		l.head = nil
	}

	l.len--
}

// MoveToFront функция для перевода значения в начало списка.
func (l *list) MoveToFront(i *ListItem) {
	if l.Front() == i {
		return
	}

	if i.Next != nil && i.Prev != nil {
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	} else if i.Prev != nil {
		i.Prev.Next = i.Next
		l.tail = i.Prev
	}

	i.Next = l.head
	l.head.Prev = i
	i.Prev = nil
	l.head = i
}
