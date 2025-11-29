package repositories

import (
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wrapped-owls/goremy-di/remy"
)

func RegisterRepositories(inj remy.Injector) {
	remy.RegisterConstructorArgs1(inj, remy.Factory[mysqlrepo.RepoMySQL], mysqlrepo.NewRepoMySQL)
}
