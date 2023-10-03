package notifier

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"log"
)

func FlushNotifications() (*models.Notifications, error) {
	if len(notificationChannel) == 0 {
		latestNotification, open := <-notificationChannel
		if !open {
			return nil, ErrNotifierClosed
		}
		return &models.Notifications{Notification: []*models.Notification{latestNotification}}, nil
	}

	numNotifications := len(notificationChannel)
	queuedNotifications := make([]*models.Notification, numNotifications)
	for i := 0; i < numNotifications; i++ {
		select {
		case notification := <-notificationChannel:
			queuedNotifications[i] = notification
		}
	}

	log.Printf("flushed notifications (BuffLenAfterFlush: %d) (Total Flused: %d)",
		len(notificationChannel), numNotifications)
	return &models.Notifications{Notification: queuedNotifications}, nil
}
