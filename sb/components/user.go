package components

import (
	"os/user"

	"roe/sb/statusbar"
	"roe/sb/util"
)

type User struct {
	statusbar.BaseComponentConfig
}

func NewUser() *User {
	base := statusbar.NewBaseComponentConfig("user", 0, 0)
	return &User{*base}
}

func (u *User) Start(update func(string), _ <-chan struct{}) {
	if user, err := user.Current(); err != nil {
		util.Warn("%s: %v", u.Name, err)
		update("")
	} else {
		update(user.Username)
	}
}
