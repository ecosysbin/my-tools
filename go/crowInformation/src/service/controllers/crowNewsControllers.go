package controllers

import (
	"fmt"

	"service/model"

	"github.com/astaxie/beego"
)

type CrowNewsControllers struct {
	beego.Controller
}

func (c *CrowNewsControllers) CrowNewsUrl() {
	url := c.GetString("url")
	if len(url) == 0 {
		c.Ctx.WriteString("crownews url is empty")
		return
	}
	fmt.Println(url)
	model.ProductNewsUrl(url)
	c.Ctx.WriteString("hello crow newsurl")
}

func (c *CrowNewsControllers) CrowNews() {
	model.ConsumeNewsUrl()
	c.Ctx.WriteString("hello crow news")
}
