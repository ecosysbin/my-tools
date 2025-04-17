package app

import (
	"sync"

	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
)

var syncWorkerMap = map[string]*v1.WorkflowAppParams{}
var lock sync.Mutex

// type SyncWorkerManager struct{}

// func StartSyncWorker(app *v1.WorkflowAppParams) {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	SyncWorkerMap[app.AppRecord.Id] = app
// }

// app delete的param怎么和app creae的param区分？，都是同一个key。只能假设一个app同一时间只能在做一个操作。
func RefreshSyncWorker(app *v1.WorkflowAppParams) {
	lock.Lock()
	defer lock.Unlock()
	syncWorkerMap[app.AppRecord.Id] = app
}

// 考虑通过自清理机制，定期删除过期的SyncWorker和超额SyncWorker
func StopSyncWorker(app *v1.WorkflowAppParams) {
	lock.Lock()
	defer lock.Unlock()
	delete(syncWorkerMap, app.AppRecord.Id)
}

func GetSyncWorker(appID string) *v1.WorkflowAppParams {
	log.Infof("map: %v", syncWorkerMap)
	return syncWorkerMap[appID]
}
