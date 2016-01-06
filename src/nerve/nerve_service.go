package nerve

import (
	log "github.com/Sirupsen/logrus"
	"time"
	"strconv"
)

type NerveService struct {
	Watcher WatcherI
	Reporter ReporterI
	Name string
	Host string
	Port int
}

func(ns *NerveService) Initialize(config NerveServiceConfiguration) error {
	var err error

	ns.Name = config.Name
	ns.Host = config.Host
	ns.Port = config.Port
	log.Debug("Service [",ns.Name,"] for Host [",ns.Host,"] Port [",ns.Port,"] initialisation")
	ns.Watcher , err = CreateWatcher(config.Watcher)
	if err != nil {
		log.Warn("Error creating Watcher in Service [",ns.Name,"]")
		return err
	}
	ns.Reporter, err = CreateReporter(config.Reporter)
	if err != nil {
		log.Warn("Error creating Reporter in Service [",ns.Name,"]")
		return err
	}
	return nil
}

func(ns *NerveService) Run(stop <-chan bool) {
	defer servicesWaitGroup.Done()
	log.Debug("Service Running [",ns.Name,"]")
	Loop:
	for {
		// Here The job to check, and report
		status, err := ns.Watcher.Check()
		if err != nil  {
			log.Warn("Check error for Service [", ns.Name, "] [",err,"]")
		}
		ns.Reporter.Report("127.0.0.1",strconv.Itoa(ns.Port),ns.Host,status);

		// Wait for the stop signal
		select {
		case hasToStop := <-stop:
			if hasToStop {
				log.Debug("Nerve: Service [",ns.Name,"]Run Close Signal Received")
			}else {
				log.Debug("Nerve: Service [",ns.Name,"]Run Close Signal Received (but a strange false one)")
			}
			break Loop
		default:
			time.Sleep(time.Second * time.Duration(ns.Watcher.GetCheckInterval()))
		}
	}
	log.Debug("Service [",ns.Name,"] stopped")
}

func CreateService(config NerveServiceConfiguration) (NerveService, error) {
	var service NerveService
	err := service.Initialize(config)
	if err != nil {
		log.Debug("Error Initializing Service [",service.Name,"]")
	}
	return service, err
}
