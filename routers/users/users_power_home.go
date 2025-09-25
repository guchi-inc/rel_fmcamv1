// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package users

import (
	"fmcam/ctrlib/middleware"

	"github.com/gin-gonic/gin"
)

type UsersRouter struct{}

// 用户资源查询
func (s *UsersRouter) InitUsers(r *gin.RouterGroup) {

	rstpCtl := ApiUserCtl.UserApis
	r.GET("/ok", rstpCtl.Pok) //当前服务测试信息
	// 路由
	r.POST("/signin", rstpCtl.SignInUser) //登录
	r.POST("/register", rstpCtl.Register) //注册用户
	r.Use(middleware.AuthorizeJWT())
	r.POST("/account/register", rstpCtl.RegisterAccount) //客户端新增 账户用户

	apiKeyG := r.Group("/apikey")
	apiKeyG.Use(middleware.AuthorizeJWT())
	apiKeyG.POST("/update", rstpCtl.UpdateApiKeys)
	apiKeyG.GET("/list", rstpCtl.GetApiKeys)
	apiKeyG.POST("/new", rstpCtl.NewApiKeys)

	//旧表 fm_api_key
	// apiKeyG.GET("/lists", rstpCtl.ListAPIKeysHandler)
	// apiKeyG.POST("/updates", rstpCtl.UpdateAPIKeyHandler)
	// apiKeyG.POST("/news", rstpCtl.GenerateAPIKeyHandler)

	menus := r.Group("/menu")
	menus.Use(middleware.AuthorizeJWT())
	menus.GET("/list", rstpCtl.MenuLists)      //后台管理 用户的菜单列表,需要为管理员
	menus.POST("/update", rstpCtl.MenusUpdate) //admin 编辑菜单权限 todo

	menus.GET("/customer/list", rstpCtl.GetMenuList) //按条件查功能

	////角色管理
	roles := r.Group("/role")
	roles.Use(middleware.AuthorizeJWT())

	roles.GET("/list", rstpCtl.RolesList)     //查询 后台角色单列表
	roles.POST("/update", rstpCtl.RoleUpdate) //更新 后台角色
	roles.POST("/new", rstpCtl.RoleInsert)    //新增 后台角色, fm_role, fm_role_limits

	roles.GET("/user/list", rstpCtl.RoleUserList)         //用户，角色 信息查询
	roles.POST("/user/update", rstpCtl.RoleUserOneUpdate) //用户，角色 信息更新
	roles.POST("/user/new", rstpCtl.RoleUserInsert)       //新增 用户，角色关系

	//token刷新和失效
	lt := r.Group("")

	lt.POST("/logout", rstpCtl.LogOut) //登出
	lt.Use(middleware.AuthorizeJWT())
	lt.POST("/token", rstpCtl.TokenUser) //更新用户token

	//限制访问路由，只有经过授权的用户组，可以访问该路由
	// service.GETByLimited(r, "/users/lists", rstpCtl.GetUsers) //有权限审核
	// service.GETByLimited(r, "/user/info", rstpCtl.GetUserInfo)

	//账号下用户查询列表
	userDep := r.Group("/user")
	userDep.POST("/new", rstpCtl.Register) //添加用户账号
	//登陆时查账户
	userDep.GET("/phone/account", rstpCtl.GetAccountsByPhone)

	userDep.GET("/ping", rstpCtl.Ping)            //当前用户信息
	userDep.GET("/dep/list", rstpCtl.GetDepments) //部门列表参数dep_name 为空，则返回全部，否则返回匹配的
	//用户密码找回 邮箱
	userDep.GET("/reback/code", rstpCtl.EmailCodeSend)             //发送邮件验证码
	userDep.POST("/reback/password", rstpCtl.ResetPasswdWithEmail) //通过邮件验证码重置密码

	userDep.Use(middleware.AuthorizeJWT())
	userDep.POST("/dep/new", rstpCtl.NewDepments)       //部门列表参数dep_name 为空，则返回全部，否则返回匹配的
	userDep.POST("/dep/update", rstpCtl.UpdateDepments) //部门列表参数dep_name 为空，则返回全部，否则返回匹配的

	//用户管理
	userDep.GET("/list", rstpCtl.GetUsersList)              //查询  租户信息相关的用户账号
	userDep.POST("/update", rstpCtl.UpdateUser)             //更新用户
	userDep.POST("/passswd/check", rstpCtl.BussTokenUser)   //检查用户旧密码 并刷新token
	userDep.POST("/passswd/update", rstpCtl.UpdatePassword) //更新自己的密码

	//管理员重置密码
	userDep.POST("/passswd/restore", rstpCtl.RestorePasswd)

	//全账号管理
	userDep.GET("/account/list", rstpCtl.GetAccountsList)      //查询账号列表
	userDep.GET("/account/info", rstpCtl.GetAccountInfo)       //查询账号和租户信息
	userDep.POST("/account/update", rstpCtl.UpdateUserAccount) //更新账号信息

	//账户修改日志
	logs := r.Group("/logs")
	logs.GET("/users/list", rstpCtl.GetUserSQLLogs)     //日志记录
	logs.GET("/devices/list", rstpCtl.GetDeviceSQLLogs) //日志记录
}
