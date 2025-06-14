package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	ticketpb "github.com/ayush-pandya/Graphql/proto/ticket"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ticketServer implements the TicketService gRPC service
type ticketServer struct {
	ticketpb.UnimplementedTicketServiceServer
	tickets map[string]*ticketpb.Ticket
	mu      sync.RWMutex
	counter int64
}

// newTicketServer creates a new ticket server with some sample data
func newTicketServer() *ticketServer {
	server := &ticketServer{
		tickets: make(map[string]*ticketpb.Ticket),
		counter: 0,
	}

	// Add some sample tickets
	now := timestamppb.Now()
	sampleTickets := []*ticketpb.Ticket{
		{
			Id:          "ticket-1",
			Title:       "Fix login bug",
			Description: "Users cannot log in with their email",
			Status:      ticketpb.TicketStatus_TICKET_STATUS_OPEN,
			Priority:    ticketpb.TicketPriority_TICKET_PRIORITY_HIGH,
			AssigneeId:  "user-123",
			Tags:        []string{"bug", "urgent"},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Id:          "ticket-2",
			Title:       "Add dark mode",
			Description: "Implement dark mode for better user experience",
			Status:      ticketpb.TicketStatus_TICKET_STATUS_IN_PROGRESS,
			Priority:    ticketpb.TicketPriority_TICKET_PRIORITY_MEDIUM,
			AssigneeId:  "user-456",
			Tags:        []string{"feature", "ui"},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for _, ticket := range sampleTickets {
		server.tickets[ticket.Id] = ticket
	}

	return server
}

// CreateTicket creates a new ticket
func (s *ticketServer) CreateTicket(ctx context.Context, req *ticketpb.CreateTicketRequest) (*ticketpb.CreateTicketResponse, error) {
	log.Printf("gRPC Microservice: Creating ticket - Title: %s", req.Title)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	ticketID := fmt.Sprintf("ticket-%d", s.counter+100) // Start from 100 to avoid conflicts

	now := timestamppb.Now()
	ticket := &ticketpb.Ticket{
		Id:          ticketID,
		Title:       req.Title,
		Description: req.Description,
		Status:      ticketpb.TicketStatus_TICKET_STATUS_OPEN,
		Priority:    req.Priority,
		AssigneeId:  req.AssigneeId,
		Tags:        req.Tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.tickets[ticketID] = ticket
	log.Printf("gRPC Microservice: Ticket created successfully - ID: %s", ticketID)

	return &ticketpb.CreateTicketResponse{
		Ticket: ticket,
	}, nil
}

// GetTicket retrieves a ticket by ID
func (s *ticketServer) GetTicket(ctx context.Context, req *ticketpb.GetTicketRequest) (*ticketpb.GetTicketResponse, error) {
	log.Printf("gRPC Microservice: Getting ticket - ID: %s", req.Id)

	s.mu.RLock()
	defer s.mu.RUnlock()

	ticket, exists := s.tickets[req.Id]
	if !exists {
		return nil, fmt.Errorf("ticket not found: %s", req.Id)
	}

	log.Printf("gRPC Microservice: Ticket retrieved successfully - ID: %s", req.Id)
	return &ticketpb.GetTicketResponse{
		Ticket: ticket,
	}, nil
}

// ListTickets retrieves all tickets
func (s *ticketServer) ListTickets(ctx context.Context, req *ticketpb.ListTicketsRequest) (*ticketpb.ListTicketsResponse, error) {
	log.Println("gRPC Microservice: Listing tickets")

	s.mu.RLock()
	defer s.mu.RUnlock()

	tickets := make([]*ticketpb.Ticket, 0, len(s.tickets))
	for _, ticket := range s.tickets {
		tickets = append(tickets, ticket)
	}

	log.Printf("gRPC Microservice: Listed %d tickets", len(tickets))
	return &ticketpb.ListTicketsResponse{
		Tickets: tickets,
	}, nil
}

// UpdateTicket updates an existing ticket
func (s *ticketServer) UpdateTicket(ctx context.Context, req *ticketpb.UpdateTicketRequest) (*ticketpb.UpdateTicketResponse, error) {
	log.Printf("gRPC Microservice: Updating ticket - ID: %s", req.Id)

	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, exists := s.tickets[req.Id]
	if !exists {
		return nil, fmt.Errorf("ticket not found: %s", req.Id)
	}

	// Update fields
	if req.Title != "" {
		ticket.Title = req.Title
	}
	if req.Description != "" {
		ticket.Description = req.Description
	}
	if req.Status != ticketpb.TicketStatus_TICKET_STATUS_UNSPECIFIED {
		ticket.Status = req.Status
	}
	if req.Priority != ticketpb.TicketPriority_TICKET_PRIORITY_UNSPECIFIED {
		ticket.Priority = req.Priority
	}
	if req.AssigneeId != "" {
		ticket.AssigneeId = req.AssigneeId
	}
	if len(req.Tags) > 0 {
		ticket.Tags = req.Tags
	}

	ticket.UpdatedAt = timestamppb.Now()
	s.tickets[req.Id] = ticket

	log.Printf("gRPC Microservice: Ticket updated successfully - ID: %s", req.Id)
	return &ticketpb.UpdateTicketResponse{
		Ticket: ticket,
	}, nil
}

// DeleteTicket deletes a ticket
func (s *ticketServer) DeleteTicket(ctx context.Context, req *ticketpb.DeleteTicketRequest) (*ticketpb.DeleteTicketResponse, error) {
	log.Printf("gRPC Microservice: Deleting ticket - ID: %s", req.Id)

	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tickets[req.Id]
	if !exists {
		return &ticketpb.DeleteTicketResponse{Success: false}, nil
	}

	delete(s.tickets, req.Id)
	log.Printf("gRPC Microservice: Ticket deleted successfully - ID: %s", req.Id)

	return &ticketpb.DeleteTicketResponse{Success: true}, nil
}

func main() {
	log.Println("🚀 Starting Ticket gRPC Microservice...")

	// Create TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()

	// Register service
	ticketService := newTicketServer()
	ticketpb.RegisterTicketServiceServer(s, ticketService)

	log.Println("✅ Ticket Service registered")

	// Start server in goroutine
	go func() {
		log.Println("🌐 Ticket gRPC Microservice listening on :50051")
		log.Println("🎫 Ready to handle ticket operations")

		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Ticket gRPC Microservice...")
	s.GracefulStop()
	log.Println("👋 Ticket gRPC Microservice stopped")
}
