// Copyright 2019 syncd Author. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"errors"
	"fast-go/global"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type WhereParam struct {
	Field   string
	Tag     string
	Prepare interface{}
}

type ExprParam struct {
	Expression string
	Where      []WhereParam
}

type QueryParam struct {
	Fields  string
	Offset  int
	Limit   int
	Order   string
	Preload string
	Where   []WhereParam
}

type WhereBuilder struct {
	Where []WhereParam
}

func (wb *WhereBuilder) AddWhereParam(field string, prepare interface{}, tag string) *WhereBuilder {
	wb.Where = append(wb.Where, WhereParam{Field: field, Prepare: prepare, Tag: tag})
	return wb
}
func Create(model interface{}) error {
	db := global.App.DB
	db = db.Create(model)
	return db.Error
}

func GetMulti(model interface{}, query QueryParam) error {
	db := global.App.DB.Offset(query.Offset)
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	if query.Fields != "" {
		db = db.Select(query.Fields)
	}
	if query.Order != "" {
		db = db.Order(query.Order)
	}
	if query.Preload != "" {
		db = db.Preload(query.Preload)
	}
	db = parseWhereParam(db, query.Where)
	db.Find(model)
	return db.Error
}

func Count(model interface{}, count *int64, query QueryParam) error {
	db := global.App.DB.Model(model)
	db = parseWhereParam(db, query.Where)
	db = db.Count(count)
	return db.Error
}

func Delete(model interface{}, query QueryParam) bool {
	if len(query.Where) == 0 {
		return false
	}
	db := global.App.DB.Model(model)
	db = parseWhereParam(db, query.Where)
	db = db.Delete(model)
	if err := db.Error; err != nil {
		return false
	}
	return true
}

func DeleteByPk(model interface{}) bool {
	db := global.App.DB.Model(model)
	db.Delete(model)
	if err := db.Error; err != nil {
		return false
	}
	return true
}

func GetOne(model interface{}, query QueryParam) bool {
	db := global.App.DB.Model(model)
	if query.Fields != "" {
		db = db.Select(query.Fields)
	}
	db = parseWhereParam(db, query.Where)
	db = db.First(model)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func GetByPk(model interface{}, id interface{}) bool {
	db := global.App.DB.Model(model)
	db = db.First(model, id)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func UpdateByPk(model interface{}) bool {
	db := global.App.DB.Model(model)
	db = db.Updates(model)
	if err := db.Error; err != nil {
		return false
	}
	return true
}

func Update(model interface{}, data interface{}, query QueryParam) error {
	db := global.App.DB.Model(model)
	db = parseWhereParam(db, query.Where)
	db = db.Updates(data)
	return db.Error
}

func parseWhereParam(db *gorm.DB, where []WhereParam) *gorm.DB {
	if len(where) == 0 {
		return db
	}
	var (
		plain   []string
		prepare []interface{}
	)
	for _, w := range where {
		tag := w.Tag
		if tag == "" {
			tag = "="
		}
		var plainFmt string
		switch tag {
		case "IN":
			plainFmt = fmt.Sprintf("%s IN (?)", w.Field)
		default:
			plainFmt = fmt.Sprintf("%s %s ?", w.Field, tag)
		}
		plain = append(plain, plainFmt)
		prepare = append(prepare, w.Prepare)
	}
	return db.Where(strings.Join(plain, " AND "), prepare...)
}

func Expression(model interface{}, result *interface{}, query ExprParam) bool {
	db := global.App.DB.Model(model)
	db = parseWhereParam(db, query.Where)
	db = db.Select(query.Expression).Scan(result)
	if err := db.Error; err != nil {
		return false
	}
	return true
}
