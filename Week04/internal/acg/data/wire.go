// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package data

import (
	"Week04/internal/acg/service/http"
	"github.com/google/wire"
)

func InitApp() (app *App, cf func(), err error) {
	panic(wire.Build(NewDB, NewRedis, http.New, NewApp))
}

