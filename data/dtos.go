package data

type AiClip struct {
	Enabled        bool   `json:"enabled" bson:"enabled"`
	FileName       string `json:"file_name" bson:"file_name"`
	CreatedAt      string `json:"created_at" bson:"created_at"`
	LastModifiedAt string `json:"last_modified_at" bson:"last_modified_at"`
	Duration       int    `json:"duration" bson:"duration"`
}

type VideoFileDto struct {
	Name            string `json:"name"`
	CreatedAt       string `json:"created_at"`
	Duration        int    `json:"duration"`
	Merged          bool   `json:"merged"`
	ObjectAppearsAt int    `json:"object_appears_at"`
}

type DetectedObjectDto struct {
	Score float32 `json:"score"`
	Label string  `json:"label"`
	X1    float32 `json:"x1" bson:"x1"`
	Y1    float32 `json:"y1" bson:"y1"`
	X2    float32 `json:"x2" bson:"x2"`
	Y2    float32 `json:"y2" bson:"y2"`
}

type AiDto struct {
	Id             string             `json:"id"`
	Module         string             `json:"module"`
	GroupId        string             `json:"group_id"`
	SourceId       string             `json:"source_id"`
	CreatedAt      string             `json:"created_at"`
	DetectedObject *DetectedObjectDto `json:"detected_object"`
	ImageFileName  string             `json:"image_file_name"`
	VideoFile      *VideoFileDto      `json:"video_file"`
	AiClip         *AiClip            `json:"ai_clip"`
}
