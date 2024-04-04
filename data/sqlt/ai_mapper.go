package sqlt

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
	"strconv"
)

type AiMapper struct {
	Config *models.Config
}

func (o *AiMapper) Map(source *AiEntity) *data.AiDto {
	ret := &data.AiDto{}
	ret.Id = strconv.FormatUint(uint64(source.ID), 10)
	ret.Module = source.Module
	ret.GroupId = source.GroupId
	ret.SourceId = source.SourceId
	ret.CreatedAt = source.CreatedAtStr
	ret.DetectedObject = &data.DetectedObjectDto{
		Score: utils.RoundFloat32(source.Score),
		Label: source.Label,
	}
	ret.ImageFileName = source.ImageFileName

	ret.VideoFile = &data.VideoFileDto{}
	ret.VideoFile.Name = source.VideoFileName
	if source.VideoFileCreatedDate != nil {
		ret.VideoFile.CreatedAt = utils.TimeToString(*source.VideoFileCreatedDate, false)
	}
	ret.VideoFile.Duration = source.VideoFileDuration
	ret.VideoFile.Merged = source.VideoFileMerged
	ret.VideoFile.ObjectAppearsAt = source.ObjectAppearsAt

	ret.AiClip = &data.AiClip{
		Enabled:        source.AiClipEnabled,
		FileName:       source.AiClipFileName,
		CreatedAt:      source.CreatedAtStr,
		LastModifiedAt: source.AiClipLastModifiedAtStr,
		Duration:       source.AiClipDuration,
	}

	return ret
}
