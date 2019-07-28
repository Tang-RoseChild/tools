package collection

import (
	"reflect"
)

type Collection struct {
	ints []int
	m    map[int]interface{}
}

func newCollection() *Collection {
	c := new(Collection)
	c.m = make(map[int]interface{})
	return c
}

func FromSlice(src interface{}, fn func(val interface{}, idx int) int) *Collection {
	c := newCollection()

	v := reflect.Indirect(reflect.ValueOf(src))
	if v.Type().Kind() != reflect.Slice {
		panic("src should be slice or pointer to slice")
	}

	for i := 0; i < v.Len(); i++ {
		key := fn(v.Index(i).Interface(), i)
		c.ints = append(c.ints, key)
		c.m[key] = v.Index(i).Interface()
	}
	return c
}

func (c *Collection) Ints() []int {
	return c.ints
}
func (c *Collection) Map() map[int]interface{} {
	return c.m
}

func (c *Collection) Exist(id int) (interface{}, bool) {
	v, ok := c.m[id]
	return v, ok
}

func (c *Collection) Add(ints []int) []int {
	return append(c.ints, ints...)
}

func (c *Collection) Minus(ints []int) *Collection {
	nc := newCollection()
	lookup := make(map[int]struct{})
	for _, v := range ints {
		lookup[v] = struct{}{}
	}
	for _, v := range c.ints {
		if _, ok := lookup[v]; !ok {
			nc.ints = append(nc.ints, v)
			nc.m[v] = c.m[v]
		}
	}
	return nc
}

func (c *Collection) Intersection(ints []int) *Collection {
	nc := newCollection()
	lookup := make(map[int]struct{})
	for _, v := range ints {
		lookup[v] = struct{}{}
	}
	for _, v := range c.ints {
		if _, ok := lookup[v]; ok {
			nc.ints = append(nc.ints, v)
			nc.m[v] = c.m[v]
		}
	}
	return nc
}

// Filter 保留fn返回为true的
func (c *Collection) Filter(fn func(v interface{}, idx int) bool) *Collection {
	nc := newCollection()
	for idx, v := range c.ints {
		if fn(c.m[v], idx) {
			nc.m[v] = c.m[v]
			nc.ints = append(nc.ints, v)
		}
	}
	return nc
}
