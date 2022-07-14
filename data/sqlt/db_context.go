package sqlt

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"mngr/models"
	"path"
)

type DbContext struct {
	Config *models.Config

	Ods   *DbSet[OdEntity]
	Frs   *DbSet[FrEntity]
	Alprs *DbSet[AlprEntity]
}

func (d *DbContext) Init() error {
	db, _ := gorm.Open(sqlite.Open(path.Join(d.Config.Db.ConnectionString, "feniks.db")), &gorm.Config{})
	d.Ods = &DbSet[OdEntity]{db}
	d.Frs = &DbSet[FrEntity]{db}
	d.Alprs = &DbSet[AlprEntity]{db}

	return nil
}
