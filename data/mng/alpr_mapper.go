package mng

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
)

type AlprMapper struct {
	Config *models.Config
}

func (a *AlprMapper) Map(source *AlprEntity) *data.AlprDto {
	ret := &data.AlprDto{}
	ret.Id = source.Id.Hex()
	ret.GroupId = source.GroupId
	ret.SourceId = source.SourceId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedPlate = &data.DetectedPlateDto{
		Plate:      source.DetectedPlate.Plate,
		Confidence: utils.RoundFloat64(source.DetectedPlate.Confidence),
	}
	ret.ImageFileName = utils.SetRelativeImagePath(a.Config, source.ImageFileName)
	ret.VideoFile = &data.VideoFileDto{}
	if source.VideoFile != nil {
		ret.VideoFile.Name = utils.SetRelativeRecordPath(a.Config, source.VideoFile.Name)
		ret.VideoFile.CreatedAt = utils.TimeToString(source.VideoFile.CreatedDate.Time(), false)
		ret.VideoFile.Duration = source.VideoFile.Duration
		ret.VideoFile.Merged = source.VideoFile.Merged
		ret.VideoFile.ObjectAppearsAt = source.VideoFile.ObjectAppearsAt
	}
	ret.AiClip = source.AiClip
	ret.AiClip.FileName = utils.SetRelativeOdAiVideoClipPath(a.Config, source.AiClip.FileName)

	return ret
}
