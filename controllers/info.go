package controllers

import (
	"net/http"

	"github.com/ditrit/badaas/httperrors"
)

const Version = "UNRELEASED"

// The information controller
type InformationController interface {
	// Return the badaas server information
	Info(response http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// check interface compliance
var _ InformationController = (*infoControllerImpl)(nil)

// The InformationController constructor
func NewInfoController() InformationController {
	return &infoControllerImpl{}
}

// The concrete implementation of the InformationController
type infoControllerImpl struct{}

// Return the badaas server information
func (*infoControllerImpl) Info(response http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {

	infos := &BadaasServerInfo{
		Status:  "OK",
		Version: Version,
	}
	return infos, nil
}

type BadaasServerInfo struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
