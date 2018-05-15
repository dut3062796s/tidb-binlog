package executor

// execute sql to mysql/tidb

import (
	"database/sql"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	pkgsql "github.com/pingcap/tidb-binlog/pkg/sql"
)

// DBConfig is the DB configuration.
type DBConfig struct {
	Host     string `toml:"host" json:"host"`
	User     string `toml:"user" json:"user"`
	Password string `toml:"password" json:"password"`
	Port     int    `toml:"port" json:"port"`
}

type mysqlExecutor struct {
	db *sql.DB
}

func newMysqlExecutor(cfg *DBConfig) (Executor, error) {
	db, err := pkgsql.OpenDB("mysql", cfg.Host, cfg.Port, cfg.User, cfg.Password)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &mysqlExecutor{db: db}, nil
}

func (m *mysqlExecutor) Execute(sqls []string, args [][]interface{}, isDDL bool) error {
	if len(sqls) == 0 {
		return nil
	}
	begin := time.Now()
	err := pkgsql.ExecuteSQLs(m.db, sqls, args, isDDL)
	if err != nil {
		return errors.Trace(err)
	}
	if cost := time.Since(begin).Seconds(); cost > 1 {
		log.Warnf("[reparo] execute sql takes %f seconds, is_ddl %v, length of sqls %d", cost, isDDL, len(sqls))
	} else {
		log.Debugf("[reparo] execute sql takes %f seconds, is_ddl %v, length of sqls %d", cost, isDDL, len(sqls))
	}
	return nil
}

func (m *mysqlExecutor) Close() error {
	return m.db.Close()
}
