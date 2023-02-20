package svc

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/btagrass/go.core/app"
	"github.com/btagrass/go.core/dao"
	"github.com/btagrass/go.core/mdl"
	"github.com/btagrass/go.core/utl"
	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	cch *cache.Cache  // 缓存
	rds *redis.Client // Redis
	db  *gorm.DB      // 数据库
)

// 初始化
func init() {
	// 缓存
	cch = cache.New(cache.NoExpiration, 5*time.Minute)
	addr := viper.GetString("redis.addr")
	if addr != "" {
		rds = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		})
	}
	// 数据库
	dsn := viper.GetString("dsn")
	if dsn != "" {
		var err error
		var dialector gorm.Dialector
		if strings.HasSuffix(dsn, ".db") {
			err = utl.MakeDir(filepath.Dir(dsn))
			if err != nil {
				logrus.Fatal(err)
			}
			dialector = sqlite.Open(dsn)
		} else if strings.HasSuffix(dsn, "&parseTime=True&loc=Local") {
			dsns := utl.Split(dsn, '/', '?')
			if len(dsns) == 3 {
				databaseName := dsns[1]
				dataSourceName := utl.Replace(dsn, databaseName, "information_schema")
				database, err := sql.Open("mysql", dataSourceName)
				if err != nil {
					logrus.Fatal(err)
				}
				defer database.Close()
				_, err = database.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", databaseName))
				if err != nil {
					logrus.Fatal(err)
				}
			}
			dialector = mysql.Open(dsn)
		}
		db, err = gorm.Open(dialector, &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: logger.New(
				log.New(io.MultiWriter(os.Stdout, app.LogFile), "", log.LstdFlags),
				logger.Config{
					SlowThreshold:             200 * time.Millisecond,
					IgnoreRecordNotFoundError: true,
					LogLevel:                  logger.LogLevel(logrus.GetLevel() - 1),
				},
			),
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			logrus.Fatal(err)
		}
		err = db.Callback().Create().Before("gorm:create").Register("gorm:id", func(d *gorm.DB) {
			if d.Statement.Schema != nil {
				id := d.Statement.Schema.LookUpField("Id")
				if id != nil {
					_, zero := id.ValueOf(d.Statement.Context, d.Statement.ReflectValue)
					if zero {
						err = id.Set(d.Statement.Context, d.Statement.ReflectValue, utl.IntId())
					}
				}
			}
		})
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

// 迁移
func Migrate(mdls []any, sqls ...string) error {
	if db == nil {
		return nil
	}
	err := db.AutoMigrate(mdls...)
	if err != nil {
		return err
	}
	for _, s := range sqls {
		err = db.Exec(s).Error
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

// 服务
type Svc[M mdl.IMdl] struct {
	*dao.Dao[M]
	Cache  *cache.Cache  // 缓存
	Redis  *redis.Client // Redis
	Prefix string        // 前缀
}

// 构造函数
func NewSvc[M mdl.IMdl](prefix string) *Svc[M] {
	return &Svc[M]{
		Cache:  cch,
		Redis:  rds,
		Prefix: prefix,
		Dao:    dao.NewDao[M](db),
	}
}

// 获取并缓存
func (s *Svc[M]) GetAndCache(id int64, expiration time.Duration) (*M, error) {
	var m *M
	key := fmt.Sprintf("%s:%d", s.Prefix, id)
	v, ok := s.Cache.Get(key)
	if ok {
		m = v.(*M)
	} else {
		var err error
		m, err = s.Get(id)
		if err != nil {
			return m, err
		}
		s.Cache.Set(key, m, expiration)
	}

	return m, nil
}

// 获取集合并缓存
func (s *Svc[M]) ListAndCache(expiration time.Duration, conds ...any) ([]M, error) {
	var ms []M
	key := fmt.Sprintf("%s:%v", s.Prefix, conds)
	v, ok := s.Cache.Get(key)
	if ok {
		ms = v.([]M)
	} else {
		var err error
		ms, _, err = s.List(conds)
		if err != nil {
			return ms, err
		}
		s.Cache.Set(key, ms, expiration)
	}

	return ms, nil
}

// 获取并Redis
func (s *Svc[M]) GetAndRedis(id int64, expiration time.Duration) (*M, error) {
	var m *M
	key := fmt.Sprintf("%s:%d", s.Prefix, id)
	err := s.Redis.Get(context.Background(), key).Scan(&m)
	if err != nil {
		m, err = s.Get(id)
		if err != nil {
			return m, err
		}
		err = s.Redis.Set(context.Background(), key, m, expiration).Err()
		if err != nil {
			return m, err
		}
	}

	return m, nil
}

// 获取集合并Redis
func (s *Svc[M]) ListAndRedis(expiration time.Duration, conds ...any) ([]M, error) {
	var ms []M
	key := fmt.Sprintf("%s:%v", s.Prefix, conds)
	err := s.Redis.Get(context.Background(), key).Scan(&ms)
	if err != nil {
		ms, _, err = s.List(conds)
		if err != nil {
			return ms, err
		}
		err = s.Redis.Set(context.Background(), key, ms, expiration).Err()
		if err != nil {
			return ms, err
		}
	}

	return ms, nil
}
