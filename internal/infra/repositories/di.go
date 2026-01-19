package repositories

import (
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/authtokenrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/authuserepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/denylistrepo"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wrapped-owls/goremy-di/remy"
)

func RegisterRepositories(inj remy.Injector) {
	// Base SQL connection wrapper
	remy.RegisterConstructorArgs1(inj, remy.Factory[mysqlrepo.RepoMySQL], mysqlrepo.NewRepoMySQL)

	// Start to inject all other constructors
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[*authuserepo.RepoMySQL],
		authuserepo.NewRepoMySQL,
	)
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[*authtokenrepo.RepoMySQL],
		authtokenrepo.NewRepoMySQL,
	)
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[*denylistrepo.RepoMySQL],
		denylistrepo.NewRepoMySQL,
	)
}
