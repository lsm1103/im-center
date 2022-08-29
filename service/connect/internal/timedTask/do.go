package timedTask

import (
	"im-center/common/timedTask"
	"im-center/service/connect/internal/svc"
	"time"
)

func Run(svc *svc.ServiceContext) {
	timedTask.Timer(3*time.Second, 30*time.Second, cleanConnection, nil, nil, nil)

	timedTask.Timer(3*time.Second, 30*time.Second, upServiceInfo, svc, nil, nil)
}
