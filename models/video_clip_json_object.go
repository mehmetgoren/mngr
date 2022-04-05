package models

import (
	"time"
)

type DetectedObject struct {
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name"`
}

type DetectedImage struct {
	Id               string           `json:"id"`
	SourceId         string           `json:"source_id"`
	CreatedAt        string           `json:"created_at"`
	DetectedObjects  []DetectedObject `json:"detected_objects"`
	Base64Image      string           `json:"base64_image"`
	VideoClipEnabled bool             `json:"video_clip_enabled"`
}

type VideoClipJsonObject struct {
	DetectedImage    *DetectedImage `json:"detected_image"`
	FileName         string         `json:"file_name"`
	CreatedAt        string         `json:"created_at"`
	LastModified     string         `json:"last_modified"`
	Duration         int            `json:"duration"`
	CreatedAtTime    time.Time      `json:"-"`
	LastModifiedTime time.Time      `json:"-"`
}
