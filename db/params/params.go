package p

import "encore.app/db"

type CreateUserParams struct {
	ID             string      `json:"id"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	Email          db.Text `json:"email"`
	ProviderID     string      `json:"provider_id"`
	ZillowUsername db.Text `json:"zillow_username"`
}

type CreateUserWithZillowUsernameParams struct {
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	ZillowUsername db.Text `json:"zillow_username"`
	ProviderID     string      `json:"provider_id"`
}

type UpdateUserZillowUsernameParams struct {
	ZillowUsername db.Text `json:"zillow_username"`
	ID             string      `json:"id"`
}

