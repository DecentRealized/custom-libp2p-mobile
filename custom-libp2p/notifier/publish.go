package notifier

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"log"
)

func QueueMessage(message *models.Message) {
	log.Printf("QueueMessage: %v", message)
	notification := &models.Notification{
		Data: &models.Notification_MessageNotification{MessageNotification: message},
	}
	QueueNotification(notification)
}

func QueueWarning(warning *models.Warning) {
	log.Printf("QueueWarning: %v", warning)
	notification := &models.Notification{
		Data: &models.Notification_WarningNotification{WarningNotification: warning},
	}
	QueueNotification(notification)
}

func QueueInfo(info string) {
	log.Printf("QueueInfo: %v", info)
	notification := &models.Notification{
		Data: &models.Notification_InfoNotification{InfoNotification: info},
	}
	QueueNotification(notification)
}

func QueueNotification(notification *models.Notification) {
	notificationChannel <- notification
	log.Printf("QueueLen: %v", len(notificationChannel))
}
