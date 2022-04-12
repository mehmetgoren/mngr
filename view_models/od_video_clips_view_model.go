package view_models

import (
	"mngr/models"
	"path"
	"strings"
)

type PreviewViewModel struct {
	ImageFileName string `json:"image_file_name"`
	ObjectNames   string `json:"object_names"`
	Id            string `json:"id"`
}

type OdVideoClipsViewModel struct {
	DetectedObjects []*models.DetectedObject `json:"detected_objects"`
	ImageFileNames  []string                 `json:"image_file_names"`
	DataFileNames   []string                 `json:"data_file_names"`
	Ids             []string                 `json:"ids"`

	Preview             *PreviewViewModel `json:"preview"`
	VideoCreatedAt      string            `json:"video_created_at"`
	VideoLastModifiedAt string            `json:"video_last_modified_at"`
	VideoFileName       string            `json:"video_file_name"`
	VideoBaseFileName   string            `json:"video_base_file_name"`
	Duration            int               `json:"duration"`
}

func Map(list []*models.ObjectDetectionJsonObject) []*OdVideoClipsViewModel {
	ret := make([]*OdVideoClipsViewModel, 0)
	dic := make(map[string]*OdVideoClipsViewModel)
	for _, item := range list {
		key := item.Video.FileName
		if len(key) == 0 {
			continue
		}
		val, ok := dic[key]
		if !ok {
			val = &OdVideoClipsViewModel{}
			val.DetectedObjects = make([]*models.DetectedObject, 0)
			val.ImageFileNames = make([]string, 0)
			val.DataFileNames = make([]string, 0)
			val.Ids = make([]string, 0)

			val.VideoCreatedAt = item.Video.CreatedAt
			val.VideoLastModifiedAt = item.Video.LastModifiedAt
			val.VideoFileName = item.Video.FileName
			val.VideoBaseFileName = path.Base(item.Video.FileName)
			val.Duration = item.Video.Duration
			dic[key] = val
			ret = append(ret, val)
		}
		val.ImageFileNames = append(val.ImageFileNames, item.ObjectDetection.ImageFileName)
		val.DataFileNames = append(val.DataFileNames, item.ObjectDetection.DataFileName)
		val.DetectedObjects = append(val.DetectedObjects, item.ObjectDetection.DetectedObjects...)
		val.Ids = append(val.Ids, item.ObjectDetection.Id)
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
