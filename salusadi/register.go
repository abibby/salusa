package salusadi

import (
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/request"
)

func Register(dp *di.DependencyProvider) {
	clog.Register(dp, nil)
	request.RegisterDI(dp)
}
