package db_connect

import (
	"core/src/lib/db/config_modules"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func CreateNeo4jConnection(dbconf config_modules.Neo4JConfiguration) (neo4j.Driver, error) {
	dbUri := "neo4j://%s:%s"
	dbUri = fmt.Sprintf(dbUri, dbconf.Host, dbconf.Port)
	driver, err := neo4j.NewDriver(dbUri, neo4j.BasicAuth(dbconf.User, dbconf.Password, ""))
	if err != nil {
		return nil, err
	}
	return driver, err
}
