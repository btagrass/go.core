package mdl

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 模型
type Model struct {
	Id        int64          `gorm:"primaryKey;autoIncrement:false;comment:编码" json:"id"` // 编码
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"-"`                               // 创建时间
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"-"`                               // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`                         // 删除时间
}

func (m *Model) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Model) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Model) ToString() string {
	return fmt.Sprintf("%+v", m)
}
