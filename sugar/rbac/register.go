/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import (
	"errors"
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
	ac.LoginFunc =handlerLogin
	ac.AddGlobalMiddle(RbacLoginMiddle(), BehaviorLog())

}

var ParamsRbac Params

type Params struct {
	whiteList  []string
	blackList  []string
	loginUrl   string
	indexUrl   string
	staticPath string
	urlPrefix  string
}

func (par *Params) WhiteList(urls ...string) {
	par.whiteList = append(par.whiteList, urls...)
}

func (par *Params) BlackList(urls ...string) {
	par.blackList = append(par.blackList, urls...)

}

func (par *Params) Url(login, index string) {
	par.loginUrl = login
	par.indexUrl = index

}

func (par *Params) Path(p, s string) {
	par.urlPrefix = p
	par.staticPath = s
}

func ParamsNew(whiteList, blackList []string, loginUrl,indexUrl, staticPath, urlPrefix string) {
	ParamsRbac = Params{
		whiteList:  whiteList,
		blackList:  blackList,
		loginUrl:   loginUrl,
		indexUrl:   indexUrl,
		staticPath: staticPath,
		urlPrefix:  urlPrefix,
	}
}
