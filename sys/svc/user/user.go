package user

import (
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"go.core/svc"
	"go.core/sys/mdl"
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
	userSvc := &UserSvc{
		Svc:       svc.NewSvc[mdl.User]("sys:users"),
		SignedKey: []byte("kskj"),
	}
	model := model.NewModel()
	model.AddDef("r", "r", "sub, obj, act")
	model.AddDef("p", "p", "sub, obj, act, id, type")
	model.AddDef("g", "g", "_, _")
	model.AddDef("e", "e", "some(where (p.eft == allow))")
	model.AddDef("m", "m", "r.sub == '300000000000001' || regexMatch(r.obj, '/sys/resources/menu') || g(r.sub, p.sub) && regexMatch(r.obj, p.obj) && regexMatch(r.act, p.act) && p.type == '2'")
	adapter, err := gormadapter.NewAdapterByDB(userSvc.Db)
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
	userSvc.Perm = perm

	return userSvc
}

// 获取用户角色集合
func (u *UserSvc) ListUserRoles(id string) ([]int64, error) {
	roles := []int64{}
	rs, err := u.Perm.GetRolesForUser(id)
	if err != nil {
		return roles, err
	}
	for _, r := range rs {
		roles = append(roles, cast.ToInt64(r))
	}

	return roles, nil
}

// 登录
func (u *UserSvc) Login(userName, password string) (*mdl.User, error) {
	var user *mdl.User
	err := u.Db.Select("id, user_name, password, frozen").First(&user, "user_name = ?", userName).Error
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
	user.Token, _ = token.SignedString(u.SignedKey)

	return user, nil
}

// 移除用户集合
func (u *UserSvc) RemoveUsers(ids []string) error {
	return u.Remove("id != 1 and id in ?", ids)
}

// 保存用户角色集合
func (u *UserSvc) SaveUser(user *mdl.User) error {
	if len(user.Password) != 60 {
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(password)
	}
	err := u.Save(user)

	return err
}

// 保存用户角色集合
func (u *UserSvc) SaveUserRoles(id string, roles []int64) error {
	_, err := u.Perm.DeleteRolesForUser(id)
	if err != nil {
		return err
	}
	_, err = u.Perm.AddRolesForUser(cast.ToString(id), cast.ToStringSlice(roles))

	return err
}
