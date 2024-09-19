package models

import "context"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (m *Model) findUserByID(ctx context.Context, id string) (User, error) {
	var user User
	err := m.conn.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
