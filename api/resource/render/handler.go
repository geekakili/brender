package render

import (
	e "brender/api/resource/common/err"
	"brender/api/resource/common/utilities"
	"encoding/json"
	"strconv"
	"time"

	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/dgraph-io/badger/v4"
	"github.com/lithammer/shortuuid/v4"
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
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}

	formData := make(map[string]interface{})
	for k, v := range r.Form {
		formData[k] = v[0]
	}
	renderObject := &RenderObject{}
	formDataJson, err := json.Marshal(formData)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}

	if err := json.Unmarshal([]byte(formDataJson), renderObject); err != nil {
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
	renderUUId := shortuuid.New()
	renderDir := fmt.Sprintf("./uploads/%s", renderUUId)
	err = utilities.EnsureDir(renderDir)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}

	filePath := fmt.Sprintf("%s/project.blend", renderDir)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}
	defer f.Close()

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		a.logger.Error().Err(err).Msg("")
		e.ServerError(w, e.FormErrResponseFailure)
		return
	}

	if a.isBusy {
		// Add a render job to a queue
		// Send a signal to start the next job in line
		// Remove finished jobs from the queue
	}
	a.errChannel = make(chan int)
	go func() {
		renderMetadata := new(RenderMetadata)
		renderMetadata.RenderObject = *renderObject
		renderMetadata.RenderDirectory = renderDir
		renderMetadata.StartTime = time.Now().UnixNano()
		renderMetadataBytes, err := json.Marshal(renderMetadata)
		if err != nil {
			a.logger.Error().Err(err).Msg("")
			a.errChannel <- 1
			return
		}

		err = a.db.Update(func(txn *badger.Txn) error {
			err := txn.Set([]byte(renderUUId), renderMetadataBytes)
			return err
		})
		if err != nil {
			a.logger.Error().Err(err).Msg("")
			a.errChannel <- 1
			return
		}
		a.runBlender(renderMetadata)
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

func (a *API) runBlender(blenderMetadata *RenderMetadata) {
	a.isBusy = true
	blendFilePath := fmt.Sprintf("%s/project.blend", blenderMetadata.RenderDirectory)
	output := fmt.Sprintf("%s/render_####", blenderMetadata.RenderDirectory)
	blenderCliArgs := []string{"-b", blendFilePath, "-x", "1", "-o", output}
	if blenderMetadata.RenderObject.StartFrame != 0 {
		blenderCliArgs = append(blenderCliArgs, "-s", strconv.Itoa(blenderMetadata.RenderObject.StartFrame))
	}

	if blenderMetadata.RenderObject.EndFrame != 0 {
		blenderCliArgs = append(blenderCliArgs, "-e", strconv.Itoa(blenderMetadata.RenderObject.EndFrame))
	}

	if blenderMetadata.RenderObject.FrameJump > 0 {
		blenderCliArgs = append(blenderCliArgs, "-j", strconv.Itoa(blenderMetadata.RenderObject.FrameJump))
	}

	if blenderMetadata.RenderObject.RenderAnimation {
		blenderCliArgs = append(blenderCliArgs, "-a")
	} else {
		blenderCliArgs = append(blenderCliArgs, "-f", blenderMetadata.RenderObject.RenderFrames)
	}
	cmd := exec.Command("/Applications/blender.app/Contents/MacOS/blender", blenderCliArgs...)

	logFilePath := fmt.Sprintf("%s/logs.txt", blenderMetadata.RenderDirectory)
	outfile, err := os.Create(logFilePath)
	if err != nil {
		a.logger.Error().Msg(fmt.Sprintf("Failed to open log file: %s ", err.Error()))
		a.errChannel <- 1
		return
	}

	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		a.logger.Error().Msg("Blender instance has failed to start")
		a.errChannel <- 1
		return
	}
	a.logger.Info().Msg(fmt.Sprintf("A Blender instance is running with PID %d", cmd.Process.Pid))
	a.logger.Info().Msg("Render process started")
	a.errChannel <- 0
	a.logger.Info().Msg("Render process running")
	err = cmd.Wait()
	if err != nil {
		a.logger.Error().Msg("Blender encountered an error while rendering")
		a.logger.Info().Msg(err.Error())
		a.isBusy = false
		return
	}
	a.logger.Info().Msg("Render process has completed")
	a.isBusy = false
}

func prepareRender(renderObject RenderObject) RenderMetadata {

}
