package pkg

type ToGoer interface {
	ToGo() string
}

type ToGoFunc func() string

func (f ToGoFunc) ToGo() string {
	return f()
}

func Raw(src string) ToGoer {
	return ToGoFunc(func() string {
		return src
	})
}
