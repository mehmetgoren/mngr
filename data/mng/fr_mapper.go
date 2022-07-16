package mng

import (
	"mngr/data"
	"mngr/utils"
)

type FrMapper struct {
}

func (f *FrMapper) Map(source *FrEntity) *data.FrDto {
	ret := &data.FrDto{}
	ret.Id = source.Id.String()
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedFace = &data.DetectedFaceDto{
		PredScore:   source.DetectedFace.PredScore,
		PredClsIdx:  source.DetectedFace.PredClsIdx,
		PredClsName: source.DetectedFace.PredClsName,
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
