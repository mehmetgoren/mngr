package mng

import "mngr/data"

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
	ret.VideoFileName = source.VideoFileName
	ret.AiClip = source.AiClip

	return ret
}
