package error

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type stack struct {
	AnotherGoroutine bool
	Line             int
	FileName         string
	Method           string
	Params           map[string]interface{} `json:"Params,omitempty"`
	Error            error                  `json:"error,omitempty"`
	Code             int
	Msg              string
}

// StackError error with stack
type StackError struct {
	Statcks []*stack
}

// NewStackErrorWithCode if no raw error, than err is nil
func NewStackErrorWithCode(code int, msg string, err error, kv ...interface{}) error {
	ser := NewStackError(err, kv...)
	se := ser.(*StackError)
	se.Statcks[0].Code = code
	se.Statcks[0].Msg = msg
	return se
}

// NewStackError kv is key value pair
func NewStackError(err error, kv ...interface{}) error {
	se := new(StackError)
	stacks := genStacks(1, stackDepth)
	stacks[0].Params, stacks[0].Error = params(kv...)
	if stacks[0].Error == nil {
		stacks[0].Error = err
	}
	se.Statcks = stacks
	return se
}

func params(kv ...interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if len(kv) == 0 {
		return m, nil
	}
	if len(kv)%2 != 0 {
		return nil, NewStackError(fmt.Errorf("invalid key value pair"), "kv", kv)
	}
	for i := 0; i < len(kv); i = i + 2 {
		k, ok := kv[i].(string)
		if !ok {
			return nil, NewStackError(fmt.Errorf("key should be string"), "kv", kv)
		}
		m[k] = kv[i+1]
	}
	return m, nil
}

func (e *StackError) Error() string {
	return formatError(e)
}

func (e *StackError) String() string {
	return formatError(e)
}

// FormatType StackError format
type FormatType int

const (
	// JSON JSON
	JSON FormatType = iota + 1
	// Pretty Pretty
	Pretty
)

var stackDepth = 5
var formatError = jsonFormat

// Init default formatType is JSON, stack depth default is stackDepth
func Init(formatType FormatType, depth int) {
	if formatType == Pretty {
		formatError = prettyFormat
	} else {
		formatError = jsonFormat
	}
	if depth > 0 {
		stackDepth = depth
	}
}

// CustomFormater custom defined formater
func CustomFormater(f func(*StackError) string) {
	formatError = f
}

func jsonFormat(e *StackError) string {
	data, _ := json.Marshal(e)
	return string(data)
}

func prettyFormat(e *StackError) string {
	if e == nil || len(e.Statcks) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(1024)
	for _, s := range e.Statcks {
		needSplitParam := len(s.Msg) > 0 || s.Code > 0 || s.Error != nil
		splits := strings.Split(s.Method, "/")
		method := splits[len(splits)-1]
		b.WriteString(s.FileName + "#" + method + ":" + strconv.Itoa(s.Line) + "\n")
		if len(s.Msg) > 0 {
			b.WriteString("msg=" + s.Msg + " ")
		}
		if s.Code > 0 {
			b.WriteString("code=" + strconv.Itoa(s.Code) + " ")
		}
		if s.Error != nil {
			b.WriteString("error=" + s.Error.Error() + " ")
		}
		if needSplitParam {
			b.WriteString("\n")
		}
		for k, v := range s.Params {
			b.WriteString(k + "=" + fmt.Sprintf("%v", v) + " ")
		}
		if len(s.Params) > 0 {
			b.WriteString("\n")
		}
		if s.AnotherGoroutine {
			b.WriteString("########## another goroutine ##########\n")
		}
	}
	return b.String()
}

// WrapErr wrap err as StackError.
// if err is not *StackError,then return NewStackError
// if err is *StackError, just add kv to stack params.
func WrapErr(err error, kv ...interface{}) error {
	switch t := err.(type) {
	case *StackError:
		stack := genSingleStack(1)
		var isSameGoroutine bool
		for _, s := range t.Statcks {
			if stack.Method == s.Method && stack.Line == s.Line && stack.FileName == s.FileName {
				isSameGoroutine = true
				s.Params, _ = params(kv...)
			}
		}
		if !isSameGoroutine {
			stacks := genStacks(1, stackDepth)
			if len(stacks) > 0 {
				stacks[len(stacks)-1].AnotherGoroutine = true
			}
			t.Statcks = append(stacks, t.Statcks...)
		}
		return t
	}
	return NewStackError(err, kv...)
}

// RawErrors return all raw errors.first err is the first
// if no raw errors, then return nil
func (e *StackError) RawErrors() []error {
	var errs []error
	for _, s := range e.Statcks {
		if s.Error != nil {
			errs = append(errs, s.Error)
		}
	}
	return errs
}

// Codes return all none 0 code. first none code is the first
// if no one,then return nil
func (e *StackError) Codes() []int {
	var codes []int
	for _, s := range e.Statcks {
		if s.Code > 0 {
			codes = append(codes, s.Code)
		}
	}
	return codes
}

// genStacks skip 0 is the caller of genStatcks,for example, funcA call genStacks
// return stack's Method is funcA
func genStacks(skip int, depth int) []*stack {
	var stacks []*stack
	pcs := make([]uintptr, depth)
	runtime.Callers(skip+2, pcs)
	fs := runtime.CallersFrames(pcs)
	for {
		f, ok := fs.Next()
		stacks = append(stacks, &stack{
			Line:     f.Line,
			FileName: f.File,
			Method:   f.Function,
		})
		if !ok {
			break
		}
	}
	return stacks
}

func genSingleStack(skip int) *stack {
	stacks := genStacks(skip+1, 1)
	if len(stacks) > 0 {
		return stacks[0]
	}
	return new(stack)
}
