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

func (n *RestNotifier) UserCreated(u entity.User) error {
	return n.notify(u, created)
}

func (n *RestNotifier) UserUpdated(u entity.User) error {
	return n.notify(u, updated)
}

func (n *RestNotifier) UserDeleted(u entity.User) error {
	return n.notify(u, deleted)
}

func (n *RestNotifier) notify(u entity.User, nt notificationType) error {
	nb := notificationBody{
		NotificationType: nt,
		User:             u,
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
