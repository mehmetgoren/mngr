package mng

import (
	"mngr/data"
	"mngr/models"
	"mngr/utils"
)

type OdMapper struct {
	Config *models.Config
}

func (o *OdMapper) Map(source *OdEntity) *data.OdDto {
	ret := &data.OdDto{}
	ret.Id = source.Id.Hex()
	ret.GroupId = source.GroupId
	ret.SourceId = source.SourceId
	ret.CreatedAt = source.CreatedAt
	ret.DetectedObject = &data.DetectedObjectDto{
		PredScore:   utils.RoundFloat32(source.DetectedObject.PredScore),
		PredClsIdx:  source.DetectedObject.PredClsIdx,
		PredClsName: source.DetectedObject.PredClsName,
		X1:          utils.RoundFloat32(source.DetectedObject.X1),
		Y1:          utils.RoundFloat32(source.DetectedObject.Y1),
		X2:          utils.RoundFloat32(source.DetectedObject.X2),
		Y2:          utils.RoundFloat32(source.DetectedObject.Y2),
	}
	if source.DetectedObject.Metadata != nil {
		ret.DetectedObject.Metadata = &data.MetadataDto{}
		if source.DetectedObject.Metadata.Colors != nil {
			ret.DetectedObject.Metadata.Colors = make([]data.ColorDto, 0)
			for _, color := range source.DetectedObject.Metadata.Colors {
				ret.DetectedObject.Metadata.Colors = append(ret.DetectedObject.Metadata.Colors, data.ColorDto{
					R: color.R,
					G: color.G,
					B: color.B,
				})
			}
		}
	}
	ret.ImageFileName = source.ImageFileName
	ret.VideoFile = &data.VideoFileDto{}
	if source.VideoFile != nil {
		ret.VideoFile.Name = source.VideoFile.Name
		ret.VideoFile.CreatedAt = utils.TimeToString(source.VideoFile.CreatedDate.Time(), false)
		ret.VideoFile.Duration = source.VideoFile.Duration
		ret.VideoFile.Merged = source.VideoFile.Merged
		ret.VideoFile.ObjectAppearsAt = source.VideoFile.ObjectAppearsAt
	}
	ret.AiClip = source.AiClip
	ret.AiClip.FileName = source.AiClip.FileName

	return ret
}
