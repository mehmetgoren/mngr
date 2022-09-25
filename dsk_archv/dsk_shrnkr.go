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
	"mngr/utils"
	"os"
	"path"
	"sort"
)

type DiskShrinker struct {
	Factory      *cmn.Factory
	Rb           *reps.RepoBucket
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
	confirmer := ActionTypeConfirmer{Config: d.Factory.Config}
	actionType := confirmer.GetActionType()
	opName := "deleting"
	if actionType == MoveToNewLocation {
		opName = "moving"
	}
	log.Println("disk shrink action type is " + opName)
	switch actionType {
	case Delete:
		delete_(d.Factory, d.Rb)
		break
	case MoveToNewLocation:
		move(d.Factory, d.Rb)
		break
	}
	return nil
}

func shrink(fc *cmn.Factory, rb *reps.RepoBucket, fn func(oldest *OldestSourceRecord, sourceId string) error) {
	sources, _ := rb.SourceRep.GetAll()
	if sources == nil {
		return
	}
	oldest := &OldestSourceRecord{}
	for _, source := range sources {
		fullPath := utils.GetRecordPathBySourceId(fc.Config, source.Id)
		deleteAllTempAiClipFiles(fullPath)

		oldest.Init(fullPath)
		if !oldest.Found {
			oldest.DeleteParentDirectoryIfEmpty(fullPath)
			continue
		}
		err := fn(oldest, source.Id)
		if err != nil {
			log.Println("an error occurred while shrinking the oldest record directory, err: " + err.Error())
		}
		oldest.DeleteParentDirectoryIfEmpty(fullPath)

		deleteAllAiData(oldest, fc, source.Id)
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
		if !file.IsDir() {
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

func deleteAllAiData(oldest *OldestSourceRecord, fc *cmn.Factory, sourceId string) {
	params := data.QueryParams{ClassName: "", NoPreparingVideoFile: false,
		Sort: models.SortInfo{Enabled: false}, Paging: models.PagingInfo{Enabled: false}}
	params.SourceId = sourceId
	params.T1 = oldest.CreateMinTime()
	params.T2 = oldest.CreateMaxTime()
	deleteOptions := &data.DeleteOptions{DeleteImage: false, DeleteVideo: false}
	rep := fc.CreateRepository()

	ods, _ := rep.QueryOds(params)
	if ods != nil {
		for _, od := range ods {
			deleteOptions.Id = od.Id
			rep.DeleteOds(deleteOptions)
		}
	}
	frs, _ := rep.QueryFrs(params)
	if frs != nil {
		for _, fr := range frs {
			deleteOptions.Id = fr.Id
			rep.DeleteFrs(deleteOptions)
		}
	}
	alprs, _ := rep.QueryAlprs(params)
	if alprs != nil {
		for _, alpr := range alprs {
			deleteOptions.Id = alpr.Id
			rep.DeleteAlprs(deleteOptions)
		}
	}

	// delete daily ai clips by source id
	rootDailyAiPath := oldest.CreateDailyPathName(utils.GetAiClipPathBySourceId(fc.Config, sourceId))
	err := os.RemoveAll(rootDailyAiPath)
	if err != nil {
		log.Println("an error occurred while deleting daily AI Clips root path, err: " + err.Error())
	} else {
		log.Println("AI Clips parent folder has been deleted: " + rootDailyAiPath)
	}
	// delete daily od images by source id
	rootDailyOdImagePath := oldest.CreateDailyPathName(utils.GetOdImagesPathBySourceId(fc.Config, sourceId))
	err = os.RemoveAll(rootDailyOdImagePath)
	if err != nil {
		log.Println("an error occurred while deleting daily Od Images root path, err: " + err.Error())
	} else {
		log.Println("Od Images parent folder has been deleted: " + rootDailyOdImagePath)
	}
	// delete daily fr images by source id
	rootDailyFrImagePath := oldest.CreateDailyPathName(utils.GetFrImagesPathBySourceId(fc.Config, sourceId))
	err = os.RemoveAll(rootDailyFrImagePath)
	if err != nil {
		log.Println("an error occurred while deleting daily Fr Images root path, err: " + err.Error())
	} else {
		log.Println("Fr Images parent folder has been deleted: " + rootDailyFrImagePath)
	}
	// delete daily alpr images by source id
	rootDailyAlprImagePath := oldest.CreateDailyPathName(utils.GetAlprImagesPathBySourceId(fc.Config, sourceId))
	err = os.RemoveAll(rootDailyAlprImagePath)
	if err != nil {
		log.Println("an error occurred while deleting daily Alpr Images root path, err: " + err.Error())
	} else {
		log.Println("Alpr Images parent folder has been deleted: " + rootDailyAlprImagePath)
	}
}

func delete_(fc *cmn.Factory, rb *reps.RepoBucket) {
	shrink(fc, rb, func(oldest *OldestSourceRecord, sourceId string) error {
		return os.RemoveAll(oldest.Path)
	})
}

func move(fc *cmn.Factory, rb *reps.RepoBucket) {
	shrink(fc, rb, func(oldest *OldestSourceRecord, sourceId string) error {
		moveLoc := path.Join(fc.Config.Archive.MoveLocation, sourceId, "record", oldest.CreateTmpFolderPathName())
		err := cp.Copy(oldest.Path, moveLoc)
		if err != nil {
			log.Println("an error occurred while copying the old data to new directory, err: " + err.Error())
		}
		return os.RemoveAll(oldest.Path)
	})
}
