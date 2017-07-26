# dagger

一个Go写的TCP通信框架

## 说明

写这个包的目的是提供一个最精简的TCP网络层，隔离业务层逻辑，只做基本的连接管理。

基本的实现思路是这样，每来一个conn，`Server`就开一个goroutine来处理它，这个goroutine里跑的`Session`就是对`net.Conn`的封装，`Session`里面维护了三个goroutine，来对应接收、发送和处理的逻辑。另外定义了`SessionCallback`和`Protocol`两个接口，用来解耦用户的上层业务逻辑和拆包封包的协议。

## TODO

例子和使用说明