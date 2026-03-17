package model

import (
	"time"
)

type Gig struct {
	ID           int64 `json:"id"`
	ClientID     int64 `json:"client_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Budget       int64 `json:"budget"`
	Status       string `json:"status"`
	FreelancerID int64 `json:"freelancer_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type Proposal struct {
	ID           int64 `json:"id"`
	GigID        int64 `json:"gig_id"`
	FreelancerID int64 `json:"freelancer_id"`
	BidAmount    int64 `json:"bid_amount"`
	Status       string `json:"status"`
}

type User struct {
	ID           int64 `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	Name         string `json:"name"`
}
