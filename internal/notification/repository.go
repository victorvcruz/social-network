package notification

import (
	"database/sql"
	"encoding/json"
	"log"
	"social_network_project/internal/account"
	"social_network_project/internal/comment"
	"social_network_project/internal/interaction"
	"social_network_project/internal/post"
)


type NotificationRepositoryClient interface {
	HandlerNotification(notification *Notification)
	NotificationPost(id *string)
	NotificationComment(id *string)
	NotificationInteraction(id *string)
	NotificationFollowAccount(id *string)
}

type Notification struct {
	Type string
	ID   string
}

type NotificationRepository struct {
	repositoryComment comment.CommentRepository
	repositoryAccount     account.AccountRepository
	repositoryPost        post.PostRepository
	repositoryInteraction interaction.InteractionRepository
}

func NewNotificationRepository(postgresDB *sql.DB) NotificationRepositoryClient {
	return &NotificationRepository{
		repositoryAccount:     account.NewAccountRepository(postgresDB),
		repositoryPost:        post.NewPostRepository(postgresDB),
		repositoryComment:     comment.NewComentRepository(postgresDB),
		repositoryInteraction: interaction.NewInteractionRepository(postgresDB),
	}
}

func (n *NotificationRepository) HandlerNotification(notification *Notification) {

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

func (n *NotificationRepository) NotificationPost(id *string) {

	emailToNotificate, err := n.repositoryAccount.FindAccountEmailFollowersByAccountID(id)
	if err != nil {
		return
	}
	log.Println(emailToNotificate)
}

func (n *NotificationRepository) NotificationComment(id *string) {

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

func (n *NotificationRepository) NotificationInteraction(id *string) {

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

func (n *NotificationRepository) NotificationFollowAccount(id *string) {

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
