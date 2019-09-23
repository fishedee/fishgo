package sqlf

import (
	gosql "database/sql"
	. "github.com/fishedee/app/log"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type SqlfResult interface {
	LastInsertId() (int64, error)
	MustLastInsertId() int64

	RowsAffected() (int64, error)
	MustRowsAffected() int64
}

type SqlfCommon interface {
	Query(data interface{}, query string, args ...interface{}) error
	MustQuery(data interface{}, query string, args ...interface{})

	Exec(query string, args ...interface{}) (SqlfResult, error)
	MustExec(query string, args ...interface{}) SqlfResult
}

type SqlfTx interface {
	SqlfCommon
	Commit() error
	MustCommit()

	Rollback() error
	MustRollback()
}

type SqlfDB interface {
	SqlfCommon
	Begin() (SqlfTx, error)
	MustBegin() SqlfTx

	Close() error
	MustClose()
}

type SqlfDBConfig struct {
	Driver                string `config:"driver"`
	SourceName            string `config:"sourceName"`
	Debug                 bool   `config:"debug"`
	MaxOpenConnection     int    `config:"maxopenconnection"`
	MaxIdleConnection     int    `config:"maxidleconnection"`
	MaxConnectionLifeTime int    `config:"maxconnectionlifttime"`
}

func NewSqlfDbTest() SqlfDB {
	log, err := NewLog(LogConfig{
		Driver: "console",
	})
	if err != nil {
		panic(err)
	}
	db, err := NewSqlfDB(log, SqlfDBConfig{
		Driver:     "sqlite3",
		SourceName: ":memory:",
		Debug:      true,
	})
	if err != nil {
		panic(err)
	}
	return db
}

func NewSqlfDB(log Log, config SqlfDBConfig) (SqlfDB, error) {
	db, err := gosql.Open(config.Driver, config.SourceName)
	if err != nil {
		return nil, err
	}
	isDebug := config.Debug
	if config.MaxOpenConnection > 0 {
		db.SetMaxOpenConns(config.MaxOpenConnection)
	}
	if config.MaxIdleConnection <= 0 {
		config.MaxIdleConnection = 100
	}
	db.SetMaxIdleConns(config.MaxIdleConnection)
	if config.MaxConnectionLifeTime <= 0 {
		//每个连接默认最长使用1天
		config.MaxConnectionLifeTime = 3600 * 24
	}
	db.SetConnMaxLifetime(time.Duration(int64(time.Second) * int64(config.MaxConnectionLifeTime)))
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &dbImplement{
		db:      db,
		log:     log,
		isDebug: isDebug,
	}, nil
}

type dbImplement struct {
	db      *gosql.DB
	log     Log
	isDebug bool
}

func (this *dbImplement) Query(data interface{}, query string, args ...interface{}) error {
	sqlRunner := func() (string, error) {
		sql, args, err := genSql(query, args)
		if err != nil {
			return query, err
		}
		rows, err := this.db.Query(sql, args...)
		if err != nil {
			return sql, err
		}
		defer rows.Close()
		err = extractResult(data, rows)
		if err != nil {
			return sql, err
		}
		return sql, nil
	}

	return runSql(this.isDebug, this.log, sqlRunner)
}

func (this *dbImplement) MustQuery(data interface{}, query string, args ...interface{}) {
	err := this.Query(data, query, args...)
	if err != nil {
		panic(err)
	}
}

func (this *dbImplement) Exec(query string, args ...interface{}) (SqlfResult, error) {
	var execResult SqlfResult
	sqlRunner := func() (string, error) {
		sql, args, err := genSql(query, args)
		if err != nil {
			return query, err
		}

		result, err := this.db.Exec(sql, args...)
		if err != nil {
			return sql, err
		}

		execResult = &resultImplement{result: result}
		return sql, nil
	}

	err := runSql(this.isDebug, this.log, sqlRunner)

	return execResult, err
}

func (this *dbImplement) MustExec(query string, args ...interface{}) SqlfResult {
	result, err := this.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *dbImplement) Begin() (SqlfTx, error) {
	tx, err := this.db.Begin()
	if err != nil {
		return nil, err
	}
	return &txImplement{tx: tx,
		isDebug: this.isDebug,
		log:     this.log,
	}, nil
}

func (this *dbImplement) MustBegin() SqlfTx {
	tx, err := this.Begin()
	if err != nil {
		panic(err)
	}
	return tx
}

func (this *dbImplement) Close() error {
	return this.db.Close()
}

func (this *dbImplement) MustClose() {
	err := this.Close()
	if err != nil {
		panic(err)
	}
}

type resultImplement struct {
	result gosql.Result
}

func (this *resultImplement) LastInsertId() (int64, error) {
	return this.result.LastInsertId()
}

func (this *resultImplement) MustLastInsertId() int64 {
	result, err := this.LastInsertId()
	if err != nil {
		panic(err)
	}
	return result
}

func (this *resultImplement) RowsAffected() (int64, error) {
	return this.result.RowsAffected()
}

func (this *resultImplement) MustRowsAffected() int64 {
	result, err := this.RowsAffected()
	if err != nil {
		panic(err)
	}
	return result
}

type txImplement struct {
	tx      *gosql.Tx
	log     Log
	isDebug bool
}

func (this *txImplement) Query(data interface{}, query string, args ...interface{}) error {
	sqlRunner := func() (string, error) {
		sql, args, err := genSql(query, args)
		if err != nil {
			return query, err
		}
		rows, err := this.tx.Query(sql, args...)
		if err != nil {
			return sql, err
		}
		defer rows.Close()
		err = extractResult(data, rows)
		if err != nil {
			return sql, err
		}
		return sql, nil
	}

	return runSql(this.isDebug, this.log, sqlRunner)
}

func (this *txImplement) MustQuery(data interface{}, query string, args ...interface{}) {
	err := this.Query(data, query, args...)
	if err != nil {
		panic(err)
	}
}

func (this *txImplement) Exec(query string, args ...interface{}) (SqlfResult, error) {
	var execResult SqlfResult
	sqlRunner := func() (string, error) {
		sql, args, err := genSql(query, args)
		if args != nil {
			return query, err
		}

		result, err := this.tx.Exec(sql, args...)
		if err != nil {
			return sql, err
		}

		execResult = &resultImplement{result: result}
		return sql, nil
	}

	err := runSql(this.isDebug, this.log, sqlRunner)

	return execResult, err
}

func (this *txImplement) MustExec(query string, args ...interface{}) SqlfResult {
	result, err := this.Exec(query, args)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *txImplement) Commit() error {
	return this.tx.Commit()
}

func (this *txImplement) MustCommit() {
	err := this.Commit()
	if err != nil {
		panic(err)
	}
}

func (this *txImplement) Rollback() error {
	return this.tx.Rollback()
}

func (this *txImplement) MustRollback() {
	err := this.Rollback()
	if err != nil {
		panic(err)
	}
}
