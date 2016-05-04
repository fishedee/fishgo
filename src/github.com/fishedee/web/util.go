package web

import (
	"database/sql"
	"net/http"
	"testing"
	"time"
)

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Testing  *testing.T
}

type Security interface {
}

type SessionStore interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
	SessionRelease(w http.ResponseWriter)
	Flush() error
}

type Session interface {
	SessionStart(w http.ResponseWriter, r *http.Request) (session SessionStore, err error)
}

type DatabaseSession interface {
	Close() error
	Sql(querystring string, args ...interface{}) DatabaseSession
	NoAutoTime() DatabaseSession
	NoAutoCondition(no ...bool) DatabaseSession
	Cascade(trueOrFalse ...bool) DatabaseSession
	Where(querystring string, args ...interface{}) DatabaseSession
	Id(id interface{}) DatabaseSession
	Distinct(columns ...string) DatabaseSession
	Select(str string) DatabaseSession
	Cols(columns ...string) DatabaseSession
	AllCols() DatabaseSession
	MustCols(columns ...string) DatabaseSession
	UseBool(columns ...string) DatabaseSession
	Omit(columns ...string) DatabaseSession
	Nullable(columns ...string) DatabaseSession
	In(column string, args ...interface{}) DatabaseSession
	Incr(column string, arg ...interface{}) DatabaseSession
	Decr(column string, arg ...interface{}) DatabaseSession
	SetExpr(column string, expression string) DatabaseSession
	Table(tableNameOrBean interface{}) DatabaseSession
	Alias(alias string) DatabaseSession
	Limit(limit int, start ...int) DatabaseSession
	Desc(colNames ...string) DatabaseSession
	Asc(colNames ...string) DatabaseSession
	OrderBy(order string) DatabaseSession
	Join(join_operator string, tablename interface{}, condition string, args ...interface{}) DatabaseSession
	GroupBy(keys string) DatabaseSession
	Having(conditions string) DatabaseSession
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Query(sql string, paramStr ...interface{}) (resultsSlice []map[string][]byte, err error)
	Insert(beans ...interface{}) (int64, error)
	InsertOne(bean interface{}) (int64, error)
	Update(bean interface{}, condiBeans ...interface{}) (int64, error)
	Delete(bean interface{}) (int64, error)
	Get(bean interface{}) (bool, error)
	Find(beans interface{}, condiBeans ...interface{}) error
	Count(bean interface{}) (int64, error)
}

type Database interface {
	DatabaseSession
	NewSession() DatabaseSession
}

type Log interface {
	Emergency(format string, v ...interface{})
	Alert(format string, v ...interface{})
	Critical(format string, v ...interface{})
	Error(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Notice(format string, v ...interface{})
	Informational(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Monitor interface {
	AscErrorCount()
	AscCriticalCount()
}

type Timer interface {
	Cron(cronspec string, handler func()) error
	Interval(duraction time.Duration, handler func()) error
	Tick(duraction time.Duration, handler func()) error
}

type Queue interface {
	Produce(topicId string, data ...interface{}) error
	Consume(topicId string, listener interface{}) error
	ConsumeInPool(topicId string, listener interface{}, poolSize int) error
	Publish(topicId string, data ...interface{}) error
	Subscribe(topicId string, listener interface{}) error
	SubscribeInPool(topicId string, listener interface{}, poolSize int) error
}

type Cache interface {
	Get(key string) (string, bool)
	Set(key string, value string, timeout time.Duration)
	Del(key string)
}
