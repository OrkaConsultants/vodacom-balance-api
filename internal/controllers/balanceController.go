package controllers

import (
	"net/http"

	"github.com/OrkaConsultants/vodacom-balance-api/internal/services"
	"github.com/gin-gonic/gin"
)

type BalanceController struct{}

var balanceService = new(services.BalanceService)

func (u BalanceController) GetBalance(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	if params["number"] != nil {
		number := params["number"][0]
		result, err := balanceService.GetBalance(number)
		_ = result
		_ = err
		ctx.JSON(200, result)
		ctx.Abort()
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Provide number parameter i.e. /api/balance?number=2782000000"})
		ctx.Abort()
		return
	}
}
