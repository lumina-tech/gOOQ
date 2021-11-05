package database

import (
	"fmt"
	"strings"

	"github.com/lumina-tech/gooq/pkg/generator/plugin/modelgen"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
)

type DockerizedDB struct {
	DB       *sqlx.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

type DatabaseConfig struct {
	Host           string
	Port           int64
	Username       string
	Password       string
	DatabaseName   string
	SSLMode        string
	MigrationPath  string
	ModelPath      string
	TablePath      string
	ModelOverrides modelgen.ModelOverride
}

func NewDockerizedDB(
	config *DatabaseConfig, dockerTag string,
) *DockerizedDB {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("Could not connect to docker: %s", err))
	}
	resource, err := pool.Run("postgres", dockerTag, []string{
		fmt.Sprintf("POSTGRES_USER=%s", config.Username),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", config.Password),
		fmt.Sprintf("POSTGRES_DB=%s", config.DatabaseName),
	})
	if err != nil {
		panic(fmt.Sprintf("could not start resource: %s", err))
	}
	result := DockerizedDB{
		pool:     pool,
		resource: resource,
	}

	if err = pool.Retry(func() error {
		hostPort := strings.Split(resource.GetHostPort(fmt.Sprintf("%d/tcp", config.Port)), ":")
		host := hostPort[0]
		port := hostPort[1]
		connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
			config.Username, config.Password, host, port, config.DatabaseName, config.SSLMode)
		fmt.Println(connStr)
		db, err := sqlx.Open("postgres", connStr)
		if err != nil {
			return err
		}
		result.DB = db
		return db.Ping()
	}); err != nil {
		panic(fmt.Sprintf("could not connect to docker: %s", err))
	}
	return &result
}

func (db *DockerizedDB) Close() error {
	return db.pool.Purge(db.resource)
}
