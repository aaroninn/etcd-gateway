package main

import (
	"context"
	"log"
	"net/http"

	rpcpb "etcdtest/proto/rpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime" // 注意v2版本
	"google.golang.org/grpc"
	insec "google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到etcd
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:2379",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insec.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial etcd:", err)
	}
	log.Println("dial to etcd succeed")

	gwmux := runtime.NewServeMux()
	// 注册etcd kv服务
	err = rpcpb.RegisterKVHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8080",
		Handler: gwmux,
	}
	// 提供gRPC-Gateway服务
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
