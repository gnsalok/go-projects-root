package repository

import (
	"context"
	"errors"

	"github.com/couchbase/gocb/v2"
	"github.com/gnsalok/go-project-root/go-db-data-api/model"
)

var (
	// ErrNotFound is returned when a user is not found in the database.
	ErrNotFound = errors.New("user not found")
)

// UserRepository defines the methods that any
// data storage provider must implement to get User information.
type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

// userRepository implements UserRepository interface.
type userRepository struct {
	bucket *gocb.Bucket
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(bucket *gocb.Bucket) UserRepository {
	return &userRepository{bucket: bucket}
}

// GetUserByID retrieves a user by their ID from Couchbase.
func (r *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	collection := r.bucket.DefaultCollection()
	getResult, err := collection.Get(id, nil)
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	err = getResult.Content(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
