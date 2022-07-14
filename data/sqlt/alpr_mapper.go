package sqlt

import (
	"mngr/data"
	"strconv"
)

type AlprMapper struct {
}

func (a *AlprMapper) Map(source *AlprEntity) *data.AlprDto {
	ret := &data.AlprDto{}
	ret.Id = strconv.FormatUint(uint64(source.ID), 10)
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAtStr
	ret.DetectedPlate = &data.DetectedPlateDto{
		Plate:      source.Plate,
		Confidence: source.Confidence,
	}
	ret.ImageFileName = source.ImageFileName
	ret.VideoFileName = source.VideoFileName
	ret.AiClip = &data.AiClip{
		Enabled:        source.AiClipEnabled,
		FileName:       source.AiClipFileName,
		CreatedAt:      source.CreatedAtStr,
		LastModifiedAt: source.AiClipLastModifiedAtStr,
		Duration:       source.AiClipDuration,
	}

	return ret
}
