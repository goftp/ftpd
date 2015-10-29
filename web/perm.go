package web

import (
	"net/url"
	"os"
	"path"

	"github.com/go-xweb/xweb"
	"github.com/goftp/server"
)

type PermAction struct {
	BaseAction

	get         xweb.Mapper `xweb:"/"`
	add         xweb.Mapper
	edit        xweb.Mapper
	del         xweb.Mapper
	updateOwner xweb.Mapper
	updateGroup xweb.Mapper
	updatePerm  xweb.Mapper
}

func (c *PermAction) Init() {
	c.AddTmplVar("isCurModule", c.IsCurModule)
}

func (c *PermAction) IsCurModule(module int) bool {
	return PERM_MODULE == module
}

func hasPerm(mode os.FileMode, idx int, rOrW string) bool {
	return string(mode.String()[idx]) == rOrW
}

func (c *PermAction) Get() error {
	p := c.GetString("path")
	var pathinfos = make([]server.FileInfo, 0)
	var err error
	var parent string
	if p != "" {
		p, err = url.QueryUnescape(p)
		if err != nil {
			return err
		}
		parent = path.Dir(p)
	} else {
		parent = "/"
		p = "/"
	}
	driver, err := Factory.NewDriver()
	if err != nil {
		return err
	}
	pathinfos, err = driver.DirContents(p)
	if err != nil {
		return err
	}

	var users []User
	err = DB.UserList(&users)
	if err != nil {
		return err
	}
	var groups []string
	err = DB.GroupList(&groups)
	if err != nil {
		return err
	}
	return c.Render("perm/list.html", &xweb.T{
		"parent":  parent,
		"path":    p,
		"infos":   pathinfos,
		"hasPerm": hasPerm,
		"users":   users,
		"groups":  groups,
	})
}

func (c *PermAction) UpdateGroup() {
	name := c.GetString("name")
	newgroup := c.GetString("newgroup")

	if name == "" || newgroup == "" {
		c.ServeJson(map[string]string{"status": "0", "error": "empty params"})
		return
	}

	name, err := url.QueryUnescape(name)
	if err != nil {
		c.ServeJson(map[string]string{"status": "0", "error": err.Error()})
		return
	}

	err = Perm.ChGroup(name, newgroup)
	if err != nil {
		c.ServeJson(map[string]string{"status": "0", "error": err.Error()})
		return
	}
	c.ServeJson(map[string]string{"status": "1"})
}

func (c *PermAction) UpdateOwner() {
	name := c.GetString("name")
	newowner := c.GetString("newowner")

	if name == "" || newowner == "" {
		c.ServeJson(map[string]string{"status": "0", "error": "empty params"})
		return
	}

	name, err := url.QueryUnescape(name)
	if err != nil {
		c.ServeJson(map[string]string{"status": "0", "error": err.Error()})
		return
	}

	err = Perm.ChOwner(name, newowner)
	if err != nil {
		c.ServeJson(map[string]string{"status": "0", "error": err.Error()})
		return
	}
	c.ServeJson(map[string]string{"status": "1"})
}

func (c *PermAction) UpdatePerm() {
	name := c.GetString("name")
	typ := c.GetString("typ")
	right := c.GetString("right")
	has := c.GetString("has")

	if name == "" || typ == "" || right == "" || has == "" {
		c.ServeJson(map[string]string{"status": "0", "error": "empty params"})
		return
	}

	name, err := url.QueryUnescape(name)
	if err != nil {
		c.ServeJson(map[string]string{"status": "0", "error": err.Error()})
		return
	}

	mode, err := Perm.GetMode(name)
	if err != nil {
		c.ServeJson(map[string]string{"status": "0", "error": err.Error()})
		return
	}

	var bs uint = 0
	if typ == "owner" {
		if right == "r" {
			bs = 8
		} else if right == "w" {
			bs = 7
		}
	} else if typ == "group" {
		if right == "r" {
			bs = 5
		} else if right == "w" {
			bs = 4
		}
	} else if typ == "other" {
		if right == "r" {
			bs = 2
		} else if right == "w" {
			bs = 1
		}
	}

	if has == "true" {
		mode = os.FileMode(uint32(mode) + uint32(1<<bs))
	} else {
		mode = os.FileMode(uint32(mode) - uint32(1<<bs))
	}

	err = Perm.ChMode(name, mode)
	if err != nil {
		c.ServeJson(map[string]string{"status": "0", "error": err.Error()})
		return
	}
	c.ServeJson(map[string]string{"status": "1"})
}
