package routes

import (
	"database/sql"
	"systemPayment/model"

	"github.com/gin-gonic/gin"
)

type CardRouter struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (c *CardRouter) createSecureCard(ctx *gin.Context, db *sql.DB) {
	var err error
	var secureCard model.SecureCardObject
	if err = ctx.BindJSON(&secureCard); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid request payload"})
		return
	}

	// Verificar que el Payer exista
	payer := model.PayerObject{ID: secureCard.PayerID}
	if err = payer.QGetPayerFromID(db); err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.JSON(404, map[string]string{"error": "Payer not found"})
		default:
			ctx.JSON(500, err.Error())
		}
		return
	}

	// Crear nueva tarjeta segura
	if err := secureCard.QCreateSecureCard(db); err != nil {
		ctx.JSON(400, map[string]string{"error": "Could not create new Secure Card"})
		return
	}

	ctx.JSON(200, secureCard)
}

func (c *CardRouter) InitializeRoutes(db *sql.DB) {
	c.Router.POST("/card/secure-card", func(ctx *gin.Context) {
		c.createSecureCard(ctx, db)
	})
}
