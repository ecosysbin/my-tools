package v1

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrorInstanceIdAlreadyExists   = errors.New("instance id already exists")
	ErrorVClusterNameAlreadyExists = errors.New("vcluster name already exists")
)

func Cause(err error) string {
	cause := errors.Cause(err)
	return strings.ReplaceAll(strings.ToUpper(cause.Error()), " ", "_")
}
