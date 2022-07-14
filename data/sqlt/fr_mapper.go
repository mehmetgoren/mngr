package sqlt

import (
	"mngr/data"
	"strconv"
)

type FrMapper struct {
}

func (f *FrMapper) Map(source *FrEntity) *data.FrDto {
	ret := &data.FrDto{}
	ret.Id = strconv.FormatUint(uint64(source.ID), 10)
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAtStr
	ret.DetectedFace = &data.DetectedFaceDto{
		PredScore:   source.PredScore,
		PredClsIdx:  source.PredClsIdx,
		PredClsName: source.PredClsName,
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
