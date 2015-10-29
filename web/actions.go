package web

import "github.com/go-xweb/xweb"

type BaseAction struct {
	*xweb.Action
}

func (a *BaseAction) IsLogin() bool {
	s := a.GetSession("userId")
	return s != nil
}

func (a *BaseAction) IsAdmin() bool {
	s := a.GetSession("userId")
	return s != nil && s.(string) == adminUser
}

func (a *BaseAction) Init() {
	a.AddTmplVars(&xweb.T{
		"IsLogin": a.IsLogin,
		"IsAdmin": a.IsAdmin,
	})
}

type MainAction struct {
	BaseAction

	get    xweb.Mapper `xweb:"/"`
	login  xweb.Mapper
	logout xweb.Mapper
}

func (a *MainAction) Get() error {
	return a.Go("get", &UserAction{})
}

func (a *MainAction) Logout() {
	a.DelSession("userId")
	a.Redirect("/")
}

func (a *MainAction) Login() error {
	if a.Method() == "GET" {
		return a.Render("login.html")
	} else if a.Method() == "POST" {
		user := new(User)
		err := a.MapForm(user, "")
		if err != nil {
			return err
		}

		if user.Name == "" || user.Pass == "" {
			return a.Render("login.html", &xweb.T{
				"msg": "用户名或者密码错误",
			})
		}

		p, err := DB.GetUser(user.Name)
		if err != nil {
			return err
		}
		if p != user.Pass {
			return a.Render("login.html", &xweb.T{
				"msg": "用户名或者密码错误",
			})
		}

		a.SetSession("userId", user.Name)
		if a.IsAdmin() {
			return a.Go("get", &UserAction{})
		}
		return a.Go("chgpass", &UserAction{})
	}
	return xweb.NotSupported()
}
