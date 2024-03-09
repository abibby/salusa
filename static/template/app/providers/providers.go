package providers

import (
	"github.com/abibby/salusa/database/model/modeldi"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/static/template/app"
	"github.com/abibby/salusa/static/template/app/models"
)

// Register registers any custom di providers
func Register(dp *di.DependencyProvider) {
	app.Kernel.Register(modeldi.Register[*models.User])
}
