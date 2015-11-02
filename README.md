# ftpd

[中文](README_CN.md)

A FTP server based on [github.com/goftp/server](http://github.com/goftp/server).

Full documentation for the package is available on [godoc](http://godoc.org/github.com/goftp/ftpd)

## Installation

    go get github.com/goftp/ftpd

Then run it:

    $GOPATH/bin/ftpd

And finally, connect to the server with any FTP client and the following
details:

    host: 127.0.0.1
    port: 2121
    username: admin
    password: 123456

More features, you can copy config.ini to the ftpd directory and modify it.