package pkg

import "fmt"

type FunctionCall struct {
	name            string
	args            []any
	hasError        bool
	returnCount     int
	returnVariables map[int]string
}

func (f *FunctionCall) WithError() *FunctionCall {
	f.hasError = true
	return f
}
func (f *FunctionCall) ReturnCount(count int) *FunctionCall {
	f.returnCount = count
	return f
}
func (f *FunctionCall) ReturnVariable(index int, name string) *FunctionCall {
	f.returnVariables[index] = name
	return f
}
func (f *FunctionCall) ToGo() string {
	src := ""
	args := ""
	for i, arg := range f.args {
		if i != 0 {
			args += ", "
		}
		if goer, ok := arg.(ToGoer); ok {
			args += goer.ToGo()
		} else {
			args += fmt.Sprintf("%#v", arg)
		}
	}

	src += "\t"

	if f.hasError {
		for i := 0; i < f.returnCount-1; i++ {
			src += "_, "
		}
		src += "err = "
	}

	src += fmt.Sprintf("%s(%s)\n", f.name, args)
	if f.hasError {
		src += "\tif err != nil {\n" +
			"\t\tpanic(err)\n" +
			"\t}"
	}
	return src
}
