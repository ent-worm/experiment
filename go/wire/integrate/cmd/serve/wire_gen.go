// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package serve

import (
	"github.com/google/wire"
	"github.com/wenzong/demo/biz/user"
	"github.com/wenzong/demo/infra/app"
	"github.com/wenzong/demo/infra/config"
	"github.com/wenzong/demo/infra/db"
	"github.com/wenzong/demo/infra/grpc"
	"github.com/wenzong/demo/infra/http"
	"github.com/wenzong/demo/infra/log"
)

// Injectors from wire.go:

func App() (*app.App, func(), error) {
	viper := config.NewConfig()
	option := http.NewOption(viper)
	defaultConn := db.NewDefaultConn(viper)
	repository := user.NewRepository(defaultConn)
	service := user.NewService(repository)
	controller := user.NewController(service)
	handler := Router(controller)
	server := http.NewServer(option, handler)
	logger, cleanup := log.New()
	v := gRPCServerOptions(logger)
	v2 := log.CtxLogger()
	userServer := user.NewServer(service, v2)
	registerServiceFunc := gRPCRegisterServiceFn(userServer)
	grpcServer := grpc.NewServer(v, registerServiceFunc)
	listener := grpc.NewListener(viper)
	appApp := app.NewApp(server, grpcServer, listener)
	return appApp, func() {
		cleanup()
	}, nil
}

// wire.go:

var ProviderSet = wire.NewSet(
	Router,
	gRPCServerOptions,
	gRPCRegisterServiceFn,
)

var Set = wire.NewSet(app.ProviderSet, config.ProviderSet, db.ProviderSet, grpc.ProviderSet, http.ProviderSet, log.ProviderSet, user.ProviderSet, ProviderSet)
