package render

import (
	e "brender/api/resource/common/err"
	"brender/util/validator"
	"encoding/json"
	"net/http"
)

// Create godoc
//
//	@summary		Render blend file
//	@description	Render blend file
//	@tags			render
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"Render form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/render [post]
func (a *API) Render(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.BadRequest(w, e.JsonDecodingFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		resp := validator.ToErrResponse(err)
		if resp == nil {
			e.ServerError(w, e.FormErrResponseFailure)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			a.logger.Error().Err(err).Msg("")
			e.ServerError(w, e.JsonEncodingFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	a.logger.Info().Msg("Rendering has started")
	w.WriteHeader(http.StatusCreated)
}
