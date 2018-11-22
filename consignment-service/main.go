package main

import (
	"context"
	pb "github.com/micro-shippy/consignment-service/proto/consignment"
	vesselPb "github.com/micro-shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
)

const (
	PORT = ":50051"
)

// 定义货轮接口
type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error) // 存放新货物
	GetAll() []*pb.Consignment                                   // 获取仓库中所有的货物
}

// 存放多批货物的仓库，实现货轮接口
type Repository struct {
	consignments []*pb.Consignment
}

// 把货物存入货轮中
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

// 获取货物列表
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// 服务
type service struct {
	repo Repository
	// consignment-service 作为客户端调用 vessel-service 的函数
	vesselClient vesselPb.VesselServiceClient
}

// 实现 consignment.pb.go 中的 ShippingServiceHandler 接口
// 使 service 作为 gRPC 的服务端

// 托运新的货物
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) (error) {

	// 检查是否有适合的货轮
	vReq := &vesselPb.Specification{
		Capacity:  int32(len(req.Containers)),
		MaxWeight: req.Weight,
	}
	vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}

	// 货物被承运
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id

	// 接收承运的货物
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}

// 获取目前所有托运的货物
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	allConsignments := s.repo.GetAll()
	res.Consignments = allConsignments
	return nil
}

func main() {
	server := micro.NewService(
		// 必须和 consignment.proto 中的 package 一致
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	// 解析命令行参数
	server.Init()
	repo := Repository{}

	// 作为 vessel-service 的客户端
	vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client())

	// 向 rRPC 服务器注册微服务
	// 此时会把我们自己实现的微服务 service 与协议中的 ShippingServiceServer 绑定
	pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vClient})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
