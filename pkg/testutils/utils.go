package testutils

import (
	"context"
	"database/sql"
	"github.com/iota-agency/iota-sdk/modules"
	"github.com/iota-agency/iota-sdk/pkg/application"
	"github.com/iota-agency/iota-sdk/pkg/application/dbutils"
	"github.com/iota-agency/iota-sdk/pkg/composables"
	"github.com/iota-agency/iota-sdk/pkg/configuration"
	"github.com/iota-agency/iota-sdk/pkg/domain/aggregates/role"
	"github.com/iota-agency/iota-sdk/pkg/domain/aggregates/user"
	"github.com/iota-agency/iota-sdk/pkg/domain/entities/permission"
	"github.com/iota-agency/iota-sdk/pkg/domain/entities/session"
	"github.com/iota-agency/iota-sdk/pkg/event"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type TestContext struct {
	SQLDB   *sql.DB
	GormDB  *gorm.DB
	Context context.Context
	Tx      *gorm.DB
	App     application.Application
}

func MockUser(permissions ...permission.Permission) *user.User {
	return &user.User{
		ID:         1,
		FirstName:  "",
		LastName:   "",
		MiddleName: "",
		Password:   "",
		Email:      "",
		AvatarID:   nil,
		Avatar:     nil,
		EmployeeID: nil,
		LastIP:     nil,
		UILanguage: "",
		LastLogin:  nil,
		LastAction: nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Roles: []*role.Role{
			{
				ID:          1,
				Name:        "admin",
				Permissions: permissions,
			},
		},
	}
}

func MockSession() *session.Session {
	return &session.Session{
		Token:     "",
		UserID:    0,
		IP:        "",
		UserAgent: "",
		ExpiresAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func GetTestContext() *TestContext {
	conf := configuration.Use()
	db, err := dbutils.ConnectDB(
		conf.DBOpts,
		gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold:             0,
				LogLevel:                  gormlogger.Error,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
				ParameterizedQueries:      true,
			},
		),
	)
	if err != nil {
		panic(err)
	}
	app := application.New(db, event.NewEventPublisher())
	if err := modules.Load(app, modules.BuiltInModules...); err != nil {
		panic(err)
	}
	if err := app.RollbackMigrations(); err != nil {
		panic(err)
	}
	if err := app.RunMigrations(); err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	tx := db.Begin()
	ctx := composables.WithTx(context.Background(), tx)
	ctx = composables.WithParams(
		ctx,
		&composables.Params{
			IP:            "",
			UserAgent:     "",
			Authenticated: true,
			Request:       nil,
			Writer:        nil,
		},
	)

	return &TestContext{
		SQLDB:   sqlDB,
		GormDB:  db,
		Tx:      tx,
		Context: ctx,
		App:     app,
	}
}
