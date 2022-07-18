package mng

import (
	"mngr/data"
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
