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
		ai := GetAiPath(dirPath)
		createDirIfNotExist(ai)
	}

	ml := GetFacePath(config)
	createDirIfNotExist(ml)
	train := GetFaceTrainPath(config)
	createDirIfNotExist(train)
	checkIfFaceFolderEmpty(train)
	test := path.Join(ml, "test")
	createDirIfNotExist(test)
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
		ai := GetAiPath(dirPath)
		createDirIfNotExist(path.Join(ai, sourceId))
		createDirIfNotExist(path.Join(ai, sourceId, "images"))
	}
}

func GetStreamPath(dirPath string) string {
	return path.Join(dirPath, "stream")
}

func GetRecordPath(dirPath string) string {
	return path.Join(dirPath, "record")
}

func GetAiPath(dirPath string) string {
	return path.Join(dirPath, "ai")
}

// GetFacePath returns only first path from config.General.DirPaths
func GetFacePath(config *models.Config) string {
	rootDirPath := GetDefaultDirPath(config)
	face := GetAiPath(rootDirPath)
	return path.Join(face, "face")
}

func GetFaceTrainPath(config *models.Config) string {
	return path.Join(GetFacePath(config), "train")
}

func GetFaceTrainPathByPerson(config *models.Config, person string) string {
	return path.Join(GetFaceTrainPath(config), person)
}

func GetAiClipPathBySource(config *models.Config, ffmpegModel models.FFmpegModel) string {
	sourceDirPath := GetSourceDirPath(config, ffmpegModel)
	return path.Join(GetRecordPath(sourceDirPath), ffmpegModel.GetSourceId(), "ai")
}

func GetAiImagesPathBySource(config *models.Config, ffmpegModel models.FFmpegModel) string {
	sourceDirPath := GetSourceDirPath(config, ffmpegModel)
	return path.Join(GetAiPath(sourceDirPath), ffmpegModel.GetSourceId(), "images")
}

func GetHourlyAiImagesPathBySource(config *models.Config, ffmpegModel models.FFmpegModel, dateStr string) string {
	di := DateIndex{}
	di.SetValuesFrom(dateStr)
	return di.GetIndexedPath(GetAiImagesPathBySource(config, ffmpegModel))
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

func checkIfFaceFolderEmpty(faceTrainDirPath string) error {
	dirs, err := ioutil.ReadDir(faceTrainDirPath)
	if err != nil {
		log.Println("an error occurred while checking the face train directory, err: " + err.Error())
		return err
	}
	if len(dirs) == 0 {
		staticPath := "static/train"
		err = cp.Copy(staticPath, path.Join(faceTrainDirPath))
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
