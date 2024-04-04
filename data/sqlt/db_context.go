package sqlt

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"mngr/models"
	"path"
)

type DbContext struct {
	Config *models.Config

	Ais *DbSet[AiEntity]
}

func (d *DbContext) Init() error {
	db, _ := gorm.Open(sqlite.Open(path.Join(d.Config.Db.ConnectionString, "feniks.db")), &gorm.Config{})
	d.Ais = &DbSet[AiEntity]{db}

	return nil
}
