package r

import "encore.app/db"

type GetAllUsersBasicRow struct {
	ID             string      `json:"id"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	ZillowUsername string `json:"zillow_username"`
}

