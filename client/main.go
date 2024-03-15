package main

import (
	"context"
	"fmt"
	"log"

	"test/user"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50054"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)

	
	createUserResp, err := client.CreateUser(context.Background(), &user.UserRequest{Name: "John Doe", Email: "john@example.com"})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	fmt.Printf("Created user: %+v\n", createUserResp)

	// Get user
	getUserResp, err := client.GetUser(context.Background(), &user.UserRequest{Id: createUserResp.Id})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	fmt.Printf("Retrieved user: %+v\n", getUserResp)

	// Update user
	updateUserResp, err := client.UpdateUser(context.Background(), &user.UserRequest{Id: createUserResp.Id, Name: "Updated Name"})
	if err != nil {
		log.Fatalf("UpdateUser failed: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", updateUserResp)

	// Delete user
	deleteUserResp, err := client.DeleteUser(context.Background(), &user.UserRequest{Id: createUserResp.Id})
	if err != nil {
		log.Fatalf("DeleteUser failed: %v", err)
	}
	fmt.Printf("Deleted user: %+v\n", deleteUserResp)
}
