package mng

import (
	"mngr/data"
)

type AlprMapper struct {
}

func (a *AlprMapper) Map(source *AlprEntity) *data.AlprDto {
	ret := &data.AlprDto{}
	ret.Id = source.Id.String()
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedPlate = &data.DetectedPlateDto{
		Plate:      source.DetectedPlate.Plate,
		Confidence: source.DetectedPlate.Confidence,
	}
	ret.ImageFileName = source.ImageFileName
	ret.VideoFile = &data.VideoFileDto{}
	if source.VideoFile != nil {
		ret.VideoFile.Name = source.VideoFile.Name
		ret.VideoFile.CreatedDate = source.VideoFile.CreatedDate.Time()
		ret.VideoFile.Duration = source.VideoFile.Duration
		ret.VideoFile.Merged = source.VideoFile.Merged
		ret.VideoFile.ObjectAppearsAt = source.VideoFile.ObjectAppearsAt
	}
	ret.AiClip = source.AiClip

	return ret
}
