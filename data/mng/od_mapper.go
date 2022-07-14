package mng

import (
	"mngr/data"
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
	ret.AiClip = source.AiClip

	return ret
}
