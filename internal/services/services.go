package services

import (
	"strings"

	"github.com/spf13/viper"
)

type VodacomUriService struct{}

func (u VodacomUriService) LoginUri() string {
	return viper.GetString("vodacom-api.login-uri")
}

func (u VodacomUriService) BalanceFor(number string) string {
	return strings.Replace(viper.GetString("vodacom-api.balance-uri"), "{number}", number, -1)
}
