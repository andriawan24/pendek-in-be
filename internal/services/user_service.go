package services

import (
	"context"
	"database/sql"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/google/uuid"
)

type userService struct {
	ctx     context.Context
	queries *database.Queries
}

type UserService interface {
	GetUserByID(id uuid.UUID) (database.User, error)
	FindUserByEmail(email string) (database.User, error)
	FindUserByGoogleID(googleID string) (database.User, error)
	InsertUser(param database.InsertUserParams) (database.User, error)
	InsertUserWithGoogle(param database.InsertUserWithGoogleParams) (database.User, error)
	UpdateUser(param database.UpdateUserParams) (database.User, error)
}

func NewUserService(ctx context.Context, queries *database.Queries) UserService {
	return &userService{
		queries: queries,
		ctx:     ctx,
	}
}

func (s *userService) GetUserByID(id uuid.UUID) (database.User, error) {
	var user database.User

	user, err := s.queries.GetUser(s.ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) FindUserByEmail(email string) (database.User, error) {
	var user database.User

	user, err := s.queries.GetUserByEmail(s.ctx, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) InsertUser(param database.InsertUserParams) (database.User, error) {
	newUser, err := s.queries.InsertUser(s.ctx, param)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *userService) UpdateUser(param database.UpdateUserParams) (database.User, error) {
	updatedUser, err := s.queries.UpdateUser(s.ctx, param)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *userService) FindUserByGoogleID(googleID string) (database.User, error) {
	user, err := s.queries.GetUserByGoogleID(s.ctx, sql.NullString{String: googleID, Valid: true})
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) InsertUserWithGoogle(param database.InsertUserWithGoogleParams) (database.User, error) {
	newUser, err := s.queries.InsertUserWithGoogle(s.ctx, param)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}
