package grpc

import (
	//"context"
	//"fmt"
	//"log"
	//"net"
	//
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"

	"github.com/faelmori/golife/internal"
	//pb "github.com/faelmori/golife/services/proto"
)

type server struct {
	//pb.UnimplementedLifecycleManagerServer
	lifecycleManager internal.LifeCycleManager
}

//func (s *server) StartProcess(ctx context.Context, req *pb.StartProcessRequest) (*pb.StartProcessResponse, error) {
//	process := internal.NewManagedProcess(req.Name, req.Command, req.Args, req.Wait, nil)
//	err := s.lifecycleManager.StartProcess(process)
//	if err != nil {
//		return nil, err
//	}
//	return &pb.StartProcessResponse{Success: true}, nil
//}

//func (s *server) StopProcess(ctx context.Context, req *pb.StopProcessRequest) (*pb.StopProcessResponse, error) {
//	process := s.lifecycleManager.GetProcess(req.Name)
//	if process == nil {
//		return nil, fmt.Errorf("process not found")
//	}
//	err := s.lifecycleManager.StopProcess(process)
//	if err != nil {
//		return nil, err
//	}
//	return &pb.StopProcessResponse{Success: true}, nil
//}

//func (s *server) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
//	status := s.lifecycleManager.Status()
//	return &pb.GetStatusResponse{Status: status}, nil
//}

//func StartGRPCServer(lifecycleManager internal.LifeCycleManager) {
//	lis, err := net.Listen("tcp", ":50051")
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//
//	s := grpc.NewServer()
//	pb.RegisterLifecycleManagerServer(s, &server{lifecycleManager: lifecycleManager})
//	reflection.Register(s)
//
//	log.Println("gRPC server is running on port 50051")
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}
