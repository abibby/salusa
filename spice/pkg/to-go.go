package pkg

type GoStringer interface {
	GoString() string
}

type GoStringFunc func() string

func (f GoStringFunc) GoString() string {
	return f()
}

func Raw(src string) GoStringer {
	return GoStringFunc(func() string {
		return src
	})
}
