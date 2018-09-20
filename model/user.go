package model

import (
	"fmt"

	"go_web_demo/pkg/auth"
	"go_web_demo/pkg/constvar"
	"gopkg.in/go-playground/validator.v9"
)

type UserModel struct {
	BaseModel
	Username string `json:"username" gorm:"cloumn:username;not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5,max=128"`
}

func (c *UserModel) TableName() string {
	return "tb_users"
}

// 创建用户
func (u *UserModel) Create() error {
	return DB.Self.Create(&u).Error
}

// 删除用户
func DeleteUser(id uint64) error {
	user := UserModel{}
	user.BaseModel.Id = id
	return DB.Self.Delete(&user).Error
}

// 修改用户
func (u *UserModel) Update() error {
	return DB.Self.Save(u).Error
}

// 获取单个用户
func GetUser(username string) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("username = ?", username).First(&u)
	return u, d.Error
}

// 获取用户群体
func ListUser(username string, offset, limit int) ([]*UserModel, uint64, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	users := make([]*UserModel, 0)
	var count uint64

	where := fmt.Sprintf("username like '%%%s%%'", username)
	if err := DB.Self.Model(&UserModel{}).Where(where).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := DB.Self.Where(where).Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

// 验证密码
func (u *UserModel) Compare(pwd string) error {
	return auth.Compare(u.Password, pwd)
}

// 对密码进行加密
func (u *UserModel) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}

// 验证实体参数是否符合数据库要求
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
