package dao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btagrass/go.core/mdl"
	"github.com/btagrass/go.core/utl"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 数据访问对象
type Dao[M mdl.IMdl] struct {
	Db *gorm.DB
}

// 构造函数
func NewDao[M mdl.IMdl](db *gorm.DB) *Dao[M] {
	return &Dao[M]{
		Db: db,
	}
}

// 获取
func (d *Dao[M]) Get(conds ...any) (*M, error) {
	var m M
	err := d.Db.First(&m, conds...).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &m, nil
}

// 获取集合
func (d *Dao[M]) List(conds ...any) ([]M, int64, error) {
	var ms []M
	var count int64
	db := d.Make(conds...).Find(&ms)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	err := db.Error
	if err != nil {
		return ms, count, err
	}

	return ms, count, nil
}

// 组装
func (d *Dao[M]) Make(conds ...any) *gorm.DB {
	db := d.Db
	if len(conds) > 0 {
		index := 0
		length := len(conds)
		cond, ok := conds[index].(map[string]any)
		if ok {
			size, ok := cond["size"]
			if ok {
				db = db.Limit(cast.ToInt(size))
				delete(cond, "size")
			}
			current, ok := cond["current"]
			if ok {
				db = db.Offset(cast.ToInt(size) * (cast.ToInt(current) - 1))
				delete(cond, "current")
			}
			var keys []string
			var values []any
			for k, v := range cond {
				value, ok := v.(string)
				if ok {
					if value != "" {
						keys = append(keys, fmt.Sprintf("%s like ?", k))
						values = append(values, fmt.Sprintf("%%%s%%", v))
					}
					delete(cond, k)
				}
			}
			if len(keys) > 0 {
				db = db.Where(strings.Join(keys, " and "), values...)
			}
			index++
		}
		order, ok := conds[length-1].(string)
		if ok && strings.Contains(order, "order by ") {
			db = db.Order(utl.Replace(order, "order by ", ""))
			length--
		}
		if index < length {
			db = db.Where(conds[index], conds[index+1:length]...)
		}
	}

	return db
}

// 移除
func (d *Dao[M]) Remove(conds ...any) error {
	err := d.Db.Delete(new(M), conds...).Error

	return err
}

// 保存
func (d *Dao[M]) Save(m M, clauses ...clause.Expression) error {
	if len(clauses) == 0 {
		clauses = []clause.Expression{
			clause.OnConflict{
				UpdateAll: true,
			},
		}
	}
	err := d.Db.Clauses(clauses...).Create(&m).Error

	return err
}

// 事务
func (d *Dao[M]) Trans(funcs ...func(tx *gorm.DB) error) error {
	err := d.Db.Transaction(func(tx *gorm.DB) error {
		for _, f := range funcs {
			err := f(tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// 更新
func (d *Dao[M]) Update(values map[string]any, conds ...any) error {
	err := d.Make(conds...).Model(new(M)).Updates(values).Error

	return err
}
