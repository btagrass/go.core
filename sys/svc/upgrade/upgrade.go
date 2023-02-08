package upgrade

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
)

// 升级服务
type UpgradeSvc struct {
}

// 构造函数
func NewUpgradeSvc() *UpgradeSvc {
	return &UpgradeSvc{}
}

// 升级
func (u *UpgradeSvc) Upgrade(ver string) (string, string, error) {
	var filePath, fileVer string
	files, err := filepath.Glob("data/upgrades/*")
	if err != nil {
		return filePath, fileVer, err
	}
	if len(files) > 0 {
		f := files[0]
		_, fVer, _ := strings.Cut(filepath.Base(f), "_")
		if u.compare(fVer, ver) {
			filePath = f
			fileVer = fVer
		}
	}

	return filePath, fileVer, nil
}

// 比较，版本格式：1.0.0
func (u *UpgradeSvc) compare(v1, v2 string) bool {
	var r bool
	v1s := strings.Split(v1, ".")
	v2s := strings.Split(v2, ".")
	for i := 0; i < len(v1s); i++ {
		v1i := cast.ToInt(v1s[i])
		v2i := cast.ToInt(v2s[i])
		if v1i > v2i {
			r = true
			break
		} else if v1i < v2i {
			break
		}
	}

	return r
}
