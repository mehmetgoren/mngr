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
	ret.Id = source.Id.String()
	ret.GroupId = source.GroupId
	ret.SourceId = source.SourceId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedFace = &data.DetectedFaceDto{
		PredScore:   utils.RoundFloat32(source.DetectedFace.PredScore),
		PredClsIdx:  source.DetectedFace.PredClsIdx,
		PredClsName: source.DetectedFace.PredClsName,
	}
	ret.ImageFileName = utils.SetRelativeImagePath(f.Config, source.ImageFileName)
	ret.VideoFile = &data.VideoFileDto{}
	if source.VideoFile != nil {
		ret.VideoFile.Name = utils.SetRelativeRecordPath(f.Config, source.VideoFile.Name)
		ret.VideoFile.CreatedAt = utils.TimeToString(source.VideoFile.CreatedDate.Time(), false)
		ret.VideoFile.Duration = source.VideoFile.Duration
		ret.VideoFile.Merged = source.VideoFile.Merged
		ret.VideoFile.ObjectAppearsAt = source.VideoFile.ObjectAppearsAt
	}
	ret.AiClip = source.AiClip
	ret.AiClip.FileName = utils.SetRelativeOdAiVideoClipPath(f.Config, source.AiClip.FileName)

	return ret
}
