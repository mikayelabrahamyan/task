package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"sort"

	mk "github.com/mikayelabrahamyan/task/go/marketplace-gen"
	"google.golang.org/grpc"
)

type server struct {
	mk.UnimplementedCreatorsServiceServer
	mk.UnimplementedProductsServiceServer
}

type Data struct {
	Creators []mk.Creator `json:"Creators"`
	Products []mk.Product `json:"Products"`
}

func GetCreatorById(creators []mk.Creator, id string) (*mk.Creator, bool) {
	for _, creator := range creators {
		if creator.Id == id {
			return &creator, true
		}
	}
	return nil, false
}

func GetProductById(products []mk.Product, id string) (*mk.Product, bool) {
	for _, product := range products {
		if product.Id == id {
			return &product, true
		}
	}
	return nil, false
}

func (s *server) GetSortedCreators(ctx context.Context, in *mk.GetSortedCreatorsRequest) (*mk.GetSortedCreatorsResponse, error) {
	file, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var data Data
	json.Unmarshal(file, &data)

	productCount := make(map[string]int)
	for _, product := range data.Products {
		productCount[product.CreatorId]++
	}

	sort.Slice(data.Creators, func(i, j int) bool {
		if in.Order == mk.SortOrder_ASCENDING {
			return productCount[data.Creators[i].Id] < productCount[data.Creators[j].Id]
		}
		return productCount[data.Creators[i].Id] > productCount[data.Creators[j].Id]
	})

	var sortedCreators []*mk.Creator
	for i := range data.Creators {
		if i >= int(in.Limit) {
			break
		}
		sortedCreators = append(sortedCreators, &data.Creators[i])
	}

	return &mk.GetSortedCreatorsResponse{Creators: sortedCreators}, nil
}

func (s *server) GetCreator(ctx context.Context, in *mk.GetCreatorRequest) (*mk.GetCreatorResponse, error) {
	file, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var data Data
	json.Unmarshal(file, &data)

	creator, exist := GetCreatorById(data.Creators, in.Id)

	if exist {
		return &mk.GetCreatorResponse{Creator: creator}, nil
	} else {
		return &mk.GetCreatorResponse{Creator: &mk.Creator{}}, nil
	}
}

func (s *server) GetCreators(ctx context.Context, in *mk.GetCreatorsRequest) (*mk.GetCreatorsResponse, error) {
	file, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var data Data
	json.Unmarshal(file, &data)

	var filteredCreators []*mk.Creator
	for _, creator := range data.Creators {
		filteredCreators = append(filteredCreators, &creator)
	}

	return &mk.GetCreatorsResponse{Creators: filteredCreators}, nil
}

func (s *server) GetProducts(ctx context.Context, in *mk.GetProductsRequest) (*mk.GetProductsResponse, error) {
	file, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var data Data
	json.Unmarshal(file, &data)

	var filteredProducts []*mk.Product
	for _, product := range data.Products {
		filteredProducts = append(filteredProducts, &product)
	}

	return &mk.GetProductsResponse{Products: filteredProducts}, nil
}

func (s *server) GetProduct(ctx context.Context, in *mk.GetProductRequest) (*mk.GetProductResponse, error) {
	file, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	var data Data
	json.Unmarshal(file, &data)

	product, exist := GetProductById(data.Products, in.Id)

	if exist {
		return &mk.GetProductResponse{Product: product}, nil
	} else {
		return &mk.GetProductResponse{Product: &mk.Product{}}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	mk.RegisterCreatorsServiceServer(s, &server{})
	mk.RegisterProductsServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
