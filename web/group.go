package web

import (
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/xsrf"
)

type GroupBaseAction struct {
	BaseAuthAction
}

func (c *GroupBaseAction) Before() {
	c.curModule = GROUP_MODULE
}

type GroupAction struct {
	GroupBaseAction
}

func (c *GroupAction) Get() error {
	groups := make([]string, 0)
	err := DB.GroupList(&groups)
	if err != nil {
		return err
	}
	return c.Render("group/list.html", renders.T{
		"groups": groups,
		"admin":  adminUser,
	})
}

type GroupAddAction struct {
	GroupBaseAction
	xsrf.Checker
}

func (c *GroupAddAction) Get() error {
	return c.Render("group/add.html", renders.T{
		"XsrfFormHtml": c.Checker.XsrfFormHtml(),
	})
}

func (c *GroupAddAction) Post() error {
	name := c.Form("name")
	err := DB.AddGroup(name)
	if err != nil {
		return err
	}
	c.Redirect("/group")
	return nil
}

type GroupEditAction struct {
	GroupBaseAction
}

func (c *GroupEditAction) Get() error {
	name := c.Form("name")
	if name == "" {
		c.Redirect("/group")
		return nil
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
	return c.Render("group/edit.html", renders.T{
		"selUsers":   selUsers,
		"otherUsers": otherUsers,
	})
}

type GroupDelAction struct {
	GroupBaseAction
}

func (c *GroupDelAction) Get() error {
	name := c.Form("name")
	err := DB.DelGroup(name)
	if err != nil {
		return err
	}
	c.Redirect("/group")
	return nil
}
