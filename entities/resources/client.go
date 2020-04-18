package resources

import "github.com/google/uuid"

type Client struct {
	ID     uuid.UUID
	Secret uuid.UUID
	Domain string
}

func NewClient(domain string) *Client {
	id := uuid.New()
	secret := uuid.New()
	return &Client{ID: id, Secret: secret, Domain: domain}
}
