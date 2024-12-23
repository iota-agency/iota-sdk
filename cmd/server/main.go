package main

import (
	internalassets "github.com/iota-agency/iota-sdk/internal/assets"
	"github.com/iota-agency/iota-sdk/internal/server"
	"github.com/iota-agency/iota-sdk/modules"
	"github.com/iota-agency/iota-sdk/modules/core"
	"github.com/iota-agency/iota-sdk/modules/core/presentation/controllers"
	"github.com/iota-agency/iota-sdk/modules/finance"
	"github.com/iota-agency/iota-sdk/modules/warehouse"
	"github.com/iota-agency/iota-sdk/pkg/application"
	"github.com/iota-agency/iota-sdk/pkg/application/dbutils"
	"github.com/iota-agency/iota-sdk/pkg/configuration"
	"github.com/iota-agency/iota-sdk/pkg/event"
	"github.com/iota-agency/iota-sdk/pkg/logging"
	_ "github.com/lib/pq"
	gormlogger "gorm.io/gorm/logger"
	"log"
)

func main() {
	conf := configuration.Use()
	logFile, logger, err := logging.FileLogger(conf.LogrusLogLevel())
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logFile.Close()

	db, err := dbutils.ConnectDB(
		conf.DBOpts,
		gormlogger.New(
			logger,
			gormlogger.Config{
				SlowThreshold:             0,
				LogLevel:                  conf.GormLogLevel(),
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
				ParameterizedQueries:      true,
			},
		),
	)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	app := application.New(db, event.NewEventPublisher())
	if err := modules.Load(app, modules.BuiltInModules...); err != nil {
		log.Fatalf("failed to load modules: %v", err)
	}
	app.RegisterNavItems(core.NavItems...)
	app.RegisterNavItems(warehouse.NavItems...)
	app.RegisterNavItems(finance.NavItems...)
	app.RegisterHashFsAssets(internalassets.HashFS)
	app.RegisterControllers(
		controllers.NewStaticFilesController(app.HashFsAssets()),
		controllers.NewGraphQLController(app),
	)

	if err := dbutils.CheckModels(db, server.RegisteredModels); err != nil {
		log.Fatal(err)
	}

	options := &server.DefaultOptions{
		Logger:        logger,
		Configuration: conf,
		Db:            db,
		Application:   app,
	}
	serverInstance, err := server.Default(options)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	if err := serverInstance.Start(conf.SocketAddress); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
