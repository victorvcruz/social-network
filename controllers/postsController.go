package controllers

import (
	"social_network_project/database/repository"
	"social_network_project/entities"
	"social_network_project/entities/response"
)

type PostsController interface {
	InsertPost(post *entities.Post) error
	FindPostsByAccountID(id *string) ([]interface{}, error)
	ChangePostDataByID(id *string, content string) error
	FindPostByID(id *string) (*response.PostResponse, error)
	ExistsPostByID(id *string) (*bool, error)
	RemovePostByID(id *string) error
}

type PostsControllerStruct struct {
	repository repository.PostRepository
}

func NewPostsController() PostsController {
	return &PostsControllerStruct{
		repository: repository.NewPostRepository(),
	}
}

func (p PostsControllerStruct) InsertPost(post *entities.Post) error {
	return p.repository.InsertPost(post)
}

func (p PostsControllerStruct) FindPostsByAccountID(id *string) ([]interface{}, error) {
	return p.repository.FindPostsByAccountID(id)
}

func (p PostsControllerStruct) ChangePostDataByID(id *string, content string) error {
	return p.repository.ChangePostDataByID(id, content)
}

func (p PostsControllerStruct) FindPostByID(id *string) (*response.PostResponse, error) {
	return p.repository.FindPostByID(id)
}

func (p PostsControllerStruct) ExistsPostByID(id *string) (*bool, error) {
	return p.repository.ExistsPostByID(id)
}

func (p PostsControllerStruct) RemovePostByID(id *string) error {
	return p.repository.RemovePostByID(id)
}
