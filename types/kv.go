package types

type Elem struct {
	key, value []byte
	next, prev *Elem
}

type List struct {
	head, tail *Elem
}

func (l *List) Push(key, value []byte) *Elem {
	e := &Elem{
		key:   key,
		value: value,
	}
	if l.head == nil {
		l.head = e
	} else {
		l.tail.next = e
		e.prev = l.tail
	}
	l.tail = e
	return e
}

func (l *List) Update(value []byte, e *Elem) {
	e.value = value
	if l.tail != e {
		if l.head == e {
			e.next.prev = nil
			l.head = e.next
		} else {
			e.prev.next = e.next
			e.next.prev = e.prev
		}
		l.tail.next = e
		e.prev = l.tail
		l.tail = e
		e.next = nil
	}
}

type Mapping struct {
	*List
	m map[string]*Elem
}

func NewMapping() *Mapping {
	return &Mapping{
		&List{},
		make(map[string]*Elem),
	}
}

func (mp *Mapping) Set(key, value []byte) {
	keystr := string(key)
	e := mp.m[keystr]
	if e == nil {
		mp.m[keystr] = mp.List.Push(key, value)
	} else {
		mp.List.Update(value, e)
	}
}

func (mp *Mapping) Get(key []byte) []byte {
	keystr := string(key)
	e := mp.m[keystr]
	if e != nil {
		return e.value
	}
	return nil
}
