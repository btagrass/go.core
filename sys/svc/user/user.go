package user

import (
	"fmt"
	"time"

	"github.com/btagrass/go.core/svc"
	"github.com/btagrass/go.core/sys/mdl"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

// 用户服务
type UserSvc struct {
	*svc.Svc[mdl.User]
	Perm      *casbin.SyncedEnforcer // 权限
	SignedKey []byte                 // 签名密钥
}

// 构造函数
func NewUserSvc() *UserSvc {
	s := &UserSvc{
		Svc:       svc.NewSvc[mdl.User]("sys:users"),
		SignedKey: []byte("kskj"),
	}
	model := model.NewModel()
	model.AddDef("r", "r", "sub, obj, act")                                                                                                                   // 请求（sub：用户编码，obj：请求路径，act：方法）
	model.AddDef("p", "p", "sub, obj, act, res")                                                                                                              // 角色资源（sub：角色编码，obj：请求路径，act：方法，id：资源编码，type：资源类型）
	model.AddDef("g", "g", "_, _")                                                                                                                            // 用户角色（_：用户编码，_：角色编码）
	model.AddDef("e", "e", "some(where (p.eft == allow))")                                                                                                    // 策略
	model.AddDef("m", "m", "r.sub == '300000000000001' || r.obj == '/mgt/sys/resources/menu' || g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act") // 匹配
	adapter, err := gormadapter.NewAdapterByDB(s.Db)
	if err != nil {
		logrus.Fatal(err)
	}
	perm, err := casbin.NewSyncedEnforcer(model, adapter)
	if err != nil {
		logrus.Fatal(err)
	}
	err = perm.LoadPolicy()
	if err != nil {
		logrus.Fatal(err)
	}
	s.Perm = perm

	return s
}

// 获取用户集合
func (s *UserSvc) ListUsers(conds map[string]any) ([]mdl.User, int64, error) {
	var users []mdl.User
	var count int64
	db := s.Make(conds).Joins("Dept").Find(&users)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	err := db.Error
	if err != nil {
		return users, count, err
	}

	return users, count, nil
}

// 获取用户角色集合
func (s *UserSvc) ListUserRoles(id string) ([]int64, error) {
	roles := []int64{}
	rs, err := s.Perm.GetRolesForUser(id)
	if err != nil {
		return roles, err
	}
	for _, r := range rs {
		roles = append(roles, cast.ToInt64(r))
	}

	return roles, nil
}

// 登录
func (s *UserSvc) Login(userName, password string) (*mdl.User, error) {
	var user *mdl.User
	err := s.Db.Select("id, user_name, password, frozen").First(&user, "user_name = ?", userName).Error
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if user.Frozen {
		return nil, fmt.Errorf("用户已被冻结")
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":    user.Id,
			"userName":  user.UserName,
			"expiresAt": time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	)
	user.Password = ""
	user.Token, _ = token.SignedString(s.SignedKey)

	return user, nil
}

// 移除用户集合
func (s *UserSvc) RemoveUsers(ids []string) error {
	return s.Remove(ids, "id != 300000000000001")
}

// 保存用户角色集合
func (s *UserSvc) SaveUser(user mdl.User) error {
	if len(user.Password) != 60 {
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(password)
	}
	err := s.Save(user)

	return err
}

// 保存用户角色集合
func (s *UserSvc) SaveUserRoles(id string, roles []int64) error {
	_, err := s.Perm.DeleteRolesForUser(id)
	if err != nil {
		return err
	}
	_, err = s.Perm.AddRolesForUser(id, cast.ToStringSlice(roles))

	return err
}
