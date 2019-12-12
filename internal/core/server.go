package core

import (
	"errors"
	"net"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/OrkaConsultants/vodacom-balance-api/docs"
	"github.com/OrkaConsultants/vodacom-balance-api/internal/models"
	"github.com/OrkaConsultants/vodacom-balance-api/internal/utils"
)

// @contact.name Orka Consultants
// @contact.url https://Orka.xyz/
// @contact.email johan@Orka.xyz

// @license.name Apache 2.0

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func SetupServer() {

	// programatically set swagger info
	docs.SwaggerInfo.Title = "Vodacom API"
	docs.SwaggerInfo.Description = "Get some balances from the Vodacom"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/" + viper.GetString("app.name")

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	// r.Use(gin.Logger())

	// swag init -g server.go -d ./internal/core/
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(jwtProcessingMiddleware())

	// Setup swagger route

	var listener net.Listener
	var err error

	// Check if port has been set in the config file
	if viper.IsSet("app.port") {
		listener, err = net.Listen("tcp", ":"+viper.GetString("app.port"))
	} else {
		// Otherwise attach to any open port
		listener, err = net.Listen("tcp", ":0")
		_, port, _ := net.SplitHostPort(listener.Addr().String())

		viper.Set("eureka.app-ip", GetOutboundIP())
		viper.Set("app.port", port)
	}

	log.Info("Service IP:", GetOutboundIP(), ":", listener.Addr().(*net.TCPAddr).Port)

	if err != nil {
		panic(err)
	}

	// CORS
	r.Use(CORSMiddleware())

	// Setup all other routes
	go SetupAPIRoutes(listener, r)

	return
}

// GetOutboundIP - Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func jwtProcessingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		JwtToken := ctx.GetHeader(viper.GetString("jwt-config.HEADER_STRING"))

		// If authentication is enabled and there is no token, abort
		if viper.GetBool("app.authentication") && JwtToken == "" {
			log.Error("No authorization header set")
			utils.NewError(ctx, http.StatusUnauthorized, errors.New("No authorization header set"))
			ctx.Abort()
			return
		}

		// Initialize a new instance of `Claims`
		claims := &models.JWTClaimsStruct{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(JwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		_ = tkn
		_ = err

		ctx.Set("username", claims.Username)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
