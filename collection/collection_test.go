package collection

import "testing"

type TestA struct {
	Name string
	Id   int
}

func TestCollection(t *testing.T) {
	arr1 := []*TestA{
		&TestA{"a", 1},
		&TestA{"b", 2},
		&TestA{"c", 3},
		&TestA{"d", 4},
	}

	c := FromSlice(arr1, func(v interface{}, idx int) int { return v.(*TestA).Id })
	for idx, v := range c.Ints() {
		if v != arr1[idx].Id {
			t.Errorf("want %v got %v", arr1[idx].Id, v)
		}
	}

	arr2 := []*TestA{
		&TestA{"b", 2},
		&TestA{"c", 3},
		&TestA{"d", 4},
		&TestA{"a", 5},
	}
	c2 := FromSlice(arr2, func(v interface{}, idx int) int { return v.(*TestA).Id })
	mm := c.Minus(c2.Ints()).Map()
	if len(mm) != 1 {
		t.Errorf("want 1 got %v", len(mm))
	}
	if mm[arr1[0].Id].(*TestA).Name != arr1[0].Name {
		t.Errorf("want %v got %v", arr1[0].Name, mm[arr1[0].Id].(*TestA).Name)
	}

	interC := c.Intersection(c2.Ints())
	for _, v := range interC.Ints() {
		val1, ok := c.Exist(v)
		if !ok {
			t.Errorf("want %v got %v", true, ok)
		}
		val2, ok := c2.Exist(v)
		if !ok {
			t.Errorf("want %v got %v", true, ok)
		}
		if val1.(*TestA).Name != val2.(*TestA).Name {
			t.Errorf("want %v got %v", val1.(*TestA).Name, val2.(*TestA).Name)
		}
	}

	fc := c.Filter(func(v interface{}, idx int) bool { return idx == 0 })
	fcm := fc.Map()
	if len(fcm) != 1 {
		t.Errorf("want %v got %v", 1, len(fcm))
	}
	v, ok := fcm[arr1[0].Id]
	if !ok {
		t.Errorf("want %v got %v", true, ok)
	}
	if v.(*TestA).Name != arr1[0].Name {
		t.Errorf("want %v got %v", arr1[0].Name, v.(*TestA).Name)
	}

}
