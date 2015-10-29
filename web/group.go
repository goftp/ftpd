package web

import "github.com/go-xweb/xweb"

type GroupAction struct {
	BaseAction

	get  xweb.Mapper `xweb:"/"`
	add  xweb.Mapper
	edit xweb.Mapper
	del  xweb.Mapper
}

func (c *GroupAction) Init() {
	c.AddTmplVar("isCurModule", c.IsCurModule)
}

func (c *GroupAction) IsCurModule(module int) bool {
	return GROUP_MODULE == module
}

func (c *GroupAction) Get() error {
	groups := make([]string, 0)
	err := DB.GroupList(&groups)
	if err != nil {
		return err
	}
	return c.Render("group/list.html", &xweb.T{
		"groups": groups,
		"admin":  adminUser,
	})
}

func (c *GroupAction) Add() error {
	if c.Method() == "GET" {
		return c.Render("group/add.html")
	} else if c.Method() == "POST" {
		name := c.GetString("name")
		err := DB.AddGroup(name)
		if err != nil {
			return err
		}
		return c.Go("get")
	}
	return xweb.NotSupported()
}

func (c *GroupAction) Edit() error {
	name := c.GetString("name")
	if name == "" {
		return c.Go("get")
	}
	var selUsers []string
	err := DB.GroupUser(name, &selUsers)
	if err != nil {
		return err
	}
	var users []User
	err = DB.UserList(&users)
	if err != nil {
		return err
	}
	var otherUsers = make([]string, 0, len(users)-len(selUsers))
	for _, user := range users {
		var hasUser bool
		for _, name := range selUsers {
			if user.Name == name {
				hasUser = true
				break
			}
		}
		if !hasUser {
			otherUsers = append(otherUsers, user.Name)
		}
	}
	return c.Render("group/edit.html", &xweb.T{
		"selUsers":   selUsers,
		"otherUsers": otherUsers,
	})
}

func (c *GroupAction) Del() error {
	name := c.GetString("name")
	err := DB.DelGroup(name)
	if err != nil {
		return err
	}
	return c.Go("get")
}
