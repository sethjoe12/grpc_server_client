package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"test/user"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type userService struct {
	users map[int64]user.UserResponse
	user.UnimplementedUserServiceServer
	mu sync.RWMutex
}

func (s *userService) CreateUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newID := uuid.New().ID()
	newUser := &user.UserResponse{
		Id:    int64(newID),
		Name:  req.Name,
		Email: req.Email,
	}
	s.users[int64(newID)] = *newUser
	return newUser, nil
}

func (s *userService) GetUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	foundUser, ok := s.users[req.Id]
	if !ok {
		return nil, fmt.Errorf("User with ID %d not found", req.Id)
	}
	return &foundUser, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updatedUser, ok := s.users[req.Id]
	if !ok {
		return nil, fmt.Errorf("User with ID %d not found", req.Id)
	}

	if req.Name != "" {
		updatedUser.Name = req.Name
	}
	if req.Email != "" {
		updatedUser.Email = req.Email
	}

	s.users[req.Id] = updatedUser
	return &updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	deletedUser, ok := s.users[req.Id]
	if !ok {
		return nil, fmt.Errorf("User with ID %d not found", req.Id)
	}

	delete(s.users, req.Id)
	return &deletedUser, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	userService := &userService{users: make(map[int64]user.UserResponse)}
	user.RegisterUserServiceServer(s, userService)
	log.Println("Server is listening on port 50054...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
