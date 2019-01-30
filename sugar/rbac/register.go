/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import (
	"errors"
	"fmt"
	"log"
	"miller-blogs/sugar"
)

func Register(ac *sugar.AdminConf) {
	if ParamsRbac.loginUrl == "" {
		panic(errors.New("RbacRegisterError: login url not found."))
	}

	if ac.AccessControl != "rbac" {
		log.Println("Warning: AccessControl must 'rbac', register error.")
		return
	}
	ac.VerifyLoginFunc = handlerVerifyLogin
	ac.LoginFunc = localLogin
	ac.AddGroupMiddle(RbacLoginMiddle(), groupParamsMiddle(), BehaviorLog())

}

var ParamsRbac Params

type Params struct {
	config     *sugar.AdminConf
	whiteList  []string
	blackList  []string
	loginUrl   string
	indexUrl   string
	staticPath string
	urlPrefix  string
	extendKey  string
	relativeKey  string

}

func (par *Params) WhiteList(urls ...string) {
	par.whiteList = append(par.whiteList, urls...)
}

func (par *Params) BlackList(urls ...string) {
	par.blackList = append(par.blackList, urls...)

}

func (par *Params) Url(relative, extend string) {
	par.extendKey = extend
	par.relativeKey = relative

}

func (par *Params) UrlKey(login, index string) {
	par.loginUrl = login
	par.indexUrl = index

}

func (par *Params) Path(p, s string) {
	par.urlPrefix = p
	par.staticPath = s
}
func (par *Params) SetAdmin(conf *sugar.AdminConf) {
	par.config = conf
	par.setParams()
}

func (par *Params) setParams() {
	conf := par.config
	par.urlPrefix = conf.Prefix
	par.staticPath = conf.Static
	fmt.Println(par.staticPath)
	par.loginUrl = conf.Prefix + "login"
	par.indexUrl = conf.Prefix + "index"
	par.whiteList = conf.WhiteUrls()
	par.blackList = conf.BlackUrls()
	par.extendKey = conf.ExtendKey
	par.relativeKey = conf.RelativeKey

}

func ParamsNew(conf *sugar.AdminConf, whiteList, blackList []string, loginUrl, indexUrl, staticPath, urlPrefix string) {
	ParamsRbac = Params{
		config:     conf,
		whiteList:  whiteList,
		blackList:  blackList,
		loginUrl:   loginUrl,
		indexUrl:   indexUrl,
		staticPath: staticPath,
		urlPrefix:  urlPrefix,
	}
}

var manuGenHandle sugar.MenuGenerator

func SetManuGenerator(f sugar.MenuGenerator) {
	if f != nil {
		manuGenHandle = f
	}

}

func secondaryMenu(sortManu sugar.SortedMenu) string {
	htmlMenu := `<ul class="nav nav-pills nav-stacked main-menu">`
	for _, menu := range sortManu {

		if menu.Children == nil {

			htmlMenu = htmlMenu + `<li class="ajax-link"><a href="`+ menu.Url +`"><i class="glyphicon `+menu.Icon+
				`"></i><span> ` + menu.Title + `</span></a>`
		}else{
			htmlMenu = htmlMenu + `<li class="accordion"><a href="`+ menu.Url +`"><i class="glyphicon `+menu.Icon+
				`"></i><span> ` + menu.Title + `</span></a>` + `<ul class="nav nav-pills nav-stacked">`

			for _, s := range menu.Children {
				htmlMenu = htmlMenu + `<li><a href="`+s.Url+`"><i class="glyphicon `+s.Icon+
					`"></i><span> ` + s.Title + `</span></a></li>`

			}
			htmlMenu += `</ul>`
		}
		htmlMenu += `</li>`
	}
	htmlMenu += `</ul>`
	return htmlMenu
}

func init() {
	manuGenHandle = secondaryMenu
}
