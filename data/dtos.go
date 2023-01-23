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

type ColorDto struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type MetadataDto struct {
	Colors []ColorDto `json:"colors"`
}

type DetectedObjectDto struct {
	PredScore   float32      `json:"pred_score"`
	PredClsIdx  int          `json:"pred_cls_idx"`
	PredClsName string       `json:"pred_cls_name"`
	X1          float32      `json:"x1" bson:"x1"`
	Y1          float32      `json:"y1" bson:"y1"`
	X2          float32      `json:"x2" bson:"x2"`
	Y2          float32      `json:"y2" bson:"y2"`
	Metadata    *MetadataDto `json:"metadata" bson:"metadata"`
}

type OdDto struct {
	Id             string             `json:"id"`
	GroupId        string             `json:"group_id"`
	SourceId       string             `json:"source_id"`
	CreatedAt      string             `json:"created_at"`
	DetectedObject *DetectedObjectDto `json:"detected_object"`
	ImageFileName  string             `json:"image_file_name"`
	VideoFile      *VideoFileDto      `json:"video_file"`
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
	VideoFile     *VideoFileDto    `json:"video_file"`
	AiClip        *AiClip          `json:"ai_clip"`
}

type DetectedPlateDto struct {
	Plate            string  `json:"plate"`
	Confidence       float64 `json:"confidence"`
	ProcessingTimeMs float64 `json:"processing_time_ms"`
	X1               int     `json:"x1" bson:"x1"`
	Y1               int     `json:"y1" bson:"y1"`
	X2               int     `json:"x2" bson:"x2"`
	Y2               int     `json:"y2" bson:"y2"`
}

type AlprDto struct {
	Id            string            `json:"id"`
	GroupId       string            `json:"group_id"`
	SourceId      string            `json:"source_id"`
	CreatedAt     string            `json:"created_at"`
	DetectedPlate *DetectedPlateDto `json:"detected_plate"`
	ImageFileName string            `json:"image_file_name"`
	VideoFile     *VideoFileDto     `json:"video_file"`
	AiClip        *AiClip           `json:"ai_clip"`
}

type AiDataDto struct {
	Id            string        `json:"id"`
	GroupId       string        `json:"group_id"`
	SourceId      string        `json:"source_id"`
	CreatedAt     string        `json:"created_at"`
	PredClsName   string        `json:"pred_cls_name"`
	PredScore     float32       `json:"pred_score"`
	ImageFileName string        `json:"image_file_name"`
	VideoFile     *VideoFileDto `json:"video_file"`
	AiClip        *AiClip       `json:"ai_clip"`
}

func (a *AiDataDto) MapFromOd(ods *OdDto) *AiDataDto {
	a.Id = ods.Id
	a.GroupId = ods.GroupId
	a.SourceId = ods.SourceId
	a.CreatedAt = ods.CreatedAt
	a.PredClsName = ods.DetectedObject.PredClsName
	a.PredScore = ods.DetectedObject.PredScore
	a.ImageFileName = ods.ImageFileName
	a.VideoFile = ods.VideoFile
	a.AiClip = ods.AiClip
	return a
}

func (a *AiDataDto) MapFromFr(frs *FrDto) *AiDataDto {
	a.Id = frs.Id
	a.GroupId = frs.GroupId
	a.SourceId = frs.SourceId
	a.CreatedAt = frs.CreatedAt
	a.PredClsName = frs.DetectedFace.PredClsName
	a.PredScore = frs.DetectedFace.PredScore
	a.ImageFileName = frs.ImageFileName
	a.VideoFile = frs.VideoFile
	a.AiClip = frs.AiClip
	return a
}

func (a *AiDataDto) MapFromAlpr(alprs *AlprDto) *AiDataDto {
	a.Id = alprs.Id
	a.GroupId = alprs.GroupId
	a.SourceId = alprs.SourceId
	a.CreatedAt = alprs.CreatedAt
	a.PredClsName = alprs.DetectedPlate.Plate
	a.PredScore = float32(alprs.DetectedPlate.Confidence)
	a.ImageFileName = alprs.ImageFileName
	a.VideoFile = alprs.VideoFile
	a.AiClip = alprs.AiClip
	return a
}
