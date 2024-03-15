package adapters

import (
	"context"
	"github.com/WildEgor/g-core/pkg/libs/notifier"
	"github.com/WildEgor/gAuth/internal/configs"
)

type Notifier struct {
	client *notifier.NotifierClient
}

func NewNotifierAdapter(config *configs.NotifierConfig) *Notifier {

	client, err := notifier.NewNotifierClient(&notifier.NotifierConfig{
		DSN:      config.DSN,
		Exchange: config.Exchange,
	})
	if err != nil {
		panic(err) // TODO: handle error
	}

	return &Notifier{
		client,
	}
}

func (n *Notifier) Notify(payload *notifier.NotificationPayload) error {
	return n.client.Notify(context.TODO(), payload)
}

func (n *Notifier) Close() error {
	return n.client.Close()
}
