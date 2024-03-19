# RPC

## Introduction

该项目选择从零实现 Go 语言官方的标准库 net/rpc，并在此基础上，新增了协议交换(protocol exchange)、注册中心(registry)、服务发现(service discovery)、负载均衡(load balance)、超时处理(timeout processing)等特性。  
提供了通过网络或其他I/O连接对一个对象的导出方法的访问。服务端注册一个对象，使它作为一个服务被暴露，服务的名字是该对象的类型名。注册之后，对象的导出方法就可以被远程访问。服务端可以注册多个不同类型的对象（服务），但注册具有相同类型的多个对象是错误的。  

## Example

一个服务端想要导出Arith类型的一个对象：  
```
package server
type Args struct {
    A, B int
}
type Quotient struct {
    Quo, Rem int
}
type Arith int
func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}
func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}
```
服务端会调用（用于HTTP服务）:  
```
arith := new(Arith)
rpc.Register(arith)
rpc.HandleHTTP()
l, e := net.Listen("tcp", ":1234")
if e != nil {
	log.Fatal("listen error:", e)
}
go http.Serve(l, nil)
```
此时，客户端可看到服务"Arith"及它的方法"Arith.Multiply"、"Arith.Divide"。要调用方法，客户端首先呼叫服务端：  
```
client, err := rpc.DialHTTP("tcp", serverAddress + ":1234")
if err != nil {
	log.Fatal("dialing:", err)
}
```
然后，客户端可以执行远程调用：  
```
// Synchronous call
args := &server.Args{7,8}
var reply int
err = client.Call("Arith.Multiply", args, &reply)
if err != nil {
	log.Fatal("arith error:", err)
}
fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
```
或
```
// Asynchronous call
quotient := new(Quotient)
divCall := client.Go("Arith.Divide", args, quotient, nil)
replyCall := <-divCall.Done	// will be equal to divCall
// check errors, print, etc.
```
## Rpc

1. RPC(Remote Procedure Call，远程过程调用)是一种计算机通信协议，允许调用不同进程空间的程序。RPC 的客户端和服务器可以在一台机器上，也可以在不同的机器上。程序员使用时，就像调用本地程序一样，无需关注内部的实现细节。
2. 只有满足如下标准的方法才能用于远程访问，其余方法会被忽略：
- 方法是导出的
- 方法有两个参数，都是导出类型或内建类型
- 方法的第二个参数是指针
- 方法只有一个error接口类型的返回值
3. 方法必须看起来像这样：  
`func (t *T) MethodName(argType T1, replyType *T2) error`  
- 其中T、T1和T2都能被encoding/gob包序列化。这些限制即使使用不同的编解码器也适用。（未来，对定制的编解码器可能会使用较宽松一点的限制）  
- 方法的第一个参数代表调用者提供的参数；第二个参数代表返回给调用者的参数。方法的返回值，如果非nil，将被作为字符串回传，在客户端看来就和errors.New创建的一样。如果返回了错误，回复的参数将不会被发送给客户端。

## Codec

1. 本包默认使用encoding/gob包来传输数据。

## Others