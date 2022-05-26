package models

type FrTrainViewModel struct {
	Name       string   `json:"name"`
	ImagePaths []string `json:"image_paths"`
}

type FrTrainScreenshotViewModel struct {
	Name        string `json:"name"`
	Base64Image string `json:"base64_image"`
}

type FrTrainRename struct {
	NewName      string `json:"new_name"`
	OriginalName string `json:"original_name"`
}

type FrTrainName struct {
	Name string `json:"name"`
}
