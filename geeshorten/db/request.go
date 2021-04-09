package db

import (
	"gorm.io/gorm"
	"time"
)

const (
	TableName = "short"
)

// 需要储存的到数据库的数据
type Request struct {
	Uid       uint64
	Shortcode string
	UrlStr    string
	Time      time.Time
}
func (r *Request) Fields() string{
	return "(Uid,Shortcode,UrlStr,Time)"
}

// 插入数据库
func (r *Request) Insert(db *gorm.DB) error {
	err := db.Table(TableName).Create(r).Error
	if  err != nil {
		return err
	}
	return nil
}

// 查询数据
func (r *Request) Select(db *gorm.DB) error {
	record := new(Request)
	err := db.Table(TableName).Select(record.Fields()).Where(
		"uid=?", r.Uid).Find(record).Limit(1).Error
	if  err != nil {
		return err
	}

	return nil
}


