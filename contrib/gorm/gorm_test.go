package gorm_test

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	contrib "github.com/imkouga/kratos-pkg/contrib/gorm"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type kratosDefaultConfig struct {
}

func (kratosDefaultConfig) GetDriver() string {
	return "sqlite"
}
func (kratosDefaultConfig) GetSource() string {
	return "test.db"
}
func (kratosDefaultConfig) GetDebug() bool {
	return true
}

type User struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name;"`
}

func TestCreateGorm(t *testing.T) {
	testdb := filepath.Join(t.TempDir(), "test.db")

	t.Logf("test tmp db path: %s", testdb)

	db, err := contrib.NewWithOptions(kratosDefaultConfig{},
		contrib.WithDriver(sqlite.Open(testdb)),
		contrib.WithLogger(log.NewHelper(log.GetLogger()), true),
		contrib.WithTracing(),
	)

	assert.NoError(t, err, "should be create gorm db failed")

	// 创建测试表
	_ = db.AutoMigrate(&User{})
	db.Migrator()

	// 测试查询数据, 一定会 error(record not found)
	t.Run("test query db with error", func(t2 *testing.T) {
		var (
			user User
			e    error
		)
		e = db.Select("id", "name").Take(&user).Error

		assert.ErrorAs(t, e, &gorm.ErrRecordNotFound, "should be record not found error")
	})

	// 插入一条，并读取回来
	t.Run("test insert data and query one user", func(t2 *testing.T) {
		name := fmt.Sprintf("test-%d", time.Now().Unix())
		creatErr := db.Create(&User{Name: name}).Error
		assert.NoError(t2, creatErr, "should be error")

		var user User
		findErr := db.Model(&User{}).First(&user).Error
		assert.NoError(t2, findErr, "should be error")
		assert.Equal(t2, name, user.Name, "should be equal")
	})

	// 清理测试现场
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()

	assert.NoError(t, os.Remove(testdb), "should be error")
}
