// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package routers

import (
	"fmcam/routers/resources"
	"fmcam/routers/sysrouter"
	"fmcam/routers/users"
)

type Routers struct {
	Users     users.UsersRouter
	AlyRouter sysrouter.AlyRouter
	Recourse  resources.ResourcesRouter
}
