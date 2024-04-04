package mng

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mngr/data"
)

type DetectedObject struct {
	Score float32 `json:"score"  bson:"score"`
	Label string  `json:"label" bson:"label"` //Index
	X1    float32 `json:"x1" bson:"x1"`
	Y1    float32 `json:"y1" bson:"y1"`
	X2    float32 `json:"x2" bson:"x2"`
	Y2    float32 `json:"y2" bson:"y2"`
}

type AiEntity struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Module         string             `json:"module" bson:"module"`       //Index
	GroupId        string             `json:"group_id" bson:"group_id"`   //Index
	SourceId       string             `json:"source_id" bson:"source_id"` //Index
	CreatedAt      string             `json:"created_at" bson:"created_at"`
	DetectedObject *DetectedObject    `json:"detected_object" bson:"detected_object"`
	ImageFileName  string             `json:"image_file_name" bson:"image_file_name"`

	VideoFile *VideoFile `json:"video_file" bson:"video_file"`

	AiClip *data.AiClip `json:"ai_clip" bson:"ai_clip"`

	CreatedDate primitive.DateTime `json:"created_date" bson:"created_date"`
}
