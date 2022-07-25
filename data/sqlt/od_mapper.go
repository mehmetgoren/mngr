package sqlt

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
	"strconv"
)

type OdMapper struct {
	Config *models.Config
}

func (o *OdMapper) Map(source *OdEntity) *data.OdDto {
	ret := &data.OdDto{}
	ret.Id = strconv.FormatUint(uint64(source.ID), 10)
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAtStr
	ret.DetectedObject = &data.DetectedObjectDto{
		PredScore:   utils.RoundFloat32(source.PredScore),
		PredClsIdx:  source.PredClsIdx,
		PredClsName: source.PredClsName,
	}
	ret.ImageFileName = utils.SetRelativeImagePath(o.Config, source.ImageFileName)

	ret.VideoFile = &data.VideoFileDto{}
	ret.VideoFile.Name = utils.SetRelativeRecordPath(o.Config, source.VideoFileName)
	if source.VideoFileCreatedDate != nil {
		ret.VideoFile.CreatedAt = utils.TimeToString(*source.VideoFileCreatedDate, false)
	}
	ret.VideoFile.Duration = source.VideoFileDuration
	ret.VideoFile.Merged = source.VideoFileMerged
	ret.VideoFile.ObjectAppearsAt = source.ObjectAppearsAt

	ret.AiClip = &data.AiClip{
		Enabled:        source.AiClipEnabled,
		FileName:       utils.SetRelativeOdAiVideoClipPath(o.Config, source.AiClipFileName),
		CreatedAt:      source.CreatedAtStr,
		LastModifiedAt: source.AiClipLastModifiedAtStr,
		Duration:       source.AiClipDuration,
	}

	return ret
}
