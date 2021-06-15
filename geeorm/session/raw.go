package session

import (
	"database/sql"
	"github.com/ZhaoxingZhang/go-gees/geecommon/log"
	"github.com/ZhaoxingZhang/go-gees/geeorm/clause"
	"github.com/ZhaoxingZhang/go-gees/geeorm/dialect"
	"github.com/ZhaoxingZhang/go-gees/geeorm/schema"
	"strings"
)

//  db = 使用 sql.Open() 方法连接数据库成功之后返回的指针
//  第二个和第三个成员变量用来拼接 SQL 语句和 SQL 语句中占位符的对应值。
//  用户调用 Raw() 方法即可改变这两个变量的值
type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	tx       *sql.Tx

	refTable *schema.Schema
	clause   clause.Clause

	sql      strings.Builder
	sqlVars  []interface{}
}

// CommonDB is a minimal function set of db
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// DB returns tx if a tx begins. otherwise return *sql.DB
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

//  封装有 2 个目的，一是统一打印日志（包括 执行的SQL 语句和错误日志）。
//  二是执行完成后，清空 (s *Session).sql 和 (s *Session).sqlVars 两个变量。
//  这样 Session 可以复用，开启一次会话，可以执行多次 SQL

// QueryRow gets a record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of records from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}