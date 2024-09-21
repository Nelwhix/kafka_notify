package models

import "context"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (m *Model) FindUserByID(ctx context.Context, id string) (User, error) {
	var user User
	err := m.Conn.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
