package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
		parents[path] = &FolderTreeItem{
			FullPath:   path,
			Label:      info.Name(),
			Size:       info.Size(),
			ModifiedAt: utils.FromDateToString(info.ModTime()),
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
	FullPath    string  `json:"fullPath"`
	SourceId    string  `json:"sourceId"`
	ClassName   string  `json:"className"`
	Score       float32 `json:"score"`
	ModifiedAt  string  `json:"modifiedAt"`
	Id          string  `json:"id"`
	Base64Image string  `json:"base64Image"`
}

type DetectedImagesParams struct {
	RootPath string `json:"rootPath"`
	SourceId string `json:"sourceId"`
}

func RegisterDetectedEndpoints(router *gin.Engine) {
	router.GET("/detectedfolders", func(ctx *gin.Context) {
		config, _ := utils.ConfigRep.GetConfig()
		path := config.AiConfig.DetectedFolder
		items, _ := newTree(path, true)
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
			splits := strings.Split(file.Name(), "_")
			if len(splits) != 11 {
				continue
			}
			item := &ImageItem{FullPath: filepath.Join(model.RootPath, file.Name()), SourceId: splits[0], ClassName: splits[1]}
			score, _ := strconv.ParseFloat(splits[2], 32)
			item.Score = float32(score)
			item.ModifiedAt = strings.Join(splits[3:9], "_")
			item.Id = splits[10]

			bytes, _ := ioutil.ReadFile(item.FullPath)
			item.Base64Image = base64.StdEncoding.EncodeToString(bytes)

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
