package routes

import (
	"database/sql"
	"log"
	"strconv"
	"systemPayment/model"

	"github.com/gin-gonic/gin"
)

type DummyRouter struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (d *DummyRouter) getDummy(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Dummy ID"})
		return
	}

	dummy := model.DummyObject{ID: id}
	if err = dummy.QGetDummy(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Dummy not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, dummy)
}

func (d *DummyRouter) createDummy(ctx *gin.Context, db *sql.DB) {
	var dummy model.DummyObject
	if err := ctx.BindJSON(&dummy); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := dummy.QCreateDummy(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new Dummy"})
		return
	}

	ctx.JSON(200, dummy)
}

func (d *DummyRouter) updateDummy(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Dummy id"})
		return
	}

	var dummy model.DummyObject
	if err = ctx.BindJSON(&dummy); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	dummy.ID = id

	if err = dummy.QUpdateDummy(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Dummy not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, dummy)
}

func (d *DummyRouter) deleteDummy(ctx *gin.Context, db *sql.DB) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid Dummy id"})
		return
	}

	dummy := model.DummyObject{ID: id}
	if err = dummy.QDeleteDummy(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Dummy not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	ctx.JSON(200, dummy)
}

func (d *DummyRouter) getDummies(ctx *gin.Context, db *sql.DB) {
	start, _ := strconv.Atoi(ctx.Query("start"))
	count, _ := strconv.Atoi(ctx.Query("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	var dummy = model.DummyObject{}
	dummies, err := dummy.QGetDummies(db, start, count)
	if err != nil {
		log.Println(err)
	}

	ctx.JSON(200, dummies)
}

func (d *DummyRouter) InitializeRoutes(db *sql.DB) {
	d.Router.GET("/dummies", func(ctx *gin.Context) {
		d.getDummies(ctx, db)
	})
	d.Router.POST("/dummy", func(ctx *gin.Context) {
		d.createDummy(ctx, db)
	})
	d.Router.GET("/dummy/:id", func(ctx *gin.Context) {
		d.getDummy(ctx, db)
	})
	d.Router.PUT("/dummy/:id", func(ctx *gin.Context) {
		d.updateDummy(ctx, db)
	})
	d.Router.DELETE("/dummy/:id", func(ctx *gin.Context) {
		d.deleteDummy(ctx, db)
	})
}
