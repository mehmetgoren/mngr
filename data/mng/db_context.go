package mng

import (
	"log"
	"mngr/models"
)

type DbContext struct {
	Config *models.Config

	Ods   *DbSet[OdEntity]
	Frs   *DbSet[FrEntity]
	Alprs *DbSet[AlprEntity]
}

func (d *DbContext) Init() error {
	cs := d.Config.Db.ConnectionString
	d.Ods = &DbSet[OdEntity]{CollectionName: "od", ConnectionString: cs}
	err := d.Ods.Open()
	if err != nil {
		return err
	}

	if err == nil {
		d.Frs = &DbSet[FrEntity]{CollectionName: "fr", ConnectionString: cs}
		err = d.Frs.Open()
		if err != nil {
			return err
		}

		d.Alprs = &DbSet[AlprEntity]{CollectionName: "alpr", ConnectionString: cs}
		err = d.Alprs.Open()
	}

	return err
}

func (d *DbContext) Close() error {
	err := d.Ods.Close()
	if err != nil {
		log.Println(err.Error())
	}
	err = d.Frs.Close()
	if err != nil {
		log.Println(err.Error())
	}
	err = d.Alprs.Close()
	if err != nil {
		log.Println(err.Error())
	}

	return err
}
