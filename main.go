package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Unknwon/goconfig"
	"github.com/goftp/file-driver"
	"github.com/goftp/ftpd/web"
	"github.com/goftp/leveldb-auth"
	"github.com/goftp/leveldb-perm"
	"github.com/goftp/qiniu-driver"
	"github.com/goftp/server"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	cfg *goconfig.ConfigFile
)

func main() {
	var err error
	cfg, err = goconfig.LoadConfigFile("config.ini", "custom.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	port, _ := cfg.Int("server", "port")
	db, err := leveldb.OpenFile("./authperm.db", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var auth = &ldbauth.LDBAuth{db}
	/*if cfg.MustValue("auth", "type") == "xorm" {
		panic("current is not supported yet")
		//auth = xormauth.NewXormAuth(orm *xorm.Engine, allowAnony bool, perm os.FileMode)
	} else {

	}*/

	var perm server.Perm
	if cfg.MustValue("perm", "type") == "leveldb" {
		perm = ldbperm.NewLDBPerm(db, "root", "root", os.ModePerm)
	} else {
		perm = server.NewSimplePerm("root", "root")
	}

	typ, _ := cfg.GetValue("driver", "type")
	var factory server.DriverFactory
	if typ == "file" {
		rootPath, _ := cfg.GetValue("file", "rootpath")
		_, err = os.Lstat(rootPath)
		if os.IsNotExist(err) {
			os.MkdirAll(rootPath, os.ModePerm)
		} else if err != nil {
			fmt.Println(err)
			return
		}
		factory = &filedriver.FileDriverFactory{
			rootPath,
			perm,
		}
	} else if typ == "qiniu" {
		accessKey, _ := cfg.GetValue("qiniu", "accessKey")
		secretKey, _ := cfg.GetValue("qiniu", "secretKey")
		bucket, _ := cfg.GetValue("qiniu", "bucket")
		factory = qiniudriver.NewQiniuDriverFactory(accessKey,
			secretKey, bucket)
	} else {
		log.Fatal("no driver type input")
	}

	// start web manage UI
	useweb, _ := cfg.Bool("web", "enable")
	if useweb {
		web.DB = auth
		web.Perm = perm
		web.Factory = factory
		weblisten, _ := cfg.GetValue("web", "listen")
		admin, _ := cfg.GetValue("admin", "user")
		pass, _ := cfg.GetValue("admin", "pass")
		ssl, _ := cfg.Bool("web", "ssl")

		go web.Web(weblisten, "static", "templates", admin, pass, ssl)
	}

	ftpName, _ := cfg.GetValue("server", "name")
	opt := &server.ServerOpts{
		Name:    ftpName,
		Factory: factory,
		Port:    port,
		Auth:    auth,
	}

	// start ftp server
	ftpServer := server.NewServer(opt)
	err = ftpServer.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
