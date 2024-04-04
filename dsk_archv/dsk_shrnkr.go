package dsk_archv

import (
	cp "github.com/otiai10/copy"
	"io/fs"
	"io/ioutil"
	"log"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
	"mngr/reps"
	"mngr/server_stats"
	"mngr/utils"
	"os"
	"path"
	"sort"
	"strings"
)

type DiskShrinker struct {
	Factory      *cmn.Factory
	Rb           *reps.RepoBucket
	DiskInfo     *server_stats.DiskInfo
	stillWorking bool
}

func (d *DiskShrinker) Shrink() error {
	if d.stillWorking {
		log.Println("Disk Shrinker is still working on, passing this round")
		return nil
	}
	d.stillWorking = true
	defer func() {
		d.stillWorking = false
	}()
	confirmer := ActionTypeConfirmer{Config: d.Factory.Config, DiskInfo: d.DiskInfo}
	actionType := confirmer.GetActionType()
	opName := "deleting"
	if actionType == MoveToNewLocation {
		opName = "moving"
	}
	log.Println("disk shrink action type is " + opName)
	switch actionType {
	case Delete:
		d.delete_()
		break
	case MoveToNewLocation:
		d.move()
		break
	}
	return nil
}

func (d *DiskShrinker) shrink(fn func(oldest *OldestSourceRecord, source *models.SourceModel) error) {
	fc, rb := d.Factory, d.Rb
	mountPoint := d.DiskInfo.MountPoint
	sources, _ := rb.SourceRep.GetAll()
	if sources == nil {
		return
	}
	oldest := &OldestSourceRecord{}
	for _, source := range sources {
		// if source does not record on this disk, it is pointless to delete AI data.
		sourceDirPath := utils.GetSourceDirPath(fc.Config, source)
		if !strings.HasPrefix(sourceDirPath, mountPoint) {
			continue
		}
		fullPath := utils.GetRecordPathBySource(fc.Config, source)
		deleteAllTempAiClipFiles(fullPath)

		oldest.Init(fullPath)
		if !oldest.Found {
			oldest.DeleteParentDirectoryIfEmpty(fullPath)
			continue
		}
		err := fn(oldest, source)
		if err != nil {
			log.Println("an error occurred while shrinking the oldest record directory, err: " + err.Error())
		}
		oldest.DeleteParentDirectoryIfEmpty(fullPath)

		deleteAllAiData(oldest, fc, source)
	}
}

func deleteAllTempAiClipFiles(fullPath string) int {
	aiFullPath := path.Join(fullPath, "ai")
	ret := 0
	files, _ := ioutil.ReadDir(aiFullPath)
	if files == nil {
		return ret
	}
	mp4Files := make([]fs.FileInfo, 0)
	for _, file := range files {
		if !file.IsDir() { //means only video file, not indexed directories.
			mp4Files = append(mp4Files, file)
		}
	}
	sort.Slice(mp4Files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})
	for j := 0; j < len(mp4Files)-1; j++ {
		os.Remove(path.Join(aiFullPath, mp4Files[j].Name()))
		ret++
	}
	return ret
}

func deleteAllAiData(oldest *OldestSourceRecord, fc *cmn.Factory, source *models.SourceModel) error {
	params := data.QueryParams{Label: "", NoPreparingVideoFile: false,
		Sort: models.SortInfo{Enabled: false}, Paging: models.PagingInfo{Enabled: false}}
	params.SourceId = source.Id
	params.T1 = oldest.CreateMinTime()
	params.T2 = oldest.CreateMaxTime()
	deleteOptions := &data.DeleteOptions{DeleteImage: false, DeleteVideo: false}
	rep := fc.CreateRepository()

	ais, err := rep.QueryAis(params)
	if ais != nil {
		for _, ai := range ais {
			deleteOptions.Id = ai.Id
			err = rep.DeleteAis(deleteOptions)
		}
	}

	// delete daily ai clips by source id
	rootDailyAiPath := oldest.CreateDailyPathName(utils.GetAiClipPathBySource(fc.Config, source))
	err = os.RemoveAll(rootDailyAiPath)
	if err != nil {
		log.Println("an error occurred while deleting daily AI Clips root path, err: " + err.Error())
	} else {
		log.Println("AI Clips parent folder has been deleted: " + rootDailyAiPath)
	}
	// delete daily AI images by source id
	rootDailyAiImagePath := oldest.CreateDailyPathName(utils.GetAiImagesPathBySource(fc.Config, source))
	err = os.RemoveAll(rootDailyAiImagePath)
	if err != nil {
		log.Println("an error occurred while deleting daily AI Images root path, err: " + err.Error())
	} else {
		log.Println("AI Images parent folder has been deleted: " + rootDailyAiImagePath)
	}

	return err
}

func (d *DiskShrinker) delete_() {
	d.shrink(func(oldest *OldestSourceRecord, source *models.SourceModel) error {
		return os.RemoveAll(oldest.Path)
	})
}

// todo: It requires two more options: 1. create multiple multiple locations 2. update all AI data to new location.
func (d *DiskShrinker) move() {
	d.shrink(func(oldest *OldestSourceRecord, source *models.SourceModel) error {
		moveLoc := path.Join(d.Factory.Config.Archive.MoveLocation, source.GetSourceId(), "record", oldest.CreateTmpFolderPathName())
		err := cp.Copy(oldest.Path, moveLoc)
		if err != nil {
			log.Println("an error occurred while copying the old data to new directory, err: " + err.Error())
		}
		return os.RemoveAll(oldest.Path)
	})
}
