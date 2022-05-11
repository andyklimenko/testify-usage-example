package changelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andyklimenko/testify-usage-example/api/entity"
)

type RestNotifier struct {
	addr string
	cli  *http.Client
}

func (n *RestNotifier) UserCreated(userID string) error {
	return n.notify(userID, created)
}

func (n *RestNotifier) UserUpdated(userID string) error {
	return n.notify(userID, updated)
}

func (n *RestNotifier) UserDeleted(userID string) error {
	return n.notify(userID, deleted)
}

func (n *RestNotifier) notify(userID string, nt notificationType) error {
	nb := notificationBody{
		NotificationType: nt,
		UserID:           userID,
	}
	bodyRaw, err := json.Marshal(nb)
	if err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}

	resp, err := n.cli.Post(n.addr, "application/json", bytes.NewReader(bodyRaw))
	if err != nil {
		return fmt.Errorf("executing request at %s: %w", n.addr, err)
	}

	defer entity.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response-code %d", resp.StatusCode)
	}

	return nil
}

func New(addr string) *RestNotifier {
	return &RestNotifier{
		addr: addr,
		cli: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}
