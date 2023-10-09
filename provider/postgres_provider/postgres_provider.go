package postgres_provider

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"parking-service/provider/viper_provider"
)

type IPostgresProvider interface {
	GetConn() *sql.DB
	PrintConn() string
	CloseAll()
	GetDB() *sql.DB
}

type postgresProvider struct {
	Db              *sql.DB
	iConfigProvider viper_provider.IConfigProvider
}

func NewPostgresProvider(IConfigProvider viper_provider.IConfigProvider) IPostgresProvider {
	pr := postgresProvider{iConfigProvider: IConfigProvider}
	pr.GetConn()
	return &pr
}

func (pr *postgresProvider) GetConn() *sql.DB {
	connStr := pr.iConfigProvider.GetConfigEnv().PGDatabase
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf(err.Error())
		return nil
	}
	pr.Db = db
	return db
}

func(pr *postgresProvider) GetDB() *sql.DB {
	return pr.Db
}

func (pr *postgresProvider) CloseAll() {
	_ = pr.Db.Close()
}

func (pr *postgresProvider) PrintConn() string {
	return "Conn"
}
