package controllers

import (
	"github.com/astaxie/beego"
)

type TestLayoutController struct {
	beego.Controller
}

func (t *TestLayoutController) Get() {
	t.Layout = "base.html"
	t.LayoutSections = make(map[string]string)
	//嫌base.html中内容太多的话可以拆分，然后在base中{{.layout_header}}引入即可
	t.LayoutSections["layout_header"] = "test_layout/test_layout_header.html"
	t.LayoutSections["layout_foot"] = "test_layout/test_layout_header.html"
	t.LayoutSections["layout_script"] = "test_layout/test_layout_header.html"

	t.TplName = "test_layout/test_layout.html"
}
