package controllers

import (
	"social_network_project/database/repository"
	"social_network_project/entities"
)

type CommentsController interface {
	InsertComment(comment *entities.Comment) error
}

type CommentsControllerStruct struct {
	repository repository.CommentRepository
}

func NewCommentsController() CommentsController {
	return &CommentsControllerStruct{
		repository: repository.NewComentRepository(),
	}
}

func (c CommentsControllerStruct) InsertComment(comment *entities.Comment) error {
	return c.repository.InsertComment(comment)
}
