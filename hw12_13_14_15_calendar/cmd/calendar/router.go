package main

import (
	"errors"
	"os"

	"github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	internalrouter "github.com/lozhkindm/otus-go-hw/hw12_13_14_15_calendar/internal/router"
)

const routerChi = "chi"

var undefinedRouterType = errors.New("undefined router type")

func NewRouter(logfile *os.File, application *app.App, config Config) (app.Router, error) {
	var router app.Router

	switch config.App.RouterType {
	case routerChi:
		router = internalrouter.NewChiRouter(logfile, application)
	default:
		return nil, undefinedRouterType
	}

	return router, nil
}
