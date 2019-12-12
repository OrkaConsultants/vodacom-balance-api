package core

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

func SetupEureka() {

	if !viper.GetBool("eureka.enabled") {
		return
	}

	client := eureka.NewClient([]string{
		viper.GetString("eureka.url") + ":" + viper.GetString("eureka.port") + "/eureka", //From a spring boot based eureka server
		// add others servers here
	})

	instance := eureka.NewInstanceInfo( //Create a new instance to register
		GetOutboundIP().String(),    //Hostname
		viper.GetString("app.name"), //App name
		GetOutboundIP().String(),    //Ip Address
		viper.GetInt("app.port"),    //port
		30,                          //ttl
		false)                       //ssl

	instance.InstanceID = instance.HostName + ":" + instance.App + ":" + viper.GetString("app.port")
	instance.VipAddress = instance.App
	instance.SecureVipAddress = instance.App

	instance.StatusPageUrl = viper.GetString("api.gatekeeper") + "/" + instance.App + "/swagger/index.html"

	// instance.Metadata = &eureka.MetaData{
	// 	Map: make(map[string]string),
	// }
	// instance.Metadata.Map["foo"] = "bar" //add metadata for example

	client.RegisterInstance(instance.App, instance) // Register new instance in your eureka(s)

	// applications, _ := client.GetApplications() // Retrieves all applications from eureka server(s)

	client.GetApplication(instance.App)                 // retrieve the application
	client.GetInstance(instance.App, instance.HostName) // retrieve the instance

	log.Info("Finished Eureka registration")

	go continueHeartBeat(instance, client)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		log.Infof("Caught signal: %+v. Unregistering Eureka instance", s)
		client.UnregisterInstance(instance.App, instance.InstanceID)
		os.Exit(0)
	}()
}

func continueHeartBeat(instance *eureka.InstanceInfo, client *eureka.Client) {
	for {
		if client.SendHeartbeat(instance.App, instance.InstanceID) != nil { // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
			log.Error("Heartbeat send failed, restarting service in 90 seconds")
			time.Sleep(time.Second * 90)
			os.Exit(0)
		}
		time.Sleep(time.Second * 10)
	}
}
