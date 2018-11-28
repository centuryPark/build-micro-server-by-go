package main

import (
	pb "github.com/micro-shippy/consignment-service/proto/consignment"
	vesselPb "github.com/micro-shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
	"os"
)

const (
	PORT = ":50051"
	DEFAULT_HOST = "localhost:27017"
)

func main() {
	// 获取容器设置的环境变量的值 - 数据库地址
	dbHost := os.Getenv("DB_HOST")
	if dbHost == ""{
		dbHost = DEFAULT_HOST
	}
	session, err := CreateSession(dbHost)
	// 创建于 MongoDB 的主会话，需在退出 main() 时候手动释放连接
	defer session.Close()
	if err != nil {
		log.Fatalf("create session error: %v\n", err)
	}

	server := micro.NewService(
		// 必须和 consignment.proto 中的 package 一致
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	// 解析命令行参数
	server.Init()
	// 作为 vessel-service 的客户端，通过rpc调用vessel-service服务
	vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client())
	// 将 server 作为微服务的服务端
	pb.RegisterShippingServiceHandler(server.Server(), &handler{session, vClient})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
