package model

// User represents a user entity in the system.
type User struct {
	ID    string `json:"id" couchbase:"id"`
	Name  string `json:"name" couchbase:"name"`
	Email string `json:"email" couchbase:"email"`
}
