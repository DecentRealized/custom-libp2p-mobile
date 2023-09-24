package notifier

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/config"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
)

var notificationChannel = make(chan *models.Notification, config.MaxNotifierQueueSize)

func Reset() error {
	QueueInfo("Notification Stream Reset")
	close(notificationChannel)
	notificationChannel = make(chan *models.Notification, config.MaxNotifierQueueSize)
	return nil
}
