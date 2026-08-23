package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "gotest/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Hubungkan client ke gRPC Server via TCP port :50051
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Gagal koneksi ke gRPC Server: %v", err)
	}
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("=== 1. UNARY RPC: CreateProduct ===")
	newProd, err := client.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:  "iPad Pro M4",
		Price: 19500000,
		Stock: 5,
	})
	if err != nil {
		log.Printf("Gagal buat produk: %v", err)
	} else {
		fmt.Printf("✅ Produk dibuat -> ID: %d, Nama: %s, Harga: Rp%.2f\n",
			newProd.GetId(), newProd.GetName(), newProd.GetPrice())
	}

	fmt.Println("\n=== 2. UNARY RPC: GetProduct ===")
	prod1, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: 1})
	if err != nil {
		log.Printf("Gagal ambil produk ID 1: %v", err)
	} else {
		fmt.Printf("✅ Detail Produk ID 1 -> %s (Rp%.2f, Stok: %d)\n",
			prod1.GetName(), prod1.GetPrice(), prod1.GetStock())
	}

	fmt.Println("\n=== 3. SERVER STREAMING RPC: ListProducts ===")
	stream, err := client.ListProducts(ctx, &pb.ListProductsRequest{})
	if err != nil {
		log.Fatalf("Gagal inisialisasi streaming: %v", err)
	}

	for {
		product, err := stream.Recv()
		if err == io.EOF {
			// Selesai membaca semua data stream dari server
			fmt.Println("📡 Streaming selesai!")
			break
		}
		if err != nil {
			log.Fatalf("Error saat membaca stream: %v", err)
		}
		fmt.Printf("📦 [STREAM ITEM] ID=%d | Nama=%s | Harga=Rp%.2f | Stok=%d\n",
			product.GetId(), product.GetName(), product.GetPrice(), product.GetStock())
	}
}
