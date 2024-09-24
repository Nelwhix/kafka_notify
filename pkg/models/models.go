package models

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type BaseModel interface {
	FindUserByID(ctx context.Context, id string) (User, error)
}

type Model struct {
	Conn *pgx.Conn
}

type Notification struct {
	From    User   `json:"from"`
	To      User   `json:"to"`
	Message string `json:"message"`
}

type UserNotifications map[string][]Notification
