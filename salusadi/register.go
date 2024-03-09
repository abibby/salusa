package salusadi

import (
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/request"
)

func Register(dp *di.DependencyProvider) {
	request.InitDI(dp)
}
