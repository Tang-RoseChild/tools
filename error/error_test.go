package error

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestGenStack(t *testing.T) {
	stack := a()
	if len(stack) != 1 {
		panic("should be 1")
	}
	slashs := strings.Split(stack[0].Method, "/")
	if strings.Split(slashs[len(slashs)-1], ".")[1] != "C" {
		panic(fmt.Sprintf("should be C but got %v", stack[0].Method))
	}
	data, _ := json.Marshal(stack)
	fmt.Println("data > ", string(data))

	s2 := d()
	slashs = strings.Split(s2.Method, "/")
	if strings.Split(slashs[len(slashs)-1], ".")[1] != "D" {
		panic(fmt.Sprintf("should be C but got %v", s2.Method))
	}
}

func a() []*stack {
	return b()
}

func b() []*stack {
	return c()
}

func c() []*stack {
	return genStacks(0, 1)
}

func d() *stack {
	return genSingleStack(0)
}

func e() error {
	return NewStackError(fmt.Errorf("test"), "name", "abc")
}

func e2() error {
	return e()
}

func e3() error {
	return e2()
}

func E4() error {
	return WrapErr(e3(), "e4", "e4")
}

func F() error {
	return fmt.Errorf("test")
}

func G() error {
	var err1 error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err1 = E4()
		wg.Done()
	}()
	wg.Wait()
	return WrapErr(err1)
}

func H() error {
	return G()
}

func TestError(t *testing.T) {
	// fmt.Print(E())
	Init(Pretty, 5)
	err := E4()
	if err == nil {
		panic("should not nil")
	}

	fmt.Println(H())
}
