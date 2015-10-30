package web

import (
	"time"

	"github.com/goftp/server"
	"github.com/lunny/tango"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"
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

type auther interface {
	AskLogin() bool
	IsLogin() bool
	LoginUserId() string
}

func auth() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		if a, ok := ctx.Action().(auther); ok {
			if a.AskLogin() {
				if !a.IsLogin() {
					ctx.Redirect("/login")
					return
				}
			}
		}
		ctx.Next()
	}
}

const (
	timeout = time.Minute * 20
)

func Web(listen, static, templates, admin, pass string, tls bool, certFile, keyFile string) {
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

	t := tango.Classic()
	t.Use(tango.Static(tango.StaticOptions{
		RootPath: static,
	}))
	t.Use(renders.New(renders.Options{
		Reload:    true, // if reload when template is changed
		Directory: templates,
	}))
	t.Use(session.New(session.Options{
		MaxAge: timeout,
	}))
	t.Use(auth())
	t.Use(binding.Bind())
	t.Use(xsrf.New(timeout))

	t.Get("/", new(MainAction))
	t.Any("/login", new(LoginAction))
	t.Get("/logout", new(LogoutAction))
	t.Group("/user", func(g *tango.Group) {
		g.Get("/", new(UserAction))
		g.Any("/add", new(UserAddAction))
		g.Any("/edit", new(UserEditAction))
		g.Any("/del", new(UserDelAction))
	})

	t.Group("/group", func(g *tango.Group) {
		g.Get("/", new(GroupAction))
		g.Get("/add", new(GroupAddAction))
		g.Get("/edit", new(GroupEditAction))
		g.Get("/del", new(GroupDelAction))
	})
	t.Group("/perm", func(g *tango.Group) {
		g.Get("/", new(PermAction))
		g.Any("/add", new(PermAddAction))
		g.Any("/edit", new(PermEditAction))
		g.Any("/del", new(PermDelAction))
		g.Any("/updateOwner", new(PermUpdateOwner))
		g.Any("/updateGroup", new(PermUpdateGroup))
		g.Any("/updatePerm", new(PermUpdatePerm))
	})

	if tls {
		t.RunTLS(certFile, keyFile, listen)
		return
	}

	t.Run(listen)
}
