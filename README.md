session [![Build Status](https://travis-ci.org/issue9/session.svg?branch=master)](https://travis-ci.org/issue9/session)
======

```go
opt := session.NewOptions(stores.NewMemory(), ...)

// 在每一个Handler中调用Start()开始一个Session操作。
sess,err :=session.Start(opt, w, req)

opt.Close()
```


### 安装

```shell
go get github.com/issue9/session
```


### 文档

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/issue9/session)
[![GoDoc](https://godoc.org/github.com/issue9/session?status.svg)](https://godoc.org/github.com/issue9/session)


### 版权

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/issue9/session/blob/master/LICENSE)
