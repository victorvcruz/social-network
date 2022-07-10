package controllers

import (
	"social_network_project/database/repository"
)

type PostsController interface {
}

type PostsControllerStruct struct {
	repository repository.PostRepository
}

func NewPostsController() PostsController {
	return &PostsControllerStruct{
		repository: repository.NewPostRepository(),
	}
}
