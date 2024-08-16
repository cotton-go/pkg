package mysql

import (
	"fmt"
	"testing"
	"time"

	"github.com/cotton-go/pkg/ssh"
	"gorm.io/gorm"
)

type users struct {
	ID        string     `gorm:"comment:用户ID;primary_key" json:"id"`                         //用户唯一标识
	Username  string     `gorm:"comment:用户名;unique_index" json:"username"`                   //用户名
	Nickname  *string    `gorm:"comment:昵称;default:NULL" json:"nickname,omitempty"`          //昵称
	Password  string     `gorm:"comment:密码" json:"-"`                                        //密码
	CreatedAt *time.Time `gorm:"comment:创建时间;default:current_timestamp()" json:"created_at"` //创建时间
	UpdatedAt *time.Time `gorm:"comment:创建时间;default:current_timestamp()" json:"updated_at"`
}

func (m users) TableName() string {
	return "users"
}

func TestMySQL(t *testing.T) {
	dsn := "root:casaos@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	conf := Config{
		DSN: dsn,
		SSH: &ssh.Config{
			// Host:     "name.asuscomm.cn",
			Port:     22,
			User:     "jun",
			Password: "337268759",
			Type:     ssh.ConfigTypeByPassword,
		},
	}
	dialector := New(conf)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var us []users
	db.Model(users{}).FindInBatches(&us, 2, func(tx *gorm.DB, batch int) error {
		fmt.Println("us=", us)
		return nil
	})
}
