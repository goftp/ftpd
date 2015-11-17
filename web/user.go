package web

import (
	"github.com/tango-contrib/flash"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
)

type UserBaseAction struct {
	BaseAuthAction
}

func (c *UserBaseAction) Before() {
	c.curModule = USER_MODULE
}

type UserAction struct {
	UserBaseAction
}

func (a *UserAction) Get() error {
	users := make([]User, 0)
	err := DB.UserList(&users)
	if err != nil {
		return err
	}

	return a.Render("user/list.html", renders.T{
		"users": users,
		"admin": adminUser,
	})
}

type ChgPassAction struct {
	UserBaseAction
	xsrf.Checker
	flash.Flash
}

func (c *ChgPassAction) Before() {
	c.curModule = CHGPASS_MODULE
}

func (c *ChgPassAction) Get() error {
	user := c.LoginUserId()
	_, err := DB.GetUser(user)
	if err != nil {
		return err
	}

	return c.Render("user/chgpass.html", renders.T{
		"user":         user,
		"userId":       c.LoginUserId(),
		"XsrfFormHtml": c.Checker.XsrfFormHtml(),
		"Flash":        c.Flash.Data(),
	})
}

func (c *ChgPassAction) Post() error {
	var user User
	errs := c.Bind(&user)
	if errs.Len() > 0 {
		return errs[0]
	}
	err := DB.ChgPass(user.Name, user.Pass)
	if err != nil {
		return err
	}
	c.Flash.Set("info", "修改密码成功")
	c.Redirect("/user/chgpass")
	return nil
}

type UserAddAction struct {
	UserBaseAction
	xsrf.Checker
}

func (a *UserAddAction) Get() error {
	return a.Render("user/add.html", renders.T{
		"XsrfFormHtml": a.Checker.XsrfFormHtml(),
	})
}

func (a *UserAddAction) Post() error {
	var user User
	errs := a.Bind(&user)
	if errs.Len() > 0 {
		return errs[0]
	}
	err := DB.AddUser(user.Name, user.Pass)
	if err != nil {
		return err
	}

	a.Redirect("/user")
	return nil
}

type UserEditAction struct {
	UserBaseAction
	xsrf.Checker
}

func (a *UserEditAction) Get() error {
	name := a.Form("name")
	pass, err := DB.GetUser(name)
	if err != nil {
		return err
	}

	return a.Render("user/edit.html", renders.T{
		"user":         &User{name, pass},
		"XsrfFormHtml": a.Checker.XsrfFormHtml(),
	})
}

func (a *UserEditAction) Post() error {
	var user User
	errs := a.Bind(&user)
	if errs.Len() > 0 {
		return errs[0]
	}
	err := DB.ChgPass(user.Name, user.Pass)
	if err != nil {
		return err
	}
	a.Redirect("/user")
	return nil
}

type UserDelAction struct {
	UserBaseAction
}

func (a *UserDelAction) Get() error {
	name := a.Form("name")
	err := DB.DelUser(name)
	if err != nil {
		return err
	}

	a.Redirect("/user")
	return nil
}
