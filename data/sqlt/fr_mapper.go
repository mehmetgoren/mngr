package sqlt

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
	"strconv"
)

type FrMapper struct {
	Config *models.Config
}

func (f *FrMapper) Map(source *FrEntity) *data.FrDto {
	ret := &data.FrDto{}
	ret.Id = strconv.FormatUint(uint64(source.ID), 10)
	ret.GroupId = source.GroupId
	ret.SourceId = source.SourceId
	ret.CreatedAt = source.CreatedAtStr
	ret.DetectedFace = &data.DetectedFaceDto{
		PredScore:   utils.RoundFloat32(source.PredScore),
		PredClsIdx:  source.PredClsIdx,
		PredClsName: source.PredClsName,
	}
	ret.ImageFileName = utils.SetRelativeImagePath(f.Config, source.ImageFileName)

	ret.VideoFile = &data.VideoFileDto{}
	ret.VideoFile.Name = utils.SetRelativeRecordPath(f.Config, source.VideoFileName)
	if source.VideoFileCreatedDate != nil {
		ret.VideoFile.CreatedAt = utils.TimeToString(*source.VideoFileCreatedDate, false)
	}
	ret.VideoFile.Duration = source.VideoFileDuration
	ret.VideoFile.Merged = source.VideoFileMerged
	ret.VideoFile.ObjectAppearsAt = source.ObjectAppearsAt

	ret.AiClip = &data.AiClip{
		Enabled:        source.AiClipEnabled,
		FileName:       utils.SetRelativeOdAiVideoClipPath(f.Config, source.AiClipFileName),
		CreatedAt:      source.CreatedAtStr,
		LastModifiedAt: source.AiClipLastModifiedAtStr,
		Duration:       source.AiClipDuration,
	}

	return ret
}
