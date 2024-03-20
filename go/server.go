package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"time"

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

	type Pair struct {
		First  int
		Second time.Time
	}

	productCount := make(map[string]Pair)
	layout := "2006-01-02T15:04:05.000000-07:00"
	for _, product := range data.Products {
		timestamp, timeErr := time.Parse(layout, product.CreateTime)
		if timeErr != nil {
			fmt.Printf("timeErr: %+v\n", timeErr)
			timestamp = time.Time{}
		}
		if pair, exists := productCount[product.CreatorId]; exists {
			pair.First++
			if timestamp.After(productCount[product.CreatorId].Second) {
				pair.Second = timestamp
			}
			productCount[product.CreatorId] = pair
		} else {
			productCount[product.CreatorId] = Pair{First: 1, Second: timestamp}
		}
	}
	sort.Slice(data.Creators, func(i, j int) bool {
		result := productCount[data.Creators[i].Id].First < productCount[data.Creators[j].Id].First

		if productCount[data.Creators[i].Id].First == productCount[data.Creators[j].Id].First {
			result = productCount[data.Creators[i].Id].Second.Before(productCount[data.Creators[j].Id].Second)
		}

		if in.Order == mk.SortOrder_ASCENDING {
			return result
		}
		return !result
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
