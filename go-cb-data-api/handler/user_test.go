package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gnsalok/go-project-root/go-db-data-api/handler"
	"github.com/gnsalok/go-project-root/go-db-data-api/model"
	"github.com/gnsalok/go-project-root/go-db-data-api/repository"
	"github.com/gnsalok/go-project-root/go-db-data-api/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetUserByID tests the GetUserByID handler.
func TestGetUserByID(t *testing.T) {
	// Initialize the mock repository
	mockRepo := new(mocks.UserRepository)

	// Sample user data
	sampleUser := &model.User{
		ID:    "user1",
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	// Setup expectations
	mockRepo.On("GetUserByID", mock.Anything, "user1").Return(sampleUser, nil)
	mockRepo.On("GetUserByID", mock.Anything, "user2").Return(nil, repository.ErrNotFound)
	mockRepo.On("GetUserByID", mock.Anything, "user3").Return(nil, errors.New("database error"))

	// Initialize handler with mock repository
	userHandler := &handler.UserHandler{Repo: mockRepo}

	// Setup Gin router with the handler
	router := gin.Default()
	router.GET("/users/:id", userHandler.GetUserByID)

	// Define test cases
	testCases := []struct {
		name         string
		userID       string
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:         "Existing User",
			userID:       "user1",
			expectedCode: http.StatusOK,
			expectedBody: sampleUser,
		},
		{
			name:         "Non-Existing User",
			userID:       "user2",
			expectedCode: http.StatusNotFound,
			expectedBody: gin.H{"error": "User not found"},
		},
		{
			name:         "Database Error",
			userID:       "user3",
			expectedCode: http.StatusInternalServerError,
			expectedBody: gin.H{"error": "Internal server error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", "/users/"+tc.userID, nil)
			assert.NoError(t, err)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the HTTP request
			router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tc.expectedCode, rr.Code)

			// Check the response body
			if tc.expectedCode == http.StatusOK {
				var user model.User
				err := json.Unmarshal(rr.Body.Bytes(), &user)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, &user)
			} else {
				var respBody map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &respBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, respBody)
			}
		})
	}

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)
}
