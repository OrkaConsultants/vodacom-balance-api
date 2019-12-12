package services

import (
	"encoding/json"
	"io/ioutil"

	"github.com/OrkaConsultants/vodacom-balance-api/internal/models"
	"github.com/OrkaConsultants/vodacom-balance-api/internal/utils"
	log "github.com/sirupsen/logrus"
)

type BalanceService struct{}

func (u BalanceService) GetBalance(number string) (*[]models.ServiceList, error) {
	log.Info("Balance: " + number)
	route := VodacomUriService.BalanceFor(VodacomUriService{}, number)

	req, err := utils.NewRequest("GET", route, nil)
	if err != nil {
		return nil, err
	}

	client := utils.NewClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var balanceResponse models.BalanceResponse
	errJSON := json.Unmarshal(body, &balanceResponse)

	if errJSON != nil {
		log.Errorf("%s", errJSON.Error())
		return nil, err
	}

	return &balanceResponse.Result.ServiceList, nil
}
