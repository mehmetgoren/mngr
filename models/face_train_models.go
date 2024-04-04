package models

type FaceTrainViewModel struct {
	Name       string   `json:"name"`
	ImagePaths []string `json:"image_paths"`
}

type FaceTrainScreenshotViewModel struct {
	Name         string   `json:"name"`
	Base64Images []string `json:"base64_images"`
}

type FaceTrainRename struct {
	NewName      string `json:"new_name"`
	OriginalName string `json:"original_name"`
}

type FaceTrainName struct {
	Name string `json:"name"`
}
