package changelog

import "github.com/andyklimenko/testify-usage-example/api/entity"

type notificationType string

const (
	created notificationType = "CREATED"
	updated notificationType = "UPDATED"
	deleted notificationType = "DELETED"
)

type notificationBody struct {
	NotificationType notificationType `json:"notification_type"`
	User             entity.User      `json:"user"`
}
