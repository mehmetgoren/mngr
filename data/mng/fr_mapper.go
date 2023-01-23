package mng

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
)

type FrMapper struct {
	Config *models.Config
}

func (f *FrMapper) Map(source *FrEntity) *data.FrDto {
	ret := &data.FrDto{}
	ret.Id = source.Id.Hex()
	ret.GroupId = source.GroupId
	ret.SourceId = source.SourceId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedFace = &data.DetectedFaceDto{
		PredScore:   utils.RoundFloat32(source.DetectedFace.PredScore),
		PredClsIdx:  source.DetectedFace.PredClsIdx,
		PredClsName: source.DetectedFace.PredClsName,
		X1:          utils.RoundFloat32(source.DetectedFace.X1),
		Y1:          utils.RoundFloat32(source.DetectedFace.Y1),
		X2:          utils.RoundFloat32(source.DetectedFace.X2),
		Y2:          utils.RoundFloat32(source.DetectedFace.Y2),
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
