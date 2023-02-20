package utl

import (
	"time"

	"github.com/yitter/idgenerator-go/idgen"
)

// 初始化
func init() {
	workerId := uint16(time.Now().Unix() % 64)
	options := idgen.NewIdGeneratorOptions(workerId)
	idgen.SetIdGenerator(options)
}

// 整型编码
func IntId() int64 {
	return idgen.NextId()
}

// 时间编码
func TimeId() string {
	return time.Now().Format("20060102150405.999999999")
}
