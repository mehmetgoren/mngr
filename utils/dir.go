package utils

import (
	cp "github.com/otiai10/copy"
	"io/ioutil"
	"log"
	"mngr/models"
	"os"
	"path"
	"regexp"
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

// CreateRequiredDirectories we are ensure that dirPaths exists
func CreateRequiredDirectories(config *models.Config) {
	for _, dirPath := range config.General.DirPaths {
		// Create HLS stream folder
		stream := GetStreamPath(dirPath)
		createDirIfNotExist(stream)

		// Create record folder
		record := GetRecordPath(dirPath)
		createDirIfNotExist(record)

		// Create object detection folder
		od := GetOdPath(dirPath)
		createDirIfNotExist(od)

		// Create facial recognition folder
		fr := GetFrPath(dirPath)
		createDirIfNotExist(fr)

		// Create automatic plate license recognizer
		alpr := GetAlprPath(dirPath)
		createDirIfNotExist(alpr)
	}

	ml := GetFrMlPath(config)
	createDirIfNotExist(ml)
	train := GetFrTrainPath(config)
	createDirIfNotExist(train)
	checkIfMlFolderEmpty(train)
	test := path.Join(ml, "test")
	createDirIfNotExist(test)

	//create DeepStack backup directory
	ds := getDeepStackBackupPath(config)
	createDirIfNotExist(ds)
}

func CreateSourceDefaultDirectories(config *models.Config, sourceId string) {
	for _, dirPath := range config.General.DirPaths {
		// Create HLS stream folder for the source
		stream := GetStreamPath(dirPath)
		createDirIfNotExist(path.Join(stream, sourceId))

		// Create record folder for the source
		record := GetRecordPath(dirPath)
		createDirIfNotExist(path.Join(record, sourceId))
		//and also short video clips folder
		createDirIfNotExist(path.Join(record, sourceId, "ai"))

		// Create object detection folder for the source
		od := GetOdPath(dirPath)
		createDirIfNotExist(path.Join(od, sourceId))
		createDirIfNotExist(path.Join(od, sourceId, "images"))

		// Create facial recognition folder for the source
		fr := GetFrPath(dirPath)
		createDirIfNotExist(path.Join(fr, sourceId))
		createDirIfNotExist(path.Join(fr, sourceId, "images"))

		// Create automatic plate license recognizer
		alpr := GetAlprPath(dirPath)
		createDirIfNotExist(path.Join(alpr, sourceId))
		createDirIfNotExist(path.Join(alpr, sourceId, "images"))
	}
}

func GetStreamPath(dirPath string) string {
	return path.Join(dirPath, "stream")
}

func GetRecordPath(dirPath string) string {
	return path.Join(dirPath, "record")
}

func GetOdPath(dirPath string) string {
	return path.Join(dirPath, "od")
}

func GetFrPath(dirPath string) string {
	return path.Join(dirPath, "fr")
}

func GetAlprPath(dirPath string) string {
	return path.Join(dirPath, "alpr")
}

// GetFrMlPath returns only first path from config.General.DirPaths
func GetFrMlPath(config *models.Config) string {
	rootDirPath := GetDefaultDirPath(config)
	fr := GetFrPath(rootDirPath)
	return path.Join(fr, "ml")
}

func GetAiClipPathBySource(config *models.Config, ffmpegModel models.FFmpegModel) string {
	sourceDirPath := GetSourceDirPath(config, ffmpegModel)
	return path.Join(GetRecordPath(sourceDirPath), ffmpegModel.GetSourceId(), "ai")
}

func GetOdImagesPathBySource(config *models.Config, ffmpegModel models.FFmpegModel) string {
	sourceDirPath := GetSourceDirPath(config, ffmpegModel)
	return path.Join(GetOdPath(sourceDirPath), ffmpegModel.GetSourceId(), "images")
}

func GetHourlyOdImagesPathBySource(config *models.Config, ffmpegModel models.FFmpegModel, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetOdImagesPathBySource(config, ffmpegModel))
}

func GetFrImagesPathBySource(config *models.Config, ffmpegModel models.FFmpegModel) string {
	sourceDirPath := GetSourceDirPath(config, ffmpegModel)
	return path.Join(GetFrPath(sourceDirPath), ffmpegModel.GetSourceId(), "images")
}
func GetHourlyFrImagesPathBySource(config *models.Config, ffmpegModel models.FFmpegModel, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetFrImagesPathBySource(config, ffmpegModel))
}

func GetFrTrainPath(config *models.Config) string {
	return path.Join(GetFrMlPath(config), "train")
}

func GetFrTrainPathByPerson(config *models.Config, person string) string {
	return path.Join(GetFrTrainPath(config), person)
}

func GetRecordPathBySource(config *models.Config, ffmpegModel models.FFmpegModel) string {
	sourceDirPath := GetSourceDirPath(config, ffmpegModel)
	return path.Join(GetRecordPath(sourceDirPath), ffmpegModel.GetSourceId())
}

func GetHourlyRecordPathBySource(config *models.Config, ffmpegModel models.FFmpegModel, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetRecordPathBySource(config, ffmpegModel))
}

func GetAlprImagesPathBySource(config *models.Config, ffmpegModel models.FFmpegModel) string {
	sourceDirPath := GetSourceDirPath(config, ffmpegModel)
	return path.Join(GetAlprPath(sourceDirPath), ffmpegModel.GetSourceId(), "images")
}
func GetHourlyAlprImagesPathBySource(config *models.Config, ffmpegModel models.FFmpegModel, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetAlprImagesPathBySource(config, ffmpegModel))
}

func getDeepStackBackupPath(config *models.Config) string {
	rootDirPath := GetDefaultDirPath(config)
	return path.Join(rootDirPath, "deepstack")
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

func IsDirExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetVideoFileAbsolutePath(v *models.VideoFile, config *models.Config, source models.FFmpegModel) string {
	sourceDirPath := GetRecordPath(GetSourceDirPath(config, source))
	return path.Join(sourceDirPath, v.SourceId, v.Year, FixZeroStr(v.Month), FixZeroStr(v.Day), FixZeroStr(v.Hour), v.Name)
}

func IsDirNameValid(fileName string) bool {
	regExpString := "\\/?%*:|\"<>"
	reg, err := regexp.Compile(regExpString)
	if err != nil {
		return false
	}
	return !reg.MatchString(fileName)
}

func GetDirPaths(config *models.Config) []string {
	return config.General.DirPaths
}

// GetDefaultDirPath The first item of the dirPaths is default and the Deepstack backup file also use it
func GetDefaultDirPath(config *models.Config) string {
	dirPaths := config.General.DirPaths
	if dirPaths == nil || len(dirPaths) == 0 || dirPaths[0] == "" {
		log.Fatal("Config.General.DirPaths is empty, the program will be terminated")
	}
	rootPath := dirPaths[0]
	return rootPath
}

func GetSourceDirPath(config *models.Config, ffmpegModel models.FFmpegModel) string {
	sourceDirPath := ffmpegModel.GetDirPath()
	if len(sourceDirPath) > 0 {
		return sourceDirPath
	}
	return GetDefaultDirPath(config)
}
