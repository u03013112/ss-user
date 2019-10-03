package user

import (
	"errors"

	"github.com/u03013112/ss-user/sql"
)

// UserAuth : 数据库存储格式
type UserAuth struct {
	sql.BaseModel
	Username string `json:"username,omitempty"`
	Passwd   string `json:"passwd,omitempty"`
}

// UserInfo : 用户信息表
type UserInfo struct {
	sql.BaseModel
	Token string `json:"token,omitempty"`
	Role  string `json:"role,omitempty"`
}

// InitDB : 初始化表格，建议在整个数据库初始化之后调用
func InitDB() {
	sql.GetInstance().AutoMigrate(&UserAuth{}, &UserInfo{})
}

// 返回用户ID或者0
func auth(username string, passwd string) (uint, error) {
	var user UserAuth
	db := sql.GetInstance().Model(&UserAuth{})
	db = db.Where(&UserAuth{Username: username, Passwd: passwd}).First(&user)
	if db.Error == nil && db.RowsAffected == 1 {
		return user.ID, nil
	}
	return 0, db.Error
}

func getUserInfo(id uint) *UserInfo {
	var u UserInfo
	db := sql.GetInstance().Model(&UserInfo{}).Where("id = ?", id).First(&u)
	if db.Error == nil && db.RowsAffected == 1 {
	} else {
		u.ID = id
		u.Token = ""
		u.Role = ""
		sql.GetInstance().Create(&u)
	}
	return &u
}

func updateUserInfo(u *UserInfo) {
	sql.GetInstance().Model(u).Where("id = ?", u.ID).Update(UserInfo{Token: u.Token, Role: u.Role})
}

func getUserInfoByToken(token string) (*UserInfo, error) {
	var u UserInfo
	db := sql.GetInstance().Model(&UserInfo{}).Where("token = ?", token).First(&u)
	if db.Error == nil && db.RowsAffected == 1 {
		return &u, nil
	}
	return nil, errors.New("token not found")
}
