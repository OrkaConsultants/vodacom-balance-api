package core

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"github.com/OrkaConsultants/vodacom-balance-api/internal/models"
)

var configPath = "."
var bootstrapConfigName = "bootstrap"
var defaultAppConfigName = "vodacom-api"

func LoadConfig() {
	readLocalConfig(bootstrapConfigName)

	if viper.GetBool("cloud-config.enabled") {
		log.Info("Loading cloud config")
		loadCloudConfig()
	} else {
		log.Info("Loading local config")
		readLocalConfig(defaultAppConfigName)
	}

	// Test if the service username has been set to ensure config has been loaded
	if viper.IsSet("vodacom-api.username") {
		log.Infof("Successfully loaded configuration for service %s", viper.GetString("app.name"))
		return
	}

	log.Fatal("Couldn't load configuration, cannot start. Terminating.")
}

func readLocalConfig(configName string) {
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.AddConfigPath(configPath)
	err := viper.MergeInConfig() // Find and merge the config file
	if err != nil {              // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s", err)
	}
}

func getCloudConfigServiceURL() string {
	eurekaServer := viper.GetString("eureka.url") + ":" + viper.GetString("eureka.port")
	configServiceURL := eurekaServer + "/eureka/apps/" + viper.GetString("cloud-config.service-name")

	body, err := fetchConfiguration(configServiceURL)

	var configServiceStruct models.EurekaAppXMLStruct
	err = xml.Unmarshal(body, &configServiceStruct)
	if err != nil {
		panic("Cannot parse configuration, message: " + err.Error())
	}

	return "http://" + configServiceStruct.Instance.IpAddr + ":" + configServiceStruct.Instance.Port.Text
}

func loadCloudConfig() {
	cloudConfigURL := getCloudConfigServiceURL()
	appName := viper.GetString("app.name")
	profile := viper.GetString("cloud-config.profile")
	branch := ""

	url := fmt.Sprintf("%s/%s/%s/%s", cloudConfigURL, appName, profile, branch)
	log.Infof("Loading config from %s", url)
	body, err := fetchConfiguration(url)
	if err != nil {
		panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}
	parseConfiguration(body)
}

// Make HTTP request to fetch configuration from config server
func fetchConfiguration(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Pass JSON bytes into struct and then into Viper
func parseConfiguration(body []byte) {
	var cloudConfig models.SpringCloudConfig
	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		panic("Cannot parse configuration, message: " + err.Error())
	}

	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
		log.Debugf("Loading config property %v => %v\n", key, value)
	}
}
