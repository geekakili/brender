package render

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/go-playground/validator/v10"

	"brender/util/logger"
)

type API struct {
	logger     *logger.Logger
	validator  *validator.Validate
	errChannel chan int
	db         *badger.DB
	isBusy     bool
}

func New(logger *logger.Logger, validator *validator.Validate, db *badger.DB) *API {
	return &API{
		logger:    logger,
		validator: validator,
		db:        db,
	}
}
