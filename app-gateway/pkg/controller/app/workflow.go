package app

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/utils"
)

type AppWorkFlow struct {
	Metadata *v1.WorkflowAppParams
	Works    []AppWork
}

type AppWork struct {
	WorkName     string
	Work         WorkFlowFn
	RollBack     WorkFlowFn
	FailCallBack WorkFlowFn
}

type WorkFlowFn func(*v1.WorkflowAppParams) error

const (
	WORK_CREATE_CREATEINSTANCE       = "create app instance"
	WORK_CREATE_UPDATECREATESTATUS   = "update creating status"
	WORK_CREATE_CHECKAPPCREATESTATUS = "check app create status"
	WORK_CREATE_UPDATERUNNIGSTATUS   = "update running status"

	WORK_UPDATE_UPDATEINSTANCE       = "update app instance"
	WORK_UPDATE_UPDATEUPDATINGSTATUS = "update updating status"
	WORK_UPDATE_CHECKAPPUPDATESTATUS = "check app update status"
	WORK_UPDATE_UPDATERUNNIGSTATUS   = "update running status"

	WORK_DELETE_DELETEINSTANCE       = "delete app instance"
	WORK_DELETE_CHECKAPPDELETESTATUS = "check app delete status"
	WORK_DELETE_UPDATEDELETEDTATUS   = "update deleted status"

	WORK_PAUSE_PAUSEINSTANCE       = "pause app instance"
	WORK_PAUSE_CHECKAPPPAUSESTATUS = "check app pause status"
	WORK_PAUSE_UPDATEDPAUSEDTATUS  = "update paused status"

	WORK_RESUME_RESUMEINSTANCE        = "resume app instance"
	WORK_RESUME_CHECKAPPPSESUMESTATUS = "check app resume status"
	WORK_RESUME_UPDATERESUMESTATUS    = "update resume status"
)

func (workFlow *AppWorkFlow) Start() {
	rollBackWorks := []WorkFlowFn{}
	for _, work := range workFlow.Works {
		if work.RollBack != nil {
			rollBackWorks = append(rollBackWorks, work.RollBack)
		}
		metadata := workFlow.Metadata
		if err := work.Work(metadata); err != nil {
			metadata.Logger.Infof("app %s process %s do err,%+v", workFlow.Metadata.AppRecord.Name, work.WorkName, err)
			// 更新失败状态
			if err := work.FailCallBack(metadata); err != nil {
				metadata.Logger.Warnf("app %s process %s FailCallBack err,%v", workFlow.Metadata.AppRecord.Name, work.WorkName, err)
				return
			}
			for _, rollBackWork := range rollBackWorks {
				if err := rollBackWork(workFlow.Metadata); err != nil {
					metadata.Logger.Warnf("rollback %s process %s do err,%v", workFlow.Metadata.AppRecord.Name, work.WorkName, err)
				}
			}
			// 回滚完成即任务退出
			return
		}
	}
}

func WorkforTimeout(workName string, period int64, timeOut int64, work WorkFlowFn) WorkFlowFn {
	return func(app *v1.WorkflowAppParams) error {
		for timeOut > 0 {
			if err := work(app); err != nil {
				if !errors.Is(err, ErrNeedRetry) {
					log.Errorf("%s err, %v, break work", workName, err)
					return err
				}
				// retry时，状态不变，event每次入库对db压力较大，不入库无法直接同步event到客户端，权衡利弊考虑使用缓存。
				// 后续考虑多协程消息同步效率，考虑使用管道通信。
				RefreshSyncWorker(app)
				log.Infof("%s err, %v, continue work", workName, err)
				time.Sleep(time.Duration(period) * time.Second)
				timeOut--
				continue
			}
			return nil
		}
		app.Conditions.Status = v1.ActionStatusFailed(app.Action)
		app.Conditions.Reason = fmt.Sprintf("%s until timeout: [%s]", workName, app.Conditions.Status)
		app.Conditions.Events = append(app.Conditions.Events, utils.ParseTimeEvent(app.Conditions.Reason))
		// 失败时将app的状态恢复到扭转之前的状态, 之前的状态为空时，则直接使用二级状态作为当前状态
		if app.Conditions.PreStatus != "" {
			app.AppRecord.Status = app.Conditions.PreStatus
		} else {
			app.AppRecord.Status = app.Conditions.Status
		}

		app.AppRecord.Reason = app.Conditions.Reason
		return errors.Errorf("%s until timeout", workName)
	}
}

// var ErrNeedRetry = errors.New("need retry")
// err信息会在等待阶段作为event展示
var ErrNeedRetry = errors.New("wait for success")
