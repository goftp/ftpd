package web

import "github.com/go-xweb/xweb"

type UserAction struct {
	BaseAction

	get     xweb.Mapper `xweb:"/"`
	add     xweb.Mapper
	edit    xweb.Mapper
	del     xweb.Mapper
	chgpass xweb.Mapper

	module int
}

func (c *UserAction) Init() {
	c.module = USER_MODULE
	c.AddTmplVar("isCurModule", c.IsCurModule)
}

func (c *UserAction) IsCurModule(module int) bool {
	return c.module == module
}

func (c *UserAction) Chgpass() error {
	c.module = CHGPASS_MODULE
	if c.Method() == "GET" {
		user := c.GetSession("userId")
		_, err := DB.GetUser(user.(string))
		if err != nil {
			return err
		}

		return c.Render("user/chgpass.html", &xweb.T{
			"user": user,
		})
	} else if c.Method() == "POST" {
		var user User
		err := c.MapForm(&user, "")
		if err != nil {
			return err
		}
		err = DB.ChgPass(user.Name, user.Pass)
		if err != nil {
			return err
		}
		return c.Go("chgpass")
	}

	return xweb.NotSupported()
}

func (a *UserAction) Get() error {
	users := make([]User, 0)
	err := DB.UserList(&users)
	if err != nil {
		return err
	}

	return a.Render("user/list.html", &xweb.T{
		"users": users,
		"admin": adminUser,
	})
}

func (a *UserAction) Add() error {
	if a.Method() == "GET" {
		return a.Render("user/add.html")
	} else if a.Method() == "POST" {
		user := new(User)
		err := a.MapForm(user, "")
		if err != nil {
			return err
		}
		err = DB.AddUser(user.Name, user.Pass)
		if err != nil {
			return err
		}
		return a.Go("get")
	}
	return xweb.NotSupported()
}

func (a *UserAction) Edit() error {
	if a.Method() == "GET" {
		name := a.GetString("name")
		pass, err := DB.GetUser(name)
		if err != nil {
			return err
		}

		return a.Render("user/edit.html", &xweb.T{
			"user": &User{name, pass},
		})
	} else if a.Method() == "POST" {
		user := new(User)
		err := a.MapForm(user, "")
		if err != nil {
			return err
		}
		err = DB.ChgPass(user.Name, user.Pass)
		if err != nil {
			return err
		}
		return a.Go("get")
	}
	return xweb.NotSupported()
}

func (a *UserAction) Del() error {
	name := a.GetString("name")
	err := DB.DelUser(name)
	if err != nil {
		return err
	}

	return a.Go("get")
}
