package components

import (
	"os/user"

	"roe/sb/statusbar"
	"roe/sb/util"
)

const userName = "user"

func startUser(cfg statusbar.ComponentConfig, update func(string), _ <-chan struct{}) {
	name := userName

	if u, err := user.Current(); err != nil {
		util.Warn("%s: %v", name, err)
		update("")
	} else {
		update(u.Username)
	}
}

func init() {
	statusbar.Register(userName, startUser)
}
