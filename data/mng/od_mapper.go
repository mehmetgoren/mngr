package mng

import (
	"mngr/data"
	"mngr/utils"
)

type OdMapper struct {
}

func (o *OdMapper) Map(source *OdEntity) *data.OdDto {
	ret := &data.OdDto{}
	ret.Id = source.Id.Hex()
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedObject = &data.DetectedObjectDto{
		PredScore:   source.DetectedObject.PredScore,
		PredClsIdx:  source.DetectedObject.PredClsIdx,
		PredClsName: source.DetectedObject.PredClsName,
	}
	ret.ImageFileName = source.ImageFileName
	ret.VideoFileName = source.VideoFileName
	if source.VideoFileCreatedDate != nil {
		ret.VideoFileCreatedAt = utils.TimeToString(*source.VideoFileCreatedDate, false)
	}
	ret.VideoFileDuration = source.VideoFileDuration
	ret.AiClip = source.AiClip

	return ret
}
