package entity

// Input Data define the struct for Input
type InputData struct {
	OriginData    map[string]interface{} `json:"origin_data"`
	TargetData    map[string]interface{} `json:"target_data"`
	ParsingFormat []Format               `json:"parsing_format"`
}

// ParsingFormat defines the structure for the parsing format JSON
type Format struct {
	Origin string `json:"origin"`
	Target string `json:"target"`
	Format string `json:"format,omitempty"`
}
