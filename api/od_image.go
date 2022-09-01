package api

import (
	"github.com/gin-gonic/gin"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
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
	CreatedAt string `json:"-"`
}

type ImagesParams struct {
	RootPath string `json:"rootPath"`
	SourceId string `json:"sourceId"`
	AiType   int    `json:"ai_type"`
}

func RegisterOdImagesEndpoints(router *gin.Engine, rb *reps.RepoBucket, factory *cmn.Factory) {
	router.GET("/aiimagesfolders/:id/:date/:aitype", func(ctx *gin.Context) {
		sourceId := ctx.Param("id")
		date := ctx.Param("date")
		aiType, _ := strconv.Atoi(ctx.Param("aitype"))
		config, _ := rb.ConfigRep.GetConfig()
		var odPath string
		switch aiType {
		case models.Od:
			odPath = utils.GetHourlyOdImagesPathBySourceId(config, sourceId, date)
			break
		case models.Fr:
			odPath = utils.GetHourlyFrImagesPathBySourceId(config, sourceId, date)
			break
		case models.Alpr:
			odPath = utils.GetHourlyAlprImagesPathBySourceId(config, sourceId, date)
			break
		}
		items, _ := newTree(odPath, true)
		ctx.JSON(http.StatusOK, items)
	})
	// it has potential security risk
	router.POST("aiimages", func(ctx *gin.Context) {
		var model ImagesParams
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		items := make([]*ImageItem, 0)
		si := models.CreateDateSort(factory.GetCreatedDateFieldName())

		switch model.AiType {
		case models.Od:
			dtos, err := factory.CreateRepository().QueryOds(*data.GetParamsByHour(model.SourceId, model.RootPath, si))
			if err != nil || dtos == nil {
				ctx.JSON(http.StatusOK, items)
				return
			}
			for _, dto := range dtos {
				item := &ImageItem{Id: dto.Id, ImagePath: dto.ImageFileName, CreatedAt: dto.CreatedAt}
				items = append(items, item)
			}
			break
		case models.Fr:
			dtos, err := factory.CreateRepository().QueryFrs(*data.GetParamsByHour(model.SourceId, model.RootPath, si))
			if err != nil {
				ctx.JSON(http.StatusOK, items)
				return
			}
			for _, dto := range dtos {
				item := &ImageItem{Id: dto.Id, ImagePath: dto.ImageFileName, CreatedAt: dto.CreatedAt}
				items = append(items, item)
			}
			break
		case models.Alpr:
			dtos, err := factory.CreateRepository().QueryAlprs(*data.GetParamsByHour(model.SourceId, model.RootPath, si))
			if err != nil {
				ctx.JSON(http.StatusOK, items)
				return
			}
			for _, dto := range dtos {
				item := &ImageItem{Id: dto.Id, ImagePath: dto.ImageFileName, CreatedAt: dto.CreatedAt}
				items = append(items, item)
			}
			break
		}
		ctx.JSON(http.StatusOK, items)
	})
}
