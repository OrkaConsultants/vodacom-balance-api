package core

import (
	"net"
	"net/http"

	"github.com/OrkaConsultants/vodacom-balance-api/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(listener net.Listener, r *gin.Engine) {

	v1 := r.Group("/api/")
	{

		hi := v1.Group("/hello")
		{
			hc := new(controllers.HelloController)
			hi.GET("", hc.SayHello)
		}
		balance := v1.Group("/balance")
		{
			bc := new(controllers.BalanceController)
			balance.GET("", bc.GetBalance)
		}
	}

	panic(http.Serve(listener, r))
}
