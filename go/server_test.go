package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"testing"

	"bou.ke/monkey"
	mk "github.com/mikayelabrahamyan/task/go/marketplace-gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func startTestServer() *bufconn.Listener {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	mk.RegisterCreatorsServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	return lis
}

func TestGetSortedCreators(t *testing.T) {
	lis := startTestServer()
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := mk.NewCreatorsServiceClient(conn)

	// wrap for mocking os.ReadFile by empty data
	func() {
		data := Data{
			Creators: []mk.Creator{},
			Products: []mk.Product{},
		}
		dataBytes, _ := json.Marshal(data)

		guard := monkey.Patch(os.ReadFile, func(filename string) ([]byte, error) {
			return dataBytes, nil
		})
		defer guard.Unpatch()

		// wrap limit test
		{
			resp, err := client.GetSortedCreators(ctx, &mk.GetSortedCreatorsRequest{Limit: 1})
			if err != nil {
				t.Fatalf("GetSortedCreators failed: %v", err)
			}

			if len(resp.Creators) != 0 {
				t.Errorf("Expected 0 creators while provided limit 1, got %d", len(resp.Creators))
			}
		}
	}()

	// wrap for mocking os.ReadFile by mock data
	func() {
		// id 1 has more products than others,
		// id 2 and 3 has equal products => need to sort base on CreateTime id 2 has more fresh product
		// id 4 has less products from 1,2,3
		// id 5 has not any products

		// need to be 1, 2, 3, 4, 5 in DESCENDING order
		// need to be 5, 4, 3, 2, 1 in ASCENDING order
		data := Data{
			Creators: []mk.Creator{
				{Id: "1", Email: "creator-one@gmail.com"},
				{Id: "2", Email: "creator-two@gmail.com"},
				{Id: "3", Email: "creator-three@gmail.com"},
				{Id: "4", Email: "creator-four@gmail.com"},
				{Id: "5", Email: "creator-five@gmail.com"},
			},
			Products: []mk.Product{
				{Id: "p1", CreatorId: "1", CreateTime: "2000-04-04T09:05:59.752712+02:00"},
				{Id: "p2", CreatorId: "1", CreateTime: "2000-04-04T09:05:59.752712+02:00"},
				{Id: "p3", CreatorId: "1", CreateTime: "2000-04-04T09:05:59.752712+02:00"},

				{Id: "p4", CreatorId: "2", CreateTime: "2023-04-04T09:05:59.752712+02:00"},
				{Id: "p5", CreatorId: "2", CreateTime: "2022-04-04T09:05:59.752712+02:00"},

				{Id: "p6", CreatorId: "3", CreateTime: "2021-04-04T09:05:59.752712+02:00"},
				{Id: "p7", CreatorId: "3", CreateTime: "2020-04-04T09:05:59.752712+02:00"},

				{Id: "p8", CreatorId: "4", CreateTime: ""},
			},
		}
		dataBytes, _ := json.Marshal(data)

		guard := monkey.Patch(os.ReadFile, func(filename string) ([]byte, error) {
			return dataBytes, nil
		})
		defer guard.Unpatch()

		// wrap limit test
		{
			resp, err := client.GetSortedCreators(ctx, &mk.GetSortedCreatorsRequest{Limit: 1, Order: mk.SortOrder_ASCENDING})
			if err != nil {
				t.Fatalf("GetSortedCreators failed: %v", err)
			}

			if len(resp.Creators) != 1 {
				t.Errorf("Expected 1 creators while provided limit 1, got %d", len(resp.Creators))
			}
		}

		// wrap SortOrder ASCENDING test
		{
			resp, err := client.GetSortedCreators(ctx, &mk.GetSortedCreatorsRequest{Limit: 1000, Order: mk.SortOrder_ASCENDING})

			if err != nil {
				t.Fatalf("GetSortedCreators failed: %v", err)
			}

			if len(resp.Creators) != 5 {
				t.Errorf("Expected 5 creators, got %d", len(resp.Creators))
			}

			if !(resp.Creators[0].Id == "5" &&
				resp.Creators[1].Id == "4" &&
				resp.Creators[2].Id == "3" &&
				resp.Creators[3].Id == "2" &&
				resp.Creators[4].Id == "1") {
				t.Errorf("Creators are not sorted by ASCENDING order correctly")
			}
		}

		// wrap SortOrder DESCENDING test
		{
			resp, err := client.GetSortedCreators(ctx, &mk.GetSortedCreatorsRequest{Limit: 1000, Order: mk.SortOrder_DESCENDING})

			if err != nil {
				t.Fatalf("GetSortedCreators failed: %v", err)
			}

			if len(resp.Creators) != 5 {
				t.Errorf("Expected 5 creators, got %d", len(resp.Creators))
			}

			if !(resp.Creators[0].Id == "1" &&
				resp.Creators[1].Id == "2" &&
				resp.Creators[2].Id == "3" &&
				resp.Creators[3].Id == "4" &&
				resp.Creators[4].Id == "5") {
				t.Errorf("Creators are not sorted by DESCENDING order correctly")
			}
		}
	}()

}
