package main

import (
	"flag"
	"io/ioutil"

	"github.com/JamesOwenHall/timed/server"

	"github.com/Sirupsen/logrus"
)

func main() {
	log := logrus.StandardLogger()

	configFileName := flag.String("c", "", "path to the configuration file")
	flag.Parse()

	if *configFileName == "" {
		log.Fatal("no configuration file specified")
	}

	configData, err := ioutil.ReadFile(*configFileName)
	if err != nil {
		log.Fatal("error reading configuration file: ", err.Error())
	}

	config, err := server.NewConfigFromYAML(configData)
	if err != nil {
		log.Fatal("error parsing configuration file: ", err.Error())
	}

	s, err := server.NewServer(log, config)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
