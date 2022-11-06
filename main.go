package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	commonutils "github.com/siddharthrajaraja/web_server_template_go/common_utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	goCtx, cancel := context.WithCancel(context.Background())

	initLogging()
	SetupConfiguration()

	defer func() {
		cancel()
	}()

	server := commonutils.NewServer()

	go func() {
		server.ServeHTTP()
	}()

	waitForShutdownSignal(server, goCtx, cancel)
}

func SetupConfiguration() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error while reading .env")
	}
}

func waitForShutdownSignal(srv *commonutils.Server, goCtx context.Context, cancel context.CancelFunc) {
	var gracefulStop = make(chan os.Signal, 3)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGQUIT)

	select {
	case <-gracefulStop:
		cancel()
		// if stop signal is received, wait for some time so that background workers get time to exit
		<-time.After(5 * time.Second)
	case <-goCtx.Done():
		// shutdown if context was cancelled by something else before shutdown signal
	}
	srv.Shutdown(goCtx)
}

func initLogging() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "debug"
	}
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)
}
