package pkg

import "fmt"

type GoStringFunc func() string

func (f GoStringFunc) GoString() string {
	return f()
}

func Raw(src string) fmt.GoStringer {
	return GoStringFunc(func() string {
		return src
	})
}
