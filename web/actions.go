package web

import (
	"fmt"

	"github.com/lunny/tango"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"
)

var _ auther = new(BaseAction)

type BaseAction struct {
	renders.Renderer
	session.Session
	tango.Ctx
	binding.Binder
	curModule int
}

func (a *BaseAction) AskLogin() bool {
	return false
}

type BaseAuthAction struct {
	BaseAction
}

func (a *BaseAuthAction) AskLogin() bool {
	return true
}

func (a *BaseAction) IsLogin() bool {
	id := a.Session.Get("userId")
	fmt.Println(id)
	return id != nil
}

func (a *BaseAction) SetLogin(user string) {
	a.Session.Set("userId", user)
}

func (a *BaseAction) LoginUserId() string {
	userId := a.Session.Get("userId")
	if userId == nil {
		return ""
	}
	return userId.(string)
}

func (a *BaseAction) Logout() {
	a.Session.Del("userId")
}

func (a *BaseAction) IsAdmin() bool {
	s := a.Session.Get("userId")
	return s != nil && s.(string) == adminUser
}

func (a *BaseAction) Render(tmpl string, vars ...renders.T) error {
	var t = renders.T{
		"IsLogin": a.IsLogin(),
		"IsAdmin": a.IsAdmin(),
		"isCurModule": func(module int) bool {
			return module == a.curModule
		},
	}
	if len(vars) > 0 {
		t.Merge(vars[0])
	}
	return a.Renderer.Render(tmpl, t)
}

type MainAction struct {
	BaseAction
}

func (a *MainAction) Get() {
	a.Redirect("/user")
}

type LoginAction struct {
	BaseAction
	xsrf.Checker
}

func (a *LoginAction) Get() error {
	return a.Render("login.html", renders.T{
		"XsrfFormHtml": a.Checker.XsrfFormHtml(),
	})
}

func (a *LoginAction) Post() error {
	var user User
	errs := a.Bind(&user)
	if errs.Len() > 0 {
		return errs[0]
	}

	if user.Name == "" || user.Pass == "" {
		return a.Render("login.html", renders.T{
			"msg": "用户名或者密码错误",
		})
	}

	p, err := DB.GetUser(user.Name)
	if err != nil {
		return err
	}
	if p != user.Pass {
		return a.Render("login.html", renders.T{
			"msg": "用户名或者密码错误",
		})
	}

	a.SetLogin(user.Name)
	if a.IsAdmin() {
		a.Redirect("/user")
		return nil
	}
	a.Redirect("/user/chgpass")
	return nil
}

type LogoutAction struct {
	BaseAuthAction
}

func (a *LogoutAction) Get() {
	a.Logout()
	a.Redirect("/")
}
