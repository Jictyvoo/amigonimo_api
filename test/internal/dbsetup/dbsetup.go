package dbsetup

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/bootstrap"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

type ConnectionFactory struct {
	rootDB *sql.DB
	dbConf config.Database
}

func NewConnectionFactory(
	ctx context.Context, conf config.Database,
) (factory ConnectionFactory, err error) {
	conf.Database = "" // Remove database name
	factory.dbConf = conf
	mysqlConf := bootstrap.MySQLConfig(factory.dbConf)

	factory.rootDB, err = bootstrap.PerformDatabaseConnection(ctx, mysqlConf, conf.Timeout)
	return factory, err
}

func (fac *ConnectionFactory) NewDatabase(ctx context.Context, dbName string) (
	newDb struct {
		Connection *sql.DB
		Name       string
	},
	err error,
) {
	newDb.Name = dbNameNormalizer(dbName)
	if _, err = fac.rootDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", newDb.Name)); err != nil {
		return newDb, err
	}

	mysqlConf := bootstrap.MySQLConfig(fac.dbConf)
	mysqlConf.DBName = newDb.Name
	mysqlConf.MultiStatements = true
	newDb.Connection, err = bootstrap.PerformDatabaseConnection(ctx, mysqlConf, fac.dbConf.Timeout)
	return newDb, err
}

func (fac *ConnectionFactory) Close() error {
	if fac.rootDB != nil {
		err := fac.rootDB.Close()
		fac.rootDB = nil
		return err
	}

	return nil
}
