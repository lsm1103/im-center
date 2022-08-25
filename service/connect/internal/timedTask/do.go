package timedTask

import (
	"im-center/service/connect/internal/svc"
	"time"
)

func Run(svc *svc.ServiceContext) {
	Timer(3*time.Second, 30*time.Second, cleanConnection, nil, nil, nil)

	Timer(3*time.Second, 30*time.Second, upServiceInfo, svc, nil, nil)
}
