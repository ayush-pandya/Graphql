package clients

import (
	"context"
	"fmt"
	"log"

	ticketpb "github.com/ayush-pandya/Graphql/proto/ticket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TicketClient wraps the gRPC client for the ticket service
type TicketClient struct {
	conn   *grpc.ClientConn
	client ticketpb.TicketServiceClient
}

// NewTicketClient creates a new ticket service client
func NewTicketClient(address string) (*TicketClient, error) {
	// Create gRPC connection
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ticket service: %w", err)
	}

	client := ticketpb.NewTicketServiceClient(conn)

	return &TicketClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (tc *TicketClient) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}

// CreateTicket creates a new ticket via gRPC
func (tc *TicketClient) CreateTicket(ctx context.Context, title, description string, priority ticketpb.TicketPriority, assigneeID string, tags []string) (*ticketpb.Ticket, error) {
	req := &ticketpb.CreateTicketRequest{
		Title:       title,
		Description: description,
		Priority:    priority,
		AssigneeId:  assigneeID,
		Tags:        tags,
	}

	resp, err := tc.client.CreateTicket(ctx, req)
	if err != nil {
		log.Printf("Error creating ticket via gRPC: %v", err)
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return resp.Ticket, nil
}

// GetTicket retrieves a ticket by ID via gRPC
func (tc *TicketClient) GetTicket(ctx context.Context, id string) (*ticketpb.Ticket, error) {
	req := &ticketpb.GetTicketRequest{
		Id: id,
	}

	resp, err := tc.client.GetTicket(ctx, req)
	if err != nil {
		log.Printf("Error getting ticket via gRPC: %v", err)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return resp.Ticket, nil
}

// ListTickets retrieves all tickets via gRPC
func (tc *TicketClient) ListTickets(ctx context.Context, pageSize int32, pageToken string) ([]*ticketpb.Ticket, string, error) {
	req := &ticketpb.ListTicketsRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	}

	resp, err := tc.client.ListTickets(ctx, req)
	if err != nil {
		log.Printf("Error listing tickets via gRPC: %v", err)
		return nil, "", fmt.Errorf("failed to list tickets: %w", err)
	}

	return resp.Tickets, resp.NextPageToken, nil
}

// UpdateTicket updates an existing ticket via gRPC
func (tc *TicketClient) UpdateTicket(ctx context.Context, id, title, description string, status ticketpb.TicketStatus, priority ticketpb.TicketPriority, assigneeID string, tags []string) (*ticketpb.Ticket, error) {
	req := &ticketpb.UpdateTicketRequest{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		AssigneeId:  assigneeID,
		Tags:        tags,
	}

	resp, err := tc.client.UpdateTicket(ctx, req)
	if err != nil {
		log.Printf("Error updating ticket via gRPC: %v", err)
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	return resp.Ticket, nil
}

// DeleteTicket deletes a ticket via gRPC
func (tc *TicketClient) DeleteTicket(ctx context.Context, id string) (bool, error) {
	req := &ticketpb.DeleteTicketRequest{
		Id: id,
	}

	resp, err := tc.client.DeleteTicket(ctx, req)
	if err != nil {
		log.Printf("Error deleting ticket via gRPC: %v", err)
		return false, fmt.Errorf("failed to delete ticket: %w", err)
	}

	return resp.Success, nil
}
