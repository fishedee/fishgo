package web

import (
	"net/http"
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
	SessionStart(w http.ResponseWriter, r *http.Request) (session SecurityStore, err error)
}

type Database interface {
	NewSession() Database
	Close() error
	Sql(querystring string, args ...interface{}) Database
	NoAutoTime() Database
	NoAutoCondition(no ...bool) Database
	Cascade(trueOrFalse ...bool) Database
	Where(querystring string, args ...interface{}) Database
	Id(id interface{}) Database
	Distinct(columns ...string) Database
	Select(str string) Database
	Cols(columns ...string) Database
	AllCols() Database
	MustCols(columns ...string) Database
	UseBool(columns ...string) Database
	Omit(columns ...string) Database
	Nullable(columns ...string) Database
	In(column string, args ...interface{}) Database
	Incr(column string, arg ...interface{}) Database
	Decr(column string, arg ...interface{}) Database
	SetExpr(column string, expression string) Database
	Table(tableNameOrBean interface{}) Database
	Alias(alias string) Database
	Limit(limit int, start ...int) Database
	Desc(colNames ...string) Database
	Asc(colNames ...string) Database
}
