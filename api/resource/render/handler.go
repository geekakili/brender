package render

import (
	e "brender/api/resource/common/err"
	"encoding/json"

	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
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
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}

	formData := make(map[string]interface{})
	for k, v := range r.Form {
		formData[k] = v[0]
	}
	form := &Form{}
	formDataJson, err := json.Marshal(formData)
	if err != nil {
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}

	if err := json.Unmarshal([]byte(formDataJson), form); err != nil {
		a.logger.Error().Err(err).Msg("")
		e.BadRequest(w, e.JsonDecodingFailure)
		return
	}

	// Get the file from the request
	file, _, err := r.FormFile("project_file")
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to process file")
		e.BadRequest(w, "Error processing file")
		return
	}
	defer file.Close()

	// Create a new file in the uploads directory
	fileName := uuid.New()
	filePath := fmt.Sprintf("./uploads/%s.blend", fileName)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}
	defer f.Close()

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}
	a.errChannel = make(chan int)
	go func() {
		a.runBlender(filePath)
	}()

	sig := <-a.errChannel
	if sig == 0 {
		w.WriteHeader(http.StatusCreated)
		return
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *API) runBlender(filePath string) {
	cmd := exec.Command("/Applications/blender.app/Contents/MacOS/blender", "-b", filePath, "-x", "1", "-o", "./test.png", "-a")

	outfile, err := os.Create("./log.txt")
	if err != nil {
		a.logger.Info().Msg(fmt.Sprintf("Failed to open log file: %s ", err.Error()))
		a.errChannel <- 1
		return
	}

	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		a.logger.Info().Msg("Blender instance has failed to start")
		a.errChannel <- 1
		return
	}
	a.logger.Info().Msg(fmt.Sprintf("A Blender instance is running with PID %d", cmd.Process.Pid))
	a.logger.Info().Msg("Render process started")
	a.errChannel <- 0
	a.logger.Info().Msg("Render process running")
	cmd.Wait()
}
