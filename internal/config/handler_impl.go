package config

import (
	"github.com/containerssh/containerssh/config"
	"github.com/containerssh/containerssh/http"
	"github.com/containerssh/containerssh/log"
)

type handler struct {
	handler RequestHandler
	logger  log.Logger
}

func (h *handler) OnRequest(request http.ServerRequest, response http.ServerResponse) error {
	requestObject := config.ConfigRequest{}
	if err := request.Decode(&requestObject); err != nil {
		return err
	}
	appConfig, err := h.handler.OnConfig(requestObject)
	if err != nil {
		return err
	}
	responseObject := config.ConfigResponseBody{
		Config: appConfig,
	}
	response.SetBody(responseObject)
	response.SetStatus(200)
	return nil
}