package reps

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"mngr/models"
	"mngr/utils"
	"path"
	"path/filepath"
	"sort"
)

type FrHandlerRepository struct {
	Config *models.Config
}

func (f *FrHandlerRepository) GetJsonObjects(sourceId string, dateString string, sorted bool) []*models.FaceRecognitionJsonObject {
	rootSourceDataPath := utils.GetFrDataPathBySourceId(f.Config, sourceId)
	t := utils.StringToTime(dateString)
	ti := utils.TimeIndex{}
	ti.SetValuesFrom(&t)
	indexedSourceDataPath := ti.GetIndexedPath(rootSourceDataPath)

	items := make([]*models.FaceRecognitionJsonObject, 0)
	filepath.Walk(indexedSourceDataPath, func(fileInfoFullPath string, file fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if file.IsDir() {
			return err
		}

		ext := path.Ext(file.Name())
		if ext != ".json" {
			return err
		}

		bytes, _ := ioutil.ReadFile(fileInfoFullPath)
		item := &models.FaceRecognitionJsonObject{}
		err = json.Unmarshal(bytes, item)
		if err != nil {
			log.Println("an error occurred while deserializing the data json file, err: ", err)
			return err
		}

		setupTimes(item.Video)

		items = append(items, item)

		return nil
	})
	if sorted {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Video.CreatedAtTime.After(items[j].Video.CreatedAtTime)
		})
	}
	return items
}
