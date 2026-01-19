package web

import (
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type DefaultController struct{}

func (ctrl DefaultController) HTTPError(statusCode int, message string) *fuego.HTTPError {
	return &fuego.HTTPError{
		Status: statusCode,
		Detail: message,
	}
}

func (ctrl DefaultController) ParamID(req *http.Request) (entities.HexID, error) {
	idStr := req.PathValue("id")
	if idStr == "" {
		return entities.HexID{}, ctrl.HTTPError(http.StatusBadRequest, "id is required")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return entities.HexID{}, ctrl.HTTPError(http.StatusBadRequest, "invalid id format")
	}

	return entities.HexID(id), nil
}
