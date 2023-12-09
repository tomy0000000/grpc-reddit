package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedditAPI struct {
	mock.Mock
}

func (m *MockRedditAPI) GetPost(id int32) (*RedditPost, error) {
	args := m.Called(id)
	return args.Get(0).(*RedditPost), args.Error(1)
}

func (m *MockRedditAPI) GetTopComments(postId int32, limit int32) ([]*RedditComment, error) {
	args := m.Called(postId, limit)
	return args.Get(0).([]*RedditComment), args.Error(1)
}

func (m *MockRedditAPI) ExpandCommentBranch(commentId int32, limit int32) ([]*RedditComment, error) {
	args := m.Called(commentId, limit)
	return args.Get(0).([]*RedditComment), args.Error(1)
}

func TestDemoFunc(t *testing.T) {
	// Initialize the mock
	mockAPI := new(MockRedditAPI)

	// Setup expectations
	mockAPI.On("GetPost", int32(1)).Return(&RedditPost{Id: 1}, nil)
	mockAPI.On("GetTopComments", int32(1), int32(10)).Return([]*RedditComment{{Id: 2}}, nil)
	mockAPI.On("ExpandCommentBranch", int32(2), int32(10)).Return([]*RedditComment{{Content: "Test Content"}}, nil)

	// Call the function
	result, err := demoFunc(mockAPI)
	assert.NoError(t, err)
	assert.Equal(t, "Test Content", result)

	// Assert that all expectations were met
	mockAPI.AssertExpectations(t)
}
