package mysql

import (
	"entgo.io/ent/dialect/sql"
	log "github.com/sirupsen/logrus"
	"go-nuts/pkg/config"
	"go-nuts/pkg/ent"
)

type DB struct {
	*ent.Client
	cfg config.MysqlConfig
}

func NewDB(config config.MysqlConfig) *DB {
	db := &DB{
		cfg: config,
	}
	return db
}

func (db *DB) Close() error {
	// Close database.
	if db.Client != nil {
		return db.Client.Close()
	}
	return nil
}

func (db *DB) Open() (err error) {
	cfg := db.cfg
	client, err := sql.Open("mysql", cfg.Dsn())
	if err != nil {
		return err
	}

	mysql := client.DB()

	// set DBCP config
	mysql.SetMaxOpenConns(cfg.ConnMax)
	mysql.SetMaxIdleConns(cfg.ConnMax)
	mysql.SetConnMaxLifetime(cfg.ConnMaxTtl())
	mysql.SetConnMaxIdleTime(cfg.ConnIdleTtl())

	// check connection
	if err := mysql.Ping(); err != nil {
		return err
	}

	opts := []ent.Option{
		ent.Driver(client),
	}
	if cfg.LogDebugSQL {
		opts = append(opts,
			ent.Debug(),
			ent.Log(func(a ...any) {
				log.Debug(a...)
			}),
		)
	}
	db.Client = ent.NewClient(opts...)
	return nil
}
