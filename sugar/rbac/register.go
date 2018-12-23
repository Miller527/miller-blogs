/*
# __author__ = "Mr.chai"
# Date: 2018/12/21
*/
package rbac

import (
	"fmt"
	"miller-blogs/sugar"
)

func Register(ac *sugar.AdminConf, all bool)  {
	if ac.AccessControl != "rbac"{
		fmt.Println("Warning: AccessControl must 'rbac', register error.")
		return
	}
	ac.VerifyLoginFunc = handleVerifyLogin
	if all{
		ac.LoginFunc = handleLogin
	}
	//ac.AddGlobalMiddle(RbacLoginMiddle(), Logger() )
}