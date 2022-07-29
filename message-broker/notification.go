package message_broker

import (
	"database/sql"
	"encoding/json"
	"log"
	"social_network_project/database/repository"
)

type Notification struct {
	Type string
	ID   string
}

type NotificationController struct {
	repositoryComment     repository.CommentRepository
	repositoryAccount     repository.AccountRepository
	repositoryPost        repository.PostRepository
	repositoryInteraction repository.InteractionRepository
}

func NewNotificationController(postgresDB *sql.DB) *NotificationController {
	return &NotificationController{
		repositoryAccount:     repository.NewAccountRepository(postgresDB),
		repositoryPost:        repository.NewPostRepository(postgresDB),
		repositoryComment:     repository.NewComentRepository(postgresDB),
		repositoryInteraction: repository.NewInteractionRepository(postgresDB),
	}
}

func (n *NotificationController) HandlerNotification(notification *Notification) {

	switch notification.Type {
	case "Post":
		n.NotificationPost(&notification.ID)
		return
	case "Comment":
		n.NotificationComment(&notification.ID)
	case "Interaction":
		n.NotificationInteraction(&notification.ID)
	case "FollowAccount":
		n.NotificationFollowAccount(&notification.ID)
	}
}

func (n *NotificationController) NotificationPost(id *string) {

	emailToNotificate, err := n.repositoryAccount.FindAccountEmailFollowersByAccountID(id)
	if err != nil {
		return
	}
	log.Println(emailToNotificate)
}

func (n *NotificationController) NotificationComment(id *string) {

	comment, err := n.repositoryComment.FindCommentByID(id)
	if err != nil {
		return
	}

	if !comment.CommentID.Valid {
		emailToNotificate, err := n.repositoryComment.FindAccountEmailOfPostByCommentID(id)
		if err != nil {
			return
		}
		log.Println(emailToNotificate)
		return
	}

	emailToNotificate, err := n.repositoryComment.FindAccountEmailOfPostAndCommentByCommentID(id)
	if err != nil {
		return
	}
	log.Println(emailToNotificate)
}

func (n *NotificationController) NotificationInteraction(id *string) {

	interaction, err := n.repositoryInteraction.FindInteractionByID(id)
	if err != nil {
		return
	}

	if interaction.PostID.Valid {
		emailToNotificate, err := n.repositoryInteraction.FindAccountEmailOfPostByInteractionID(id)
		if err != nil {
			return
		}
		log.Println(emailToNotificate)
		return
	}

	if interaction.CommentID.Valid {
		emailToNotificate, err := n.repositoryInteraction.FindAccountEmailOfCommentByInteractionID(id)
		if err != nil {
			return
		}
		log.Println(emailToNotificate)
		return
	}
}

func (n *NotificationController) NotificationFollowAccount(id *string) {

	emailToNotificate, err := n.repositoryAccount.FindAccountEmailByID(id)
	if err != nil {
		return
	}
	log.Println(emailToNotificate)
	return
}

func CreateNotificationJson(typeN, id string) string {
	not := &Notification{
		Type: typeN,
		ID:   id,
	}
	jsonStr, _ := json.Marshal(not)

	return string(jsonStr)
}
