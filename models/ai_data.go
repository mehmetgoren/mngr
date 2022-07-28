package models

const (
	Od   = 0
	Fr   = 1
	Alpr = 2
)

const (
	Ascending  = 1
	Descending = -1
)

type SortInfo struct {
	Enabled bool   `json:"enabled"`
	Field   string `json:"field"`
	Sort    int    `json:"sort"`
}

func CreateDateSort(field string) SortInfo {
	return SortInfo{
		Enabled: true,
		Field:   field,
		Sort:    Descending,
	}
}

type PagingInfo struct {
	Enabled bool `json:"enabled"`
	Page    int  `json:"page"`
	Take    int  `json:"take"`
}

type QueryAiDataParams struct {
	AiType               int    `json:"ai_type"`
	SourceId             string `json:"source_id"`
	DateTimeStr          string `json:"date_time_str"`
	PredClassName        string `json:"pred_class_name"`
	NoPreparingVideoFile bool   `json:"no_preparing_video_file"`
}

type QueryAiDataAdvancedParams struct {
	AiType               int        `json:"ai_type"`
	SourceId             string     `json:"source_id"`
	StartDateTimeStr     string     `json:"start_date_time_str"`
	EndDateTimeStr       string     `json:"end_date_time_str"`
	PredClassName        string     `json:"pred_class_name"`
	NoPreparingVideoFile bool       `json:"no_preparing_video_file"`
	Sort                 SortInfo   `json:"sort"`
	Paging               PagingInfo `json:"paging"`
}

type AiDataDeleteOptions struct {
	AiType      int    `json:"ai_type"`
	Id          string `json:"id"`
	DeleteImage bool   `json:"delete_image"`
	DeleteVideo bool   `json:"delete_video"`
}
