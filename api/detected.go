package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type FolderTreeItem struct {
	FullPath   string            `json:"fullPath"`
	Label      string            `json:"label"`
	Size       int64             `json:"size"`
	ModifiedAt string            `json:"modifiedAt"`
	Children   []*FolderTreeItem `json:"children"`
	Parent     *FolderTreeItem   `json:"-"`
	Icon       string            `json:"icon"`
}

// Create directory hierarchy.
func newTree(root string, onlyFolder bool) (result *FolderTreeItem, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*FolderTreeItem)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if onlyFolder && !info.IsDir() {
			return err
		}
		parents[path] = &FolderTreeItem{
			FullPath:   path,
			Label:      info.Name(),
			Size:       info.Size(),
			ModifiedAt: utils.TimeToString(info.ModTime(), true),
			Children:   make([]*FolderTreeItem, 0),
			Icon:       "folder",
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
			sort.Slice(parent.Children, func(i, j int) bool {
				item1, _ := strconv.Atoi(parent.Children[i].Label)
				item2, _ := strconv.Atoi(parent.Children[j].Label)
				return item1 < item2
			})
		}
	}
	return
}

type ImageItem struct {
	FullPath   string  `json:"fullPath"`
	SourceId   string  `json:"sourceId"`
	ClassIndex int     `json:"classIndex"`
	ClassName  string  `json:"className"`
	Score      float32 `json:"score"`
	ModifiedAt string  `json:"modifiedAt"`
	Id         string  `json:"id"`
	ImagePath  string  `json:"imagePath"`
}

type DetectedImagesParams struct {
	RootPath string `json:"rootPath"`
	SourceId string `json:"sourceId"`
}

func RegisterDetectedEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/detectedfolders", func(ctx *gin.Context) {
		config, _ := rb.ConfigRep.GetConfig()
		odPath, _ := utils.GetOdFolder(config)
		items, _ := newTree(odPath, true)
		ctx.JSON(http.StatusOK, items)
	})
	// it has potential security risk
	router.POST("detectedimages", func(ctx *gin.Context) {
		var model DetectedImagesParams
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		files, _ := ioutil.ReadDir(model.RootPath)
		items := make([]*ImageItem, 0)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			splits := strings.Split(fileName, "_")
			if len(splits) != 9 {
				continue
			}
			sourceId := splits[0]
			if sourceId != model.SourceId {
				continue
			}
			item := &ImageItem{FullPath: filepath.Join(model.RootPath, fileName), SourceId: sourceId}
			item.ClassIndex = -1
			item.ClassName = ""
			item.Score = .0
			item.ModifiedAt = strings.Join(splits[1:7], "_")
			item.Id = splits[8]

			bytes, _ := ioutil.ReadFile(item.FullPath)
			serverRoot := utils.GetDetectedFolderPath()
			imgPath := path.Join(serverRoot, fileName)
			ioutil.WriteFile(imgPath, bytes, 0777)
			item.ImagePath = path.Join(utils.GetDetectedFolderName(), fileName)

			items = append(items, item)

			sort.Slice(items, func(i, j int) bool {
				item1 := items[i].ModifiedAt
				item2 := items[j].ModifiedAt
				return item1 > item2
			})
		}
		ctx.JSON(http.StatusOK, items)
	})
}
