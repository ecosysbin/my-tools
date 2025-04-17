package app

import (
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/repo"
)

func (ac *AppController) ListAppConfig() ([]repo.AppConfig, error) {
	return ac.controller.AppRepo().ListAppConfig()
}

func (ac *AppController) AddConfig(config repo.AppConfig) error {
	return ac.controller.AppRepo().AddConfig(config)
}

func (ac *AppController) DeleteConfig(appType string) error {
	return ac.controller.AppRepo().DeleteConfig(appType)
}
