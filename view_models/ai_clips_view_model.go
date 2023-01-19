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
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name"`
}

type AiClipViewModel struct {
	AiType         int                  `json:"ai_type"`
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
			m[d.PredClsName] = true
		}
		a := make([]string, 0)
		for k, _ := range m {
			a = append(a, k)
		}
		v.Preview = &AiPreviewViewModel{Id: v.Ids[0], ObjectNames: strings.Join(a, ","), ImageFileName: v.ImageFileNames[0]}
	}
}

func Map(sourceId string, dtos []data.AiDto) []*AiClipViewModel {
	ret := make([]*AiClipViewModel, 0)
	dic := make(map[string]*AiClipViewModel)
	for _, dto := range dtos {
		aiClip := dto.GetAiClip()
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
		val.AiType = dto.GetAiType()
		val.ImageFileNames = append(val.ImageFileNames, dto.GetImageFileName())
		do := &AiObjectViewModel{}
		do.PredClsIdx = dto.GetPredClsIdx()
		do.PredClsName = dto.GetPredClsName()
		do.PredScore = dto.GetPredScore()
		val.AiObjects = append(val.AiObjects, do)
		val.Ids = append(val.Ids, dto.GetId())
		val.SourceId = sourceId
	}
	setAiPreviewViewModel(ret)
	return ret
}

type AiClipQueryViewModel struct {
	AiType   int    `json:"ai_type"`
	SourceId string `json:"source_id"`
	Date     string `json:"date"`
}
