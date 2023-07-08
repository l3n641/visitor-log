package models

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// ID 自增ID主键
type ID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

// SoftDeletes 软删除
type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type LocalTime time.Time

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

type LocalTimeStamp int64

func (l LocalTimeStamp) MarshalJSON() ([]byte, error) {
	tTime := time.Unix(int64(l), 0)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}
