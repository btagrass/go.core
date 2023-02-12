package dao

import (
	"errors"
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
func (d *Dao[M]) List(conds ...any) ([]M, error) {
	var ms []M
	db := d.Db
	if len(conds) > 0 {
		index := 0
		length := len(conds)
		cond, ok := conds[index].(map[string]any)
		if ok {
			db = db.Where(cond)
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
	err := db.Find(&ms).Error

	return ms, err
}

// 分页集合
func (d *Dao[M]) Page(conds ...any) ([]M, int64, error) {
	var ms []M
	var count int64
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
			db = db.Where(cond)
			index++
		}
		order, ok := conds[length-1].(string)
		if ok && strings.Contains(order, "order by ") {
			db = db.Order(utl.Replace(order, "order by ", ""))
			length--
		}
		if index < length {
			db = db.Where(conds[index], conds[index+1:length]...).Count(&count)
		}
	}
	err := db.Find(&ms).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		return ms, count, err
	}

	return ms, count, nil
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
func (d *Dao[M]) Update(m M, values map[string]any) error {
	err := d.Db.Model(&m).Updates(values).Error

	return err
}
