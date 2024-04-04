package view_models

import (
	"mngr/data"
	"path"
	"strings"
)

type AiPreviewViewModel struct {
	ImageFileName string `json:"image_file_name"`
	ObjectNames   string `json:"object_names"`
	Id            string `json:"id"`
}

type AiObjectViewModel struct {
	Score float32 `json:"score"`
	Label string  `json:"label"`
}

type AiClipViewModel struct {
	Module         string               `json:"module"`
	AiObjects      []*AiObjectViewModel `json:"ai_objects"`
	ImageFileNames []string             `json:"image_file_names"`
	Ids            []string             `json:"ids"`

	Preview             *AiPreviewViewModel `json:"preview"`
	VideoCreatedAt      string              `json:"video_created_at"`
	VideoLastModifiedAt string              `json:"video_last_modified_at"`
	VideoFileName       string              `json:"video_file_name"`
	VideoBaseFileName   string              `json:"video_base_file_name"`
	Duration            int                 `json:"duration"`

	SourceId string `json:"source_id"`
}

func createAiClipViewModelFrom(source *data.AiClip) *AiClipViewModel {
	ret := &AiClipViewModel{}
	ret.AiObjects = make([]*AiObjectViewModel, 0)
	ret.ImageFileNames = make([]string, 0)
	ret.Ids = make([]string, 0)
	ret.VideoCreatedAt = source.CreatedAt
	ret.VideoLastModifiedAt = source.LastModifiedAt
	ret.VideoFileName = source.FileName
	ret.VideoBaseFileName = path.Base(source.FileName)
	ret.Duration = source.Duration
	return ret
}

func setAiPreviewViewModel(ret []*AiClipViewModel) {
	for _, v := range ret {
		m := map[string]bool{}
		for _, d := range v.AiObjects {
			m[d.Label] = true
		}
		a := make([]string, 0)
		for k, _ := range m {
			a = append(a, k)
		}
		v.Preview = &AiPreviewViewModel{Id: v.Ids[0], ObjectNames: strings.Join(a, ","), ImageFileName: v.ImageFileNames[0]}
	}
}

func Map(sourceId string, dtos []*data.AiDto) []*AiClipViewModel {
	ret := make([]*AiClipViewModel, 0)
	dic := make(map[string]*AiClipViewModel)
	for _, dto := range dtos {
		aiClip := dto.AiClip
		if aiClip == nil {
			continue
		}
		key := aiClip.FileName
		if len(key) == 0 {
			continue
		}
		val, ok := dic[key]
		if !ok {
			val = createAiClipViewModelFrom(aiClip)
			dic[key] = val
			ret = append(ret, val)
		}
		val.Module = dto.Module
		val.ImageFileNames = append(val.ImageFileNames, dto.ImageFileName)
		do := &AiObjectViewModel{}
		do.Label = dto.DetectedObject.Label
		do.Score = dto.DetectedObject.Score
		val.AiObjects = append(val.AiObjects, do)
		val.Ids = append(val.Ids, dto.Id)
		val.SourceId = sourceId
	}
	setAiPreviewViewModel(ret)
	return ret
}

type AiClipQueryViewModel struct {
	Module   string `json:"module"`
	SourceId string `json:"source_id"`
	Date     string `json:"date"`
}
