package api

import (
	"database/sql"
	"fmt"
	"log"
	"systemPayment/api/routes"
	"systemPayment/database"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (a *App) Run(addr string) {
	router := a.Router

	err := router.Run(addr)
	if err != nil {
		return
	}
}

func (a *App) Initialize(user string, password string, dbhost string, dbname string) {
	var err error

	fmt.Println("Initializing App...")

	a.DB = database.DBInit(user, password, dbhost, dbname)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = gin.Default()

	var dummy = routes.DummyRouter{Router: a.Router, DB: a.DB}
	var payer = routes.PayerRouter{Router: a.Router, DB: a.DB}
	var card = routes.CardRouter{Router: a.Router, DB: a.DB}

	dummy.InitializeRoutes(dummy.DB)
	payer.InitializeRoutes(payer.DB)
	card.InitializeRoutes(card.DB)

	// schedule.RunCronJobs(a.DB)

	log.Println("*************** App Running ***************")
}
