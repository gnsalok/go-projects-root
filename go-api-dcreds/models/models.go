// models/models.go
package models

type DynamicCredential struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	TTL  int    `json:"ttl" bson:"ttl"`
	// Add other fields as necessary
}

type CreateDynamicCredentialRequest struct {
	Name string `json:"name" binding:"required"`
	TTL  int    `json:"ttl" binding:"required,gt=0"`
	// Add other fields with validation tags
}

type UpdateDynamicCredentialRequest struct {
	Name string `json:"name" binding:"required"`
	TTL  int    `json:"ttl" binding:"required,gt=0"`
	// Add other fields with validation tags
}

type UpdateTTLRequest struct {
	TTL int `json:"ttl" binding:"required,gt=0"`
}
