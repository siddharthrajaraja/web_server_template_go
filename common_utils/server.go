package commonutils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	r       *mux.Router
	address string
	srv     *http.Server
}

func NewServer() *Server {
	return &Server{
		r:       mux.NewRouter(),
		address: fmt.Sprintf("%s:%s", viper.GetString("SERVER_HOST"), viper.GetString("SERVER_PORT")),
	}
}

func (s *Server) ServeHTTP() {
	s.srv = &http.Server{
		Handler: s.r,
		Addr:    s.address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Minute,
		ReadTimeout:  time.Minute,
	}

	logrus.Info("Server starting at addr: ", s.address)
	err := s.srv.ListenAndServe()

	if err == http.ErrServerClosed {
		logrus.Info("Server shut dowon")
	} else if err != nil {
		logrus.Error(err)
	}

}

func (s *Server) Shutdown(goCtx context.Context) {
	if s.srv == nil {
		logrus.Error("Server not yet started.")
		return
	}

	if err := s.srv.Shutdown(goCtx); err != nil {
		logrus.Fatalf("Server failed to shutdown: %v", err.Error())
	}
}
