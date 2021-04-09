package db


import (
"database/sql"
"fmt"
_ "github.com/go-sql-driver/mysql"
"strconv"
"gorm.io/gorm"
"gorm.io/driver/mysql"
)

var tableCount int

// DB是一个数据库（操作）句柄，代表一个具有零到多个底层连接的连接池。它可以安全的被多个go程同时使用。
var MyDB *gorm.DB

type Db struct {
	Username string
	Password string
	Address  string
	Dbname   string
}

func InitConn(d Db) {
	fmt.Println("连接mysql服务器...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", d.Username, d.Password, d.Address, d.Dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	MyDB = db
}

// 创建数据库表
func CreateTables(DB *sql.DB, n int) {
	tableCount = n
	str := "CREATE TABLE IF NOT EXISTS %s (" +
		"uid BIGINT UNSIGNED," +
		"shortcode VARCHAR(20) NOT NULL," +
		"urlstr VARCHAR(500) NOT NULL," +
		"time DATETIME," +
		"INDEX (uid)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8"
	for i := 0; i < n; i++ {
		sql := fmt.Sprintf(str, "short_"+strconv.Itoa(i))

		_, err := DB.Exec(sql) // ignore_security_alert TODO: fix with gorm
		if err != nil {
			panic(err)
		}
	}
}

// 获取启动初始uid
func GetInitUid(DB *sql.DB, n int) int64 {
	var uid int64 = 0
	var id int64
	str := "select uid from %s order by uid desc limit 1"
	for i := 0; i < n; i++ {
		mySql := fmt.Sprintf(str, "short_"+strconv.Itoa(i))

		err := DB.QueryRow(mySql).Scan(&id) // ignore_security_alert,TODO: fix with gorm
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			panic(err)
		}

		if id > uid {
			uid = id
		}
	}

	return uid
}
