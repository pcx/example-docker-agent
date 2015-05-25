package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"

	"github.com/pcx/st-agent/conf"
	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/server"
)

func main() {
	app := cli.NewApp()
	app.Name = "example-docker-agent"
	app.Usage = "Manage Docker like a ninja!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "MachineID, m",
			Value: "",
			Usage: "Authentication & identification for the controlled machine",
		},
		cli.StringFlag{
			Name:  "HubURL, u",
			Value: "",
			Usage: "Used to connect to a Central Hub, format: scheme://host:port",
		},
		cli.StringFlag{
			Name:   "DockerHost, d",
			Value:  "unix:///var/run/docker.sock",
			Usage:  "Host address to connect to Docker daemon",
			EnvVar: "DOCKER_HOST",
		},
		cli.BoolFlag{
			Name:   "DockerTLSVerify, s",
			Usage:  "Connect to Docker daemon over TLS, needs --DockerCertPath",
			EnvVar: "DOCKER_TLS_VERIFY",
		},
		cli.StringFlag{
			Name:  "DockerCertPath, p",
			Value: "",
			Usage: "Path to dir with TLS certs for connecting to Docker daemon," +
				" needs --DockerTLSverify",
			EnvVar: "DOCKER_CERT_PATH",
		},
	}

	app.Action = func(ctx *cli.Context) {
		config, err := conf.MakeConfig(ctx)
		if err != nil {
			log.Error(err)
			fmt.Println()

			cli.ShowAppHelp(ctx)
			os.Exit(1)
		}
		log.Debug("Starting example-docker-agent...")
		// log.EnableTimestamps()
		log.EnableDebug()

		s := server.NewServer(config)
		s.Start()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

		<-sigChan
		log.Infof("Gracefully shutting down...")
		s.Stop()
	}

	app.Run(os.Args)
}
