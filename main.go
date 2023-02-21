package main

import (
	"os"

	"systempayment/controller"
	"systempayment/database"
	_ "systempayment/docs"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Swagger System Payment
//	@version		1.0
//	@description	API implementation.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

//	@securitydefinitions.oauth2.application	OAuth2Application
//	@tokenUrl								https://example.com/oauth/token
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information

//	@securitydefinitions.oauth2.implicit	OAuth2Implicit
//	@authorizationUrl						https://example.com/oauth/authorize
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information

//	@securitydefinitions.oauth2.password	OAuth2Password
//	@tokenUrl								https://example.com/oauth/token
//	@scope.read								Grants read access
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information

//	@securitydefinitions.oauth2.accessCode	OAuth2AccessCode
//	@tokenUrl								https://example.com/oauth/token
//	@authorizationUrl						https://example.com/oauth/authorize
//	@scope.admin							Grants read and write access to administrative information

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	var user string
	var password string
	var dbhost string
	var dbname string
	var port string

	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbhost = os.Getenv("DATABASE_HOST")
	dbname = os.Getenv("POSTGRES_DB")
	port = os.Getenv("APPLICATION_PORT")

	r := gin.Default()
	database.DBInit(user, password, dbhost, dbname)

	c := controller.NewController()

	v1 := r.Group("/api/v1", CORSMiddleware())
	{
		payer := v1.Group("/payer")
		{
			payer.POST("/new", c.NewPayer)
			payer.GET("/payers", c.Payers)
			payer.GET(":id", c.GetPayer)
			payer.PUT("/update/:id", c.UpdatePayer)
		}
		product := v1.Group("/product")
		{
			product.POST("/new", c.NewProduct)
			product.GET("/products", c.Products)
			product.GET(":id", c.GetProduct)
			product.PUT("/update/:id", c.UpdateProduct)
		}
		order := v1.Group("/order")
		{
			order.POST("/new", c.NewOrder)
			order.GET("/orders", c.Orders)
			order.GET(":id", c.GetOrder)
		}
		dlocal := v1.Group("/dlocal")
		{
			dlocal.POST("/card", c.CreateCard)
			dlocal.POST("/secure-payment", c.MakeSecurePayment)
			dlocal.POST("/payment", c.MakePayment)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(port)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// func auth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if len(c.GetHeader("Authorization")) == 0 {
// 			httputil.NewError(c, http.StatusUnauthorized, errors.New("authorization is required header"))
// 			c.Abort()
// 		}
// 		c.Next()
// 	}
// }
