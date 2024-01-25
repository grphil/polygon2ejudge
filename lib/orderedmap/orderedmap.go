package orderedmap

import (
	"bytes"
	"fmt"
)

type UnquotedStr string
type Comment bool
type NoVal bool

type Pair struct {
	Key   string
	Value interface{}
}

type OrderedMap struct {
	values []*Pair
}

func New() *OrderedMap {
	return &OrderedMap{}
}

func (m *OrderedMap) Set(key string, value interface{}) {
	for _, v := range m.values {
		if v.Key == key {
			v.Value = value
			return
		}
	}
	m.values = append(m.values, &Pair{key, value})
}

func (m *OrderedMap) Remove(key string) {
	for i, p := range m.values {
		if p.Key == key {
			m.values = append(m.values[:i], m.values[i+1:]...)
		}
	}
}

func (m *OrderedMap) Get(key string) (interface{}, bool) {
	for _, v := range m.values {
		if v.Key == key {
			return v.Value, true
		}
	}
	return nil, false
}

func (m *OrderedMap) GetStr(key string) string {
	v, ok := m.Get(key)
	if !ok {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	case UnquotedStr:
		return string(v.(UnquotedStr))
	default:
		return ""
	}
}

func (m *OrderedMap) GetBool(key string) bool {
	v, ok := m.Get(key)
	if !ok {
		return false
	}
	switch v.(type) {
	case NoVal:
		return true
	case bool:
		return v.(bool)
	default:
		return false
	}
}

func (m *OrderedMap) GetInt(key string) int {
	v, ok := m.Get(key)
	if !ok {
		return 0
	}

	val, ok := v.(int)
	if ok {
		return val
	} else {
		return 0
	}
}

func (m *OrderedMap) Write(b *bytes.Buffer) {
	m.WritePrefix(b, "")
}

func (m *OrderedMap) WritePrefix(b *bytes.Buffer, pref string) {
	for _, p := range m.values {
		b.WriteString(pref)
		var txt string
		switch p.Value.(type) {
		case string:
			txt = fmt.Sprintf("%s = \"%s\"", p.Key, p.Value)
		case UnquotedStr, int:
			txt = fmt.Sprintf("%s = %v", p.Key, p.Value)
		case bool:
			val := p.Value.(bool)
			if val {
				txt = p.Key
			} else {
				txt = fmt.Sprintf("%s = 0", p.Key)
			}
		default:
			txt = p.Key
		}
		b.WriteString(txt)
		b.WriteRune('\n')
	}
	b.WriteRune('\n')
}

func (m *OrderedMap) Update(o *OrderedMap) {
	for _, p := range o.values {
		m.Set(p.Key, p.Value)
	}
}

func (m *OrderedMap) Clone() *OrderedMap {
	res := New()
	for _, p := range m.values {
		res.values = append(res.values, &Pair{p.Key, p.Value})
	}
	return res
}
