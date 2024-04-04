package mng

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
)

type AiMapper struct {
	Config *models.Config
}

func (a *AiMapper) Map(source *AiEntity) *data.AiDto {
	ret := &data.AiDto{}
	ret.Id = source.Id.Hex()
	ret.Module = source.Module
	ret.GroupId = source.GroupId
	ret.SourceId = source.SourceId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedObject = &data.DetectedObjectDto{
		Score: utils.RoundFloat32(source.DetectedObject.Score),
		Label: source.DetectedObject.Label,
		X1:    utils.RoundFloat32(source.DetectedObject.X1),
		Y1:    utils.RoundFloat32(source.DetectedObject.Y1),
		X2:    utils.RoundFloat32(source.DetectedObject.X2),
		Y2:    utils.RoundFloat32(source.DetectedObject.Y2),
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
