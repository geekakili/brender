package render

type Form struct {
	StartFrame      int    `json:"start_frame,string"`
	EndFrame        int    `json:"end_frame,string"`
	FrameJump       int    `json:"frame_jump,string"`
	RenderFrames    string `json:"render_frames"`
	RenderEngine    string `json:"render_engine"`
	OutputFormat    string `json:"output_format"`
	RenderAnimation bool   `json:"render_animation,string"`
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
