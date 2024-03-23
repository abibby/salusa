package salusadi

import (
	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/request"
)

func Register[T auth.User](dp *di.DependencyProvider) {
	clog.Register(dp, nil)
	request.Register(dp)
	auth.Register[T](dp)
}
