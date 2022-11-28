package utils

import (
	cp "github.com/otiai10/copy"
	"io/ioutil"
	"log"
	"mngr/models"
	"os"
	"path"
	"strings"
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
	checkIfMlFolderEmpty(train)
	test := path.Join(ml, "test")
	createDirIfNotExist(test)

	// Create automatic plate license recognizer
	alpr := GetAlprPath(config)
	createDirIfNotExist(alpr)

	//create DeepStack backup directory
	ds := getDeepStackBackupPath(config)
	createDirIfNotExist(ds)
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
	createDirIfNotExist(path.Join(od, sourceId, "images"))

	// Create facial recognition folder for the source
	fr := GetFrPath(config)
	createDirIfNotExist(path.Join(fr, sourceId))
	createDirIfNotExist(path.Join(fr, sourceId, "images"))

	// Create automatic plate license recognizer
	alpr := GetAlprPath(config)
	createDirIfNotExist(path.Join(alpr, sourceId))
	createDirIfNotExist(path.Join(alpr, sourceId, "images"))
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

func GetAiClipPathBySourceId(config *models.Config, sourceId string) string {
	return path.Join(GetRecordPath(config), sourceId, "ai")
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

func GetFrTrainPath(config *models.Config) string {
	return path.Join(GetFrMlPath(config), "train")
}

func GetFrTrainPathByPerson(config *models.Config, person string) string {
	return path.Join(GetFrTrainPath(config), person)
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

// alpr ends

func SetRelativeImagePath(config *models.Config, fullImagePath string) string {
	return strings.Replace(fullImagePath, config.General.RootFolderPath+"/", "", -1)
}

func SetRelativeOdAiVideoClipPath(config *models.Config, fullVideoPath string) string {
	return strings.Replace(fullVideoPath, config.General.RootFolderPath+"/", "", -1)
}

func SetRelativeRecordPath(config *models.Config, fullRecordPath string) string {
	return strings.Replace(fullRecordPath, config.General.RootFolderPath+"/record", "playback", -1)
}

func getDeepStackBackupPath(config *models.Config) string {
	return path.Join(config.General.RootFolderPath, "deepstack")
}

func checkIfMlFolderEmpty(mlTrainDirPath string) error {
	dirs, err := ioutil.ReadDir(mlTrainDirPath)
	if err != nil {
		log.Println("an error occurred while checking the ml train directory, err: " + err.Error())
		return err
	}
	if len(dirs) == 0 {
		staticPath := "static/train"
		err = cp.Copy(staticPath, path.Join(mlTrainDirPath))
		if err != nil {
			log.Println("an error occurred while copying the ml train directory, err: " + err.Error())
		}
	}
	return err
}
