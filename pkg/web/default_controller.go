package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type DefaultController struct{}

func (ctrl DefaultController) HTTPError(statusCode int, err error) *fuego.HTTPError {
	if err == nil {
		err = errors.New(http.StatusText(statusCode))
	}

	return &fuego.HTTPError{
		Err:    err,
		Title:  http.StatusText(statusCode),
		Status: statusCode,
		Detail: err.Error(),
	}
}

func (ctrl DefaultController) HandleError(err error) error {
	return MapError(err)
}

func (ctrl DefaultController) ParamID(req *http.Request) (entities.HexID, error) {
	idStr := req.PathValue("id")
	if idStr == "" {
		return entities.HexID{}, ctrl.HTTPError(http.StatusBadRequest, errors.New("id is required"))
	}

	return ctrl.ParseHexID(idStr)
}

func (ctrl DefaultController) ParseHexID(rawValue string) (entities.HexID, error) {
	id, err := entities.ParseHexID(rawValue)
	if err != nil {
		return entities.HexID{}, ctrl.HTTPError(
			http.StatusBadRequest,
			fmt.Errorf("invalid id format: %w", err),
		)
	}

	return id, nil
}
