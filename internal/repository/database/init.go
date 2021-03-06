package database

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-campaign-service/internal/repository/configuration"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database/dbrepo"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database/inmemorydb"
	"github.com/StephanHCB/go-campaign-service/internal/repository/database/mysqldb"
)

var (
	ActiveRepository dbrepo.Repository
)

// use this for mocking

var fatalFunc = fatal
var infoFunc = info

// only exported so you can use it in test code - use Open() everywhere else
func SetRepository(repository dbrepo.Repository) {
	ActiveRepository = repository
}

func Open() {
	var r dbrepo.Repository
	if configuration.DatabaseUse() == "mysql" {
		infoFunc("Opening mysql database...")
		r = mysqldb.Create()
	} else {
		infoFunc("Opening inmemory database...")
		r = inmemorydb.Create()
	}
	r.Open()
	SetRepository(r)
}

func Close() {
	infoFunc("Closing database...")
	GetRepository().Close()
	SetRepository(nil)
}

func MigrateIfEnabled() {
	if configuration.MigrateDatabase() {
		infoFunc("Migrating database...")
		GetRepository().Migrate()
	} else {
		infoFunc("Not migrating database. Provide --database.migrate=true command line switch or configuration variable to enable.")
	}
}

func GetRepository() dbrepo.Repository {
	if ActiveRepository == nil {
		fatalFunc("You must Open() the database before using it. This is an error in your implementation.")
	}
	return ActiveRepository
}

func info(msg string) {
	aulogging.Logger.NoCtx().Info().Print(msg)
}

func fatal(msg string) {
	aulogging.Logger.NoCtx().Fatal().Print(msg)
}
