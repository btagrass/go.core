package mdl

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 模型接口
type IMdl interface {
	GetId() int64 // 获取编码
}

// 模型
type Mdl struct {
	Id        int64          `gorm:"primaryKey;autoIncrement:false;comment:编码" json:"id"` // 编码
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"-"`                               // 创建时间
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"-"`                               // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`                         // 删除时间
}

func (m Mdl) GetId() int64 {
	return m.Id
}

func (m *Mdl) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Mdl) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m)
}

func (m *Mdl) ToString() string {
	return fmt.Sprintf("%+v", m)
}
