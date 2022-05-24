package utils

import (
	"mngr/models"
	"os"
	"path"
)

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateRequiredDirectories(config *models.Config) {
	// Create HLS stream folder
	stream := GetStreamPath(config)
	createDirIfNotExist(stream)

	// Create record folder
	record := GetRecordPath(config)
	createDirIfNotExist(record)

	// Create object detection folder
	od := GetOdPath(config)
	createDirIfNotExist(od)

	// Create facial recognition folder
	fr := GetFrPath(config)
	createDirIfNotExist(fr)
	ml := GetFrMlPath(config)
	createDirIfNotExist(ml)
	train := GetFrTrainPath(config)
	createDirIfNotExist(train)
	test := path.Join(ml, "test")
	createDirIfNotExist(test)

	// Create automatic plate license recognizer
	alpr := GetAlprPath(config)
	createDirIfNotExist(alpr)
}

func CreateSourceDefaultDirectories(config *models.Config, sourceId string) {
	// Create HLS stream folder for the source
	stream := GetStreamPath(config)
	createDirIfNotExist(path.Join(stream, sourceId))

	// Create record folder for the source
	record := GetRecordPath(config)
	createDirIfNotExist(path.Join(record, sourceId))
	//and also short video clips folder
	createDirIfNotExist(path.Join(record, sourceId, "ai"))

	// Create object detection folder for the source
	od := GetOdPath(config)
	createDirIfNotExist(path.Join(od, sourceId))
	createDirIfNotExist(path.Join(od, sourceId, "data"))
	createDirIfNotExist(path.Join(od, sourceId, "images"))
	createDirIfNotExist(path.Join(od, sourceId, "videos"))

	// Create facial recognition folder for the source
	fr := GetFrPath(config)
	createDirIfNotExist(path.Join(fr, sourceId))
	createDirIfNotExist(path.Join(fr, sourceId, "data"))
	createDirIfNotExist(path.Join(fr, sourceId, "images"))
	createDirIfNotExist(path.Join(fr, sourceId, "videos"))

	// Create automatic plate license recognizer
	alpr := GetAlprPath(config)
	createDirIfNotExist(path.Join(alpr, sourceId))
	createDirIfNotExist(path.Join(alpr, sourceId, "data"))
	createDirIfNotExist(path.Join(alpr, sourceId, "images"))
	createDirIfNotExist(path.Join(alpr, sourceId, "videos"))
}

func GetStreamPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "stream")
}

func GetRecordPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "record")
}

func GetOdPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "od")
}

func GetFrPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "fr")
}

func GetFrMlPath(config *models.Config) string {
	fr := GetFrPath(config)
	return path.Join(fr, "ml")
}

// od starts
func GetOdImagesPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(GetOdPath(config), sourceId, "images")
}

func GetHourlyOdImagesPathBySourceId(config *models.Config, sourceId string, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetOdImagesPathBySourceId(config, sourceId))
}

func GetOdDataPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(GetOdPath(config), sourceId, "data")
}

// od ends

// fr starts
func getFrPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "fr")
}

func GetFrImagesPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(getFrPath(config), sourceId, "images")
}
func GetHourlyFrImagesPathBySourceId(config *models.Config, sourceId string, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetFrImagesPathBySourceId(config, sourceId))
}

func GetFrDataPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(getFrPath(config), sourceId, "data")
}

func GetFrTrainPath(config *models.Config) string {
	return path.Join(GetFrMlPath(config), "train")
}

//fr ends

// record starts
func GetRecordPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(GetRecordPath(config), sourceId)
}

func GetHourlyRecordPathBySourceId(config *models.Config, sourceId string, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetRecordPathBySourceId(config, sourceId))
}

// record ends

// alpr starts
func GetAlprPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "alpr")
}

func getAlprPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "alpr")
}

func GetAlprImagesPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(getAlprPath(config), sourceId, "images")
}
func GetHourlyAlprImagesPathBySourceId(config *models.Config, sourceId string, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetAlprImagesPathBySourceId(config, sourceId))
}

func GetAlprDataPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(getAlprPath(config), sourceId, "data")
}

// alpr ends
