// Copyright (c) 2024 xueyingdu zco-inc All Right reserved. privacy_zco@outlook.com
package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"

	"fmcam/common/configs"
	"fmcam/common/inits"
	_ "fmcam/models/schema"
	_ "fmcam/systems/genclients"
	_ "fmcam/systems/genclients/runtime"
	"time"

	_ "entgo.io/ent/dialect/sql"

	gqs "github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // 注册 mysql dialect
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	WebPorts = configs.WebManagerPort
	ctx      = context.Background()
)

var (
	mode = setMode()
	ge   errgroup.Group
)

const Version = "1.025.092501"

func setMode() bool {
	gin.SetMode(configs.VersionMode) //gin.DebugMode) // TestMode
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	return true
}

var (
	logger = inits.NewLogger("main")
)

func init() {

	fmt.Printf("this is main init.\n")

	dbs := inits.InitStorage(configs.GCONFIG)

	gqs.Dialect("mysql")
	fmt.Printf("dbs init:%v \n", dbs != nil)

}

func main() {

	logger.Infof("start new project:%v", time.Now().String())

	mainser := inits.NewMainService()
	webser := inits.NewWebService()
	operser := inits.NewOperateService()
	mgs := inits.NewManager(
		mainser,
		webser,
		operser,
	)

	if err := mgs.Run(ctx); err != nil {
		logger.Fatalf("services terminated: %#v", err)
	}

}
