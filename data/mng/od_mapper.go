package mng

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
)

type OdMapper struct {
	Config *models.Config
}

func (o *OdMapper) Map(source *OdEntity) *data.OdDto {
	ret := &data.OdDto{}
	ret.Id = source.Id.Hex()
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedObject = &data.DetectedObjectDto{
		PredScore:   utils.RoundFloat32(source.DetectedObject.PredScore),
		PredClsIdx:  source.DetectedObject.PredClsIdx,
		PredClsName: source.DetectedObject.PredClsName,
	}
	ret.ImageFileName = utils.SetRelativeImagePath(o.Config, source.ImageFileName)
	ret.VideoFile = &data.VideoFileDto{}
	if source.VideoFile != nil {
		ret.VideoFile.Name = utils.SetRelativeRecordPath(o.Config, source.VideoFile.Name)
		ret.VideoFile.CreatedAt = utils.TimeToString(source.VideoFile.CreatedDate.Time(), false)
		ret.VideoFile.Duration = source.VideoFile.Duration
		ret.VideoFile.Merged = source.VideoFile.Merged
		ret.VideoFile.ObjectAppearsAt = source.VideoFile.ObjectAppearsAt
	}
	ret.AiClip = source.AiClip
	ret.AiClip.FileName = utils.SetRelativeOdAiVideoClipPath(o.Config, source.AiClip.FileName)

	return ret
}
