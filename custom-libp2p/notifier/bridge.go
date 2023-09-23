package notifier

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"google.golang.org/protobuf/proto"
)

type FlushNotificationsBridgeOutput = models.Notifications

func FlushNotificationsBridge(proto.Message) (proto.Message, error) {
	return FlushNotifications()
}
