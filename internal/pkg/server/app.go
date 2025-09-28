package server

import (
	"fmt"

	"github.com/g0shi4ek/RIP_backend/internal/app/config"
	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
   Config  *config.Config
   Router  *gin.Engine
   Handler domain.IChargingHandler
}

func NewApp(c *config.Config, r *gin.Engine, h domain.IChargingHandler) *Application {
   return &Application{
      Config:  c,
      Router:  r,
      Handler: h,
   }
}

func (a *Application) RunApp() {
   logrus.Info("Server start up")

   a.Handler.RegisterChargingHandler(a.Router)
   
   serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
   if err := a.Router.Run(serverAddress); err != nil {
      logrus.Fatal(err)
   }
   logrus.Info("Server down")
}