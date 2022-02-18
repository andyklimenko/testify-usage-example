package changelog

type notificationType string

const (
	created notificationType = "CREATED"
	updated notificationType = "CREATED"
	deleted notificationType = "CREATED"
)

type notificationBody struct {
	NotificationType notificationType `json:"notification_type"`
	UserID           string           `json:"user_id"`
}
