package user

import (
	"crypto/md5"
	"fmt"
)

const Salt = "myEasySalt"

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"username"`
	Password string `json:"-"`
	RoleID   int    `json:"-"`
}

func (user User) ValidPassword(password string) bool {
	return user.Password == GetMD5Password(password)
}

func GetMD5Password(password string) string {
	if password == "" {
		return ""
	}
	data := []byte(password + Salt)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (user User) IsAdmin() bool {
	//TODO заменить айдишник на нейм роли
	return user.RoleID == 1
}
