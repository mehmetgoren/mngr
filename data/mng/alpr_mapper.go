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
		Plate:            source.DetectedPlate.Plate,
		Confidence:       utils.RoundFloat64(source.DetectedPlate.Confidence),
		ProcessingTimeMs: source.DetectedPlate.ProcessingTimeMs,
		X1:               source.DetectedPlate.Coordinates.X0,
		Y1:               source.DetectedPlate.Coordinates.Y0,
		X2:               source.DetectedPlate.Coordinates.X1,
		Y2:               source.DetectedPlate.Coordinates.Y1,
	}
	ret.ImageFileName = source.ImageFileName
	ret.VideoFile = &data.VideoFileDto{}
	if source.VideoFile != nil {
		ret.VideoFile.Name = source.VideoFile.Name
		ret.VideoFile.CreatedAt = utils.TimeToString(source.VideoFile.CreatedDate.Time(), false)
		ret.VideoFile.Duration = source.VideoFile.Duration
		ret.VideoFile.Merged = source.VideoFile.Merged
		ret.VideoFile.ObjectAppearsAt = source.VideoFile.ObjectAppearsAt
	}
	ret.AiClip = source.AiClip
	ret.AiClip.FileName = source.AiClip.FileName

	return ret
}
