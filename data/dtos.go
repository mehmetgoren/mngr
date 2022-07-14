package data

type AiClip struct {
	Enabled        bool   `json:"enabled" bson:"enabled"`
	FileName       string `json:"file_name" bson:"file_name"` //Index
	CreatedAt      string `json:"created_at" bson:"created_at"`
	LastModifiedAt string `json:"last_modified_at" bson:"last_modified_at"`
	Duration       int    `json:"duration" bson:"duration"`
}

type DetectedObjectDto struct {
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name"`
}

type OdDto struct {
	Id             string             `json:"id"`
	GroupId        string             `json:"group_id"`
	SourceId       string             `json:"source_id"`
	CreatedAt      string             `json:"created_at"`
	DetectedObject *DetectedObjectDto `json:"detected_object"`
	ImageFileName  string             `json:"image_file_name"`
	VideoFileName  string             `json:"video_file_name"`
	AiClip         *AiClip            `json:"ai_clip"`
}

type DetectedFaceDto struct {
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name"`
	X1          float32 `json:"x1"`
	Y1          float32 `json:"y1"`
	X2          float32 `json:"x2"`
	Y2          float32 `json:"y2"`
}

type FrDto struct {
	Id            string           `json:"id"`
	GroupId       string           `json:"group_id"`
	SourceId      string           `json:"source_id"`
	CreatedAt     string           `json:"created_at"`
	DetectedFace  *DetectedFaceDto `json:"detected_face"`
	ImageFileName string           `json:"image_file_name"`
	VideoFileName string           `json:"video_file_name"`
	AiClip        *AiClip          `json:"ai_clip"`
}

type DetectedPlateDto struct {
	Plate            string  `json:"plate"`
	Confidence       float64 `json:"confidence"`
	ProcessingTimeMs float64 `json:"processing_time_ms"`
}

type AlprDto struct {
	Id            string            `json:"id"`
	GroupId       string            `json:"group_id"`
	SourceId      string            `json:"source_id"`
	CreatedAt     string            `json:"created_at"`
	DetectedPlate *DetectedPlateDto `json:"detected_plate"`
	ImageFileName string            `json:"image_file_name"`
	VideoFileName string            `json:"video_file_name"`
	AiClip        *AiClip           `json:"ai_clip"`
}
