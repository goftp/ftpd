package web

import (
	"github.com/go-xweb/xweb"
	"github.com/goftp/server"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	USER_MODULE = iota + 1
	GROUP_MODULE
	PERM_MODULE
	CHGPASS_MODULE
)

var (
	DB        UserDB
	Perm      server.Perm
	Factory   server.DriverFactory
	adminUser string
)

func Web(listen, static, templates, admin, pass string, ssl bool) {
	_, err := DB.GetUser(admin)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = DB.AddUser(admin, pass)
		}
	}
	if err != nil {
		panic(err)
	}
	adminUser = admin

	app := xweb.RootApp()
	filter := xweb.NewLoginFilter(app, "userId", "/login")
	filter.AddAnonymousUrls("/login")
	app.AddFilter(filter)

	xweb.SetStaticDir(static)
	xweb.SetTemplateDir(templates)
	xweb.AddAction(&MainAction{})
	xweb.AutoAction(&UserAction{}, &GroupAction{}, &PermAction{})

	if ssl {
		//xweb.RunTLS(listen, config)
	}
	xweb.Run(listen)
}
