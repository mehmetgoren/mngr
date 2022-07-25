package sqlt

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
	"strconv"
)

type AlprMapper struct {
	Config *models.Config
}

func (a *AlprMapper) Map(source *AlprEntity) *data.AlprDto {
	ret := &data.AlprDto{}
	ret.Id = strconv.FormatUint(uint64(source.ID), 10)
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAtStr
	ret.DetectedPlate = &data.DetectedPlateDto{
		Plate:      source.Plate,
		Confidence: utils.RoundFloat64(source.Confidence),
	}
	ret.ImageFileName = utils.SetRelativeImagePath(a.Config, source.ImageFileName)

	ret.VideoFile = &data.VideoFileDto{}
	ret.VideoFile.Name = utils.SetRelativeRecordPath(a.Config, source.VideoFileName)
	if source.VideoFileCreatedDate != nil {
		ret.VideoFile.CreatedAt = utils.TimeToString(*source.VideoFileCreatedDate, false)
	}
	ret.VideoFile.Duration = source.VideoFileDuration
	ret.VideoFile.Merged = source.VideoFileMerged
	ret.VideoFile.ObjectAppearsAt = source.ObjectAppearsAt

	ret.AiClip = &data.AiClip{
		Enabled:        source.AiClipEnabled,
		FileName:       utils.SetRelativeOdAiVideoClipPath(a.Config, source.AiClipFileName),
		CreatedAt:      source.CreatedAtStr,
		LastModifiedAt: source.AiClipLastModifiedAtStr,
		Duration:       source.AiClipDuration,
	}

	return ret
}
