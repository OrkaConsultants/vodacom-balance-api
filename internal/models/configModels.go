package models

import "encoding/xml"

// Structs having same structure as response from Spring Cloud Config
type SpringCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []PropertySource `json:"propertySources"`
}

type PropertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

type EurekaAppXMLStruct struct {
	XMLName  xml.Name `xml:"application"`
	Text     string   `xml:",chardata"`
	Name     string   `xml:"name"`
	Instance struct {
		Text             string `xml:",chardata"`
		InstanceId       string `xml:"instanceId"`
		HostName         string `xml:"hostName"`
		App              string `xml:"app"`
		IpAddr           string `xml:"ipAddr"`
		Status           string `xml:"status"`
		Overriddenstatus string `xml:"overriddenstatus"`
		Port             struct {
			Text    string `xml:",chardata"`
			Enabled string `xml:"enabled,attr"`
		} `xml:"port"`
		SecurePort struct {
			Text    string `xml:",chardata"`
			Enabled string `xml:"enabled,attr"`
		} `xml:"securePort"`
		CountryId      string `xml:"countryId"`
		DataCenterInfo struct {
			Text  string `xml:",chardata"`
			Class string `xml:"class,attr"`
			Name  string `xml:"name"`
		} `xml:"dataCenterInfo"`
		LeaseInfo struct {
			Text                  string `xml:",chardata"`
			RenewalIntervalInSecs string `xml:"renewalIntervalInSecs"`
			DurationInSecs        string `xml:"durationInSecs"`
			RegistrationTimestamp string `xml:"registrationTimestamp"`
			LastRenewalTimestamp  string `xml:"lastRenewalTimestamp"`
			EvictionTimestamp     string `xml:"evictionTimestamp"`
			ServiceUpTimestamp    string `xml:"serviceUpTimestamp"`
		} `xml:"leaseInfo"`
		Metadata struct {
			Text  string `xml:",chardata"`
			Class string `xml:"class,attr"`
		} `xml:"metadata"`
		HomePageUrl                   string `xml:"homePageUrl"`
		StatusPageUrl                 string `xml:"statusPageUrl"`
		HealthCheckUrl                string `xml:"healthCheckUrl"`
		VipAddress                    string `xml:"vipAddress"`
		SecureVipAddress              string `xml:"secureVipAddress"`
		IsCoordinatingDiscoveryServer string `xml:"isCoordinatingDiscoveryServer"`
		LastUpdatedTimestamp          string `xml:"lastUpdatedTimestamp"`
		LastDirtyTimestamp            string `xml:"lastDirtyTimestamp"`
		ActionType                    string `xml:"actionType"`
	} `xml:"instance"`
}
