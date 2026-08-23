package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "gotest/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedProductServiceServer
	mu       sync.RWMutex
	products map[uint32]*pb.ProductResponse
	nextID   uint32
}

// 1. Unary RPC: GetProduct
func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, exists := s.products[req.GetId()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "produk dengan ID %d tidak ditemukan", req.GetId())
	}

	return product, nil
}

// 2. Unary RPC: CreateProduct
func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	if req.GetName() == "" || req.GetPrice() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nama dan harga produk harus valid")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	newProduct := &pb.ProductResponse{
		Id:    id,
		Name:  req.GetName(),
		Price: req.GetPrice(),
		Stock: req.GetStock(),
	}

	s.products[id] = newProduct
	fmt.Printf("[gRPC SERVER] 📦 Produk baru dibuat: ID=%d, Nama=%s\n", id, req.GetName())
	return newProduct, nil
}

// 3. Server Streaming RPC: ListProducts (Mengirim data secara streaming ke client)
func (s *server) ListProducts(req *pb.ListProductsRequest, stream grpc.ServerStreamingServer[pb.ProductResponse]) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("[gRPC SERVER] 📡 Streaming data produk ke client...")
	for _, p := range s.products {
		if err := stream.Send(p); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond) // Simulasi latency streaming
	}
	return nil
}

func main() {
	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Gagal membuka port: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := &server{
		products: map[uint32]*pb.ProductResponse{
			1: {Id: 1, Name: "MacBook Air M2", Price: 17500000, Stock: 8},
			2: {Id: 2, Name: "Sony WH-1000XM5", Price: 4999000, Stock: 12},
		},
		nextID: 3,
	}

	pb.RegisterProductServiceServer(grpcServer, s)

	fmt.Printf("🚀 gRPC Server berjalan di port %s (TCP)\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC Server gagal berjalan: %v", err)
	}
}
