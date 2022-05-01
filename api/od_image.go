package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"os"
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

		splits := strings.Split(path, "/")
		l := len(splits)
		sep := "_"
		newPath := splits[l-4] + sep + splits[l-3] + sep + splits[l-2] + sep + splits[l-1]
		parents[path] = &FolderTreeItem{
			FullPath:   newPath,
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
				return item1 > item2
			})
		}
	}
	return
}

type ImageItem struct {
	Id        string `json:"id"`
	ImagePath string `json:"imagePath"`
}

type ImagesParams struct {
	RootPath string `json:"rootPath"`
	SourceId string `json:"sourceId"`
}

func RegisterOdImagesEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/odimagesfolders/:id", func(ctx *gin.Context) {
		sourceId := ctx.Param("id")
		config, _ := rb.ConfigRep.GetConfig()
		odPath := utils.GetOdImagesPathBySourceId(config, sourceId)
		items, _ := newTree(odPath, true)
		ctx.JSON(http.StatusOK, items)
	})
	// it has potential security risk
	router.POST("odimages", func(ctx *gin.Context) {
		var model ImagesParams
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		config, _ := rb.ConfigRep.GetConfig()
		odhRep := reps.OdHandlerRepository{Config: config}
		jsonObjects := odhRep.GetJsonObjects(model.SourceId, model.RootPath, true)
		items := make([]*ImageItem, 0)
		for _, jsonObject := range jsonObjects {
			od := jsonObject.ObjectDetection
			item := &ImageItem{Id: od.Id, ImagePath: od.ImageFileName}
			items = append(items, item)
		}
		ctx.JSON(http.StatusOK, items)
	})
}
