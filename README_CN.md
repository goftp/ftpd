# ftpd

[English](README.md)

这是一个基于 [github.com/goftp/server](http://github.com/goftp/server) 编写的Ftp服务器程序。

文档可以通过 [godoc](http://godoc.org/github.com/goftp/ftpd) 获取。

## 安装

    go get github.com/goftp/ftpd

然后运行

    $GOPATH/bin/ftpd

最后，通过FTP客户端连接即可：

    host: 127.0.0.1
    port: 2121
    username: admin
    password: 123456

如需要进一步修改，可以拷贝config.ini文件到ftpd目录下，然后修改其中的配置