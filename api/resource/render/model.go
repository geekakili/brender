package render

type RenderObject struct {
	StartFrame      int    `json:"start_frame,string"`
	EndFrame        int    `json:"end_frame,string"`
	FrameJump       int    `json:"frame_jump,string"`
	RenderFrames    string `json:"render_frames"`
	RenderEngine    string `json:"render_engine"`
	OutputFormat    string `json:"output_format"`
	RenderAnimation bool   `json:"render_animation,string"`
}

type RenderMetadata struct {
	RenderObject    RenderObject
	RenderDirectory string
	StartTime       int64
	StopTime        int64
}
