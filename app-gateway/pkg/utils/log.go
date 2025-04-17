package utils

import (
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
)

var Logger *log.Logger

func InitLogger() {
	Logger = &log.Logger{}
}
