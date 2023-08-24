package render

import (
	"github.com/go-playground/validator/v10"

	"brender/util/logger"
)

type API struct {
	logger    *logger.Logger
	validator *validator.Validate
}

func New(logger *logger.Logger, validator *validator.Validate) *API {
	return &API{
		logger:    logger,
		validator: validator,
	}
}
