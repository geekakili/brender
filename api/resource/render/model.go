package render

type Form struct {
	StartFrame      int    `json:"start_frame" form:"default, 1"`
	EndFrame        int    `json:"end_frame" form:"default, 1"`
	FrameJump       int    `json:"frame_jump" form:"default, 0"`
	RenderFrames    string `json:"render_frames"`
	RenderEngine    string `json:"render_engine"`
	OutputFormat    string `json:"output_format"`
	RenderAnimation bool   `json:"render_animation" form:"default, false"`
}

type RenderObject struct {
	StartFrame      int
	EndFrame        int
	FrameJump       int
	RenderFrames    string
	RenderEngine    string
	OutputFormat    string
	RenderAnimation bool
}

func (renderForm *Form) ToModel() *RenderObject {

	return &RenderObject{
		StartFrame:      renderForm.StartFrame,
		EndFrame:        renderForm.EndFrame,
		FrameJump:       renderForm.FrameJump,
		RenderFrames:    renderForm.RenderFrames,
		RenderEngine:    renderForm.RenderEngine,
		OutputFormat:    renderForm.OutputFormat,
		RenderAnimation: renderForm.RenderAnimation,
	}
}
