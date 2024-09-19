package models

import "github.com/jackc/pgx/v5"

type Model struct {
	conn *pgx.Conn
}

type Notification struct {
	From    User   `json:"from"`
	To      User   `json:"to"`
	Message string `json:"message"`
}
