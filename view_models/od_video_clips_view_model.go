package view_models

import (
	"mngr/data"
	"path"
	"strings"
)

type PreviewViewModel struct {
	ImageFileName string `json:"image_file_name"`
	ObjectNames   string `json:"object_names"`
	Id            string `json:"id"`
}

type DetectedObjectViewModel struct {
	PredScore   float32 `json:"pred_score"`
	PredClsIdx  int     `json:"pred_cls_idx"`
	PredClsName string  `json:"pred_cls_name"`
}

type OdVideoClipsViewModel struct {
	DetectedObjects []*DetectedObjectViewModel `json:"detected_objects"`
	ImageFileNames  []string                   `json:"image_file_names"`
	Ids             []string                   `json:"ids"`

	Preview             *PreviewViewModel `json:"preview"`
	VideoCreatedAt      string            `json:"video_created_at"`
	VideoLastModifiedAt string            `json:"video_last_modified_at"`
	VideoFileName       string            `json:"video_file_name"`
	VideoBaseFileName   string            `json:"video_base_file_name"`
	Duration            int               `json:"duration"`
}

func Map(list []*data.OdDto) []*OdVideoClipsViewModel {
	ret := make([]*OdVideoClipsViewModel, 0)
	dic := make(map[string]*OdVideoClipsViewModel)
	for _, item := range list {
		if item.AiClip == nil {
			continue
		}
		key := item.AiClip.FileName
		if len(key) == 0 {
			continue
		}
		val, ok := dic[key]
		if !ok {
			val = &OdVideoClipsViewModel{}
			val.DetectedObjects = make([]*DetectedObjectViewModel, 0)
			val.ImageFileNames = make([]string, 0)
			val.Ids = make([]string, 0)

			val.VideoCreatedAt = item.AiClip.CreatedAt
			val.VideoLastModifiedAt = item.AiClip.LastModifiedAt
			val.VideoFileName = item.AiClip.FileName
			val.VideoBaseFileName = path.Base(item.AiClip.FileName)
			val.Duration = item.AiClip.Duration
			dic[key] = val
			ret = append(ret, val)
		}
		val.ImageFileNames = append(val.ImageFileNames, item.ImageFileName)
		do := &DetectedObjectViewModel{}
		do.PredClsIdx = item.DetectedObject.PredClsIdx
		do.PredClsName = item.DetectedObject.PredClsName
		do.PredScore = item.DetectedObject.PredScore
		val.DetectedObjects = append(val.DetectedObjects, do)
		val.Ids = append(val.Ids, item.Id)
	}

	for _, v := range ret {
		m := map[string]bool{}
		for _, d := range v.DetectedObjects {
			m[d.PredClsName] = true
		}
		a := make([]string, 0)
		for k, _ := range m {
			a = append(a, k)
		}
		v.Preview = &PreviewViewModel{Id: v.Ids[0], ObjectNames: strings.Join(a, ","), ImageFileName: v.ImageFileNames[0]}
	}
	return ret
}
