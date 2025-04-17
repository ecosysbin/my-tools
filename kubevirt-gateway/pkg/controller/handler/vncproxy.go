package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	gcpctx "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/context"
	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/response"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/framework"
)

var _ framework.ReverseProxyHandlerInterface = &VncProxyHandler{}

type VncProxyHandler struct {
	controller framework.Interface
}

func New(controller framework.Interface) *VncProxyHandler {
	return &VncProxyHandler{
		controller: controller,
	}
}

func (w *VirtualServerHandler) VncProxy(c *gcpctx.GCPContext) {
	vncServer := w.controller.ComponentConfig().GetHttpVncServer()
	target, err := url.Parse(vncServer)
	if err != nil {
		log.Infof("parse url %s err, %v", vncServer, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	uri, f := strings.CutPrefix(c.Request.RequestURI, "/vnc/v1")
	if !f {
		return
	}
	targetUrl := fmt.Sprintf("%s%s", target, uri)
	newRequest, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		log.Infof("new request err, %v", err)
		return
	}
	newRequest.Header = c.Request.Header
	proxy.ServeHTTP(c.Writer, newRequest)
}

func (w *VirtualServerHandler) PodProxy(c *gcpctx.GCPContext) {
	vncServer := w.controller.ComponentConfig().GetHttpVncServer()
	target, err := url.Parse(vncServer)
	if err != nil {
		log.Infof("parse url %s err, %v", vncServer, err)
		response.Response(c, response.ErrCreateResponse(err.Error()), nil)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	uri, f := strings.CutPrefix(c.Request.RequestURI, "/vnc/v1")
	if !f {
		return
	}
	targetUrl := fmt.Sprintf("%s%s", target, uri)
	newRequest, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		log.Infof("new request err, %v", err)
		return
	}

	newRequest.Header = c.Request.Header
	proxy.ServeHTTP(c.Writer, newRequest)
}

func (w *VirtualServerHandler) VncProxyVerify(c *gcpctx.GCPContext) {
	username := c.GetUesrName()
	// 进行授权校验
	instanceId := c.Param("instanceId")
	if _, err := w.virtualServerManager.Repo.GetVmById(username, instanceId); err != nil {
		log.Infof("can not find vm %s by user %s", instanceId, username)
		response.Response(c, response.ErrAuthorization, nil)
		return
	}
	w.VncProxy(c)
}
