package service

import "errors"

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID    int64
	Email string
	Role  string
}

type AuthService struct {
	users map[int64]User
}

func NewAuthService() *AuthService {
	return &AuthService{
		users: map[int64]User{
			42: {
				ID:    42,
				Email: "user42@example.com",
				Role:  "user",
			},
			99: {
				ID:    99,
				Email: "admin99@example.com",
				Role:  "admin",
			},
		},
	}
}

func (s *AuthService) ValidateUser(userID int64) (*User, error) {
	if userID <= 0 {
		return nil, ErrUserNotFound
	}

	user, ok := s.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}

	return &user, nil
}
