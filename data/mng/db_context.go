package mng

import (
	"log"
	"mngr/models"
)

type DbContext struct {
	Config *models.Config
	Ais    *DbSet[AiEntity]
}

func (d *DbContext) Init() error {
	cs := d.Config.Db.ConnectionString
	d.Ais = &DbSet[AiEntity]{CollectionName: "ai", ConnectionString: cs}
	err := d.Ais.Open()

	return err
}

func (d *DbContext) Close() error {
	err := d.Ais.Close()
	if err != nil {
		log.Println(err.Error())
	}

	return err
}
