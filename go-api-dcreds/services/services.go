// services/services.go
package services

import (
	"errors"
	"fmt"
	"test-go/models"

	"github.com/google/uuid"
)

var (
	// In-memory data store. Replace with persistent DB in production.
	dynCredsStore = make(map[string]*models.DynamicCredential)
)

// CreateDynamicCredential creates a new dynamic credential.
func CreateDynamicCredential(req models.CreateDynamicCredentialRequest) (*models.DynamicCredential, error) {
	id := uuid.New().String()
	cred := &models.DynamicCredential{
		ID:   id,
		Name: req.Name,
		TTL:  req.TTL,
	}
	dynCredsStore[id] = cred
	return cred, nil
}

// GetDynamicCredential retrieves a dynamic credential by ID.
func GetDynamicCredential(id string) (*models.DynamicCredential, error) {
	cred, exists := dynCredsStore[id]
	if !exists {
		return nil, errors.New("dynamic credential not found")
	}
	return cred, nil
}

// UpdateDynamicCredential updates an existing dynamic credential.
func UpdateDynamicCredential(id string, req models.UpdateDynamicCredentialRequest) (*models.DynamicCredential, error) {
	cred, exists := dynCredsStore[id]
	if !exists {
		return nil, errors.New("dynamic credential not found")
	}
	cred.Name = req.Name
	cred.TTL = req.TTL
	// Update other fields as necessary
	return cred, nil
}

// DeleteDynamicCredential deletes a dynamic credential by ID.
func DeleteDynamicCredential(id string) error {
	_, exists := dynCredsStore[id]
	if !exists {
		return errors.New("dynamic credential not found")
	}
	delete(dynCredsStore, id)
	return nil
}

// UpdateTTLForAllWorkspaces updates the TTL across all Terraform workspaces.
func UpdateTTLForAllWorkspaces(id string, ttl int) error {
	// Implement actual Terraform workspace update logic here
	// For demonstration, we'll simulate with a print statement
	fmt.Printf("Updating TTL to %d for dynamic credential %s in all Terraform workspaces\n", ttl, id)
	// Return error if Terraform update fails
	return nil
}
