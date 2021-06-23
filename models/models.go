package models

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gogs.buffalo-robot.com/zouhy/micro/lib"
	ztime "gogs.buffalo-robot.com/zouhy/micro/time"
	"gorm.io/gorm"
)

var Node *snowflake.Node

func init() {
	// hostName, err := os.Hostname()
	// if err != nil {
	// 	panic(err)
	// }
	ip := lib.GetIP()
	var seed int
	var err error
	if ip == "" {
		seed = 1
	} else {
		arr := strings.Split(ip, ".")
		lastOne := arr[3]
		seed, err = strconv.Atoi(lastOne)
		if err != nil {
			panic(err)
		}
	}
	idGenerator, err := snowflake.NewNode(int64(seed))
	if err != nil {
		panic(err)
	}
	Node = idGenerator
}

//Model 所有表格都有的基础模型
type Model struct {
	ID        Int64Str       `gorm:"column:id;type:bigint;not null;primaryKey;" json:"-" redis:"primary"`                     // 主键
	CreatedAt ztime.Time     `gorm:"column:created_at;type:datetime;not null;" json:"-" form:"created_at"`                    // 创建时间
	UpdatedAt ztime.Time     `gorm:"column:updated_at;type:datetime; null;" json:"-" form:"updated_at"`                       // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index:delete_at_index;null;" json:"-" form:"deleted_at" ` //删除时间
	CreatedBy string         `gorm:"column:created_by;type:varchar(50);default:0;not null;" json:"-" form:"created_by"`       // 创建人
	UpdatedBy string         `gorm:"column:updated_by;type:varchar(50);default:0;not null;" json:"-" form:"updated_by"`       // 更新人
	DeletedBy string         `gorm:"column:deleted_by;type:varchar(50);default:0;not null;" json:"-" form:"deleted_by"`       // 删除人
}

func ModelNameMap() map[string]string {
	return map[string]string{
		"ID":        "id",
		"CreatedAt": "created_at",
		"UpdatedAt": "updated_at",
		"DeletedAt": "deleted_at",
		"CreatedBy": "created_by",
		"UpdatedBy": "updated_by",
		"DeletedBy": "deleted_by",
	}
}

//BeforeCreate 基本类的创建钩子
func (model *Model) BeforeCreate(tx *gorm.DB) error {

	sid := Int64Str(Node.Generate().Int64())

	// if err := tx.Model(model).Updates(map[string]interface{}{"id": id, "_id": sid, "CreatedAt": ToJSONTime(time.Now())}).Error; err != nil {
	// 	return err
	// }
	model.ID = sid

	model.CreatedAt = ztime.Now()
	return nil
}

//BeforeUpdate 基本类的更新钩子
func (model *Model) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = ztime.Now()

	return nil
}

//GetID 获得唯一ID
func GetID() string {
	return primitive.NewObjectID().Hex()
}

func GetSFID() uint64 {
	return uint64(Node.Generate().Int64())
}
