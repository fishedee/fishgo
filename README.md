# fishgo

## 问题

由于最新的go版本(1.5)仍然没有好好地解决包管理的问题

缺少包依赖声明，缺少依赖包的版本声明

经常导致dev与idc的依赖包不一致导致编译有问题

目前解决办法是利用git的版本来确定各依赖包的版本号

以保证dev与idc的依赖包总是保持一致的

## 安装

```
git clone git@github.com:fishedee/fishgo.git --recursive
export GOPATH=xxxx
cd fishgo
./install.sh
```

## 功能

* golang.org/x/net
* beego 1.5
* xorm 0.4.4
* beego 的扩展功能

## 使用

（有空再补充吧，有点懒）
