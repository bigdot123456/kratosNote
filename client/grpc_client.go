package main

import (
	pb "callServer/api"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	Address = "127.0.0.1:50052"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:9000", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	c := pb.NewDemoClient(conn)

	req := new(pb.HelloReq)
	req.Name = "kratos grpc"
	r, err := c.SayHelloURL(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r.Content)
	fmt.Println("Go another!")

	//req1 := new(pb.HelloReq)
	//req1.Name = "grpc goofly" // If we don't provide req1 name, it will be error, since service will check name field!
	//r1, err1 := c.Create(context.Background(), req1)
	//if err1 != nil {
	//	fmt.Println(err1)
	//	return
	//}
	//fmt.Println(r1.Content)

}
