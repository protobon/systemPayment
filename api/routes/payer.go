package routes

import (
	"database/sql"
	"log"
	"strconv"
	"systemPayment/model"

	"github.com/gin-gonic/gin"
	"gopkg.in/validator.v2"
)

type PayerRouter struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (p *PayerRouter) getPayer(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Fatal(err)
		ctx.JSON(400, map[string]string{"error": "Invalid Payer ID"})
		return
	}

	payer := model.PayerObject{ID: id}
	if err = payer.QGetPayer(db); err != nil {
		log.Fatal(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Payer not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, payer)
}

func (p *PayerRouter) createPayer(ctx *gin.Context, db *sql.DB) {
	var err error
	var payer model.PayerObject
	if err = ctx.BindJSON(&payer); err != nil {
		log.Fatal(err)
		ctx.JSON(400, map[string]string{"message": "Invalid request payload.", "error": err.Error()})
		return
	}
	if err = validator.Validate(payer); err != nil {
		ctx.JSON(400, map[string]string{"message": "Request body incomplete.", "error": err.Error()})
		return
	}

	if err = payer.QCreatePayer(db); err != nil {
		log.Fatal(err)
		ctx.JSON(400, map[string]string{"message": "Could not create new Payer", "error": err.Error()})
		return
	}

	ctx.JSON(200, payer)
}

func (p *PayerRouter) updatePayer(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Fatal(err)
		ctx.JSON(400, map[string]string{"message": "Invalid Payer id", "error": err.Error()})
		return
	}

	var payer model.PayerObject
	if err = ctx.BindJSON(&payer); err != nil {
		log.Fatal(err)
		ctx.JSON(400, map[string]string{"message": "Invalid request payload", "error": err.Error()})
		return
	}

	payer.ID = id

	if err = payer.QUpdatePayer(db); err != nil {
		log.Fatal(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"message": "Payer not found", "error": err.Error()})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, payer)
}

func (p *PayerRouter) deletePayer(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Fatal(err)
		ctx.JSON(400, map[string]string{"message": "Invalid Payer id", "error": err.Error()})
		return
	}

	payer := model.PayerObject{ID: id}
	if err = payer.QDeletePayer(db); err != nil {
		log.Fatal(err)
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"message": "Payer not found", "error": err.Error()})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, payer)
}

func (p *PayerRouter) getPayers(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	var err error
	var payers []model.PayerObject
	var payer = model.PayerObject{}
	payers, err = payer.QGetPayers(db, start, count)
	if err != nil {
		log.Fatal(err)
		ctx.JSON(500, map[string]string{"message": "Error parsing Payers", "error": err.Error()})
	}

	ctx.JSON(200, payers)
}

func (p *PayerRouter) InitializeRoutes(db *sql.DB) {
	p.Router.GET("/payers", func(ctx *gin.Context) {
		p.getPayers(ctx, db)
	})
	p.Router.POST("/payer", func(ctx *gin.Context) {
		p.createPayer(ctx, db)
	})
	p.Router.GET("/payer/:id", func(ctx *gin.Context) {
		p.getPayer(ctx, db)
	})
	p.Router.PUT("/payer/:id", func(ctx *gin.Context) {
		p.updatePayer(ctx, db)
	})
	p.Router.DELETE("/payer/:id", func(ctx *gin.Context) {
		p.deletePayer(ctx, db)
	})
}
