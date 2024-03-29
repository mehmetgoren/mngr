package dsk_archv

import (
	"github.com/go-co-op/gocron"
	"log"
	"mngr/data/cmn"
	"mngr/reps"
	"mngr/server_stats"
	"time"
)

type DiskUsageChecker struct {
	Factory   *cmn.Factory
	Rb        *reps.RepoBucket
	stats     *server_stats.ServerStats
	scheduler *gocron.Scheduler
}

func (d *DiskUsageChecker) StartScheduler() {
	if d.stats == nil {
		d.stats = &server_stats.ServerStats{}
	}
	config := d.Factory.Config
	limitPercent := config.Archive.LimitPercent
	if limitPercent <= 0 {
		limitPercent = 1
	} else if limitPercent >= 100 {
		limitPercent = 99
	}
	d.scheduler = gocron.NewScheduler(time.UTC)
	d.scheduler.Every(1).Minute().Do(func() {
		err := d.stats.InitDiskInfos(config)
		if err != nil {
			log.Println("an error occurred while getting disk usage info for DiskUsageChecker, err: " + err.Error())
			return
		}
		for _, disk := range d.stats.Disks {
			if int(disk.UsagePercent) >= limitPercent {
				ds := &DiskShrinker{Factory: d.Factory, Rb: d.Rb, DiskInfo: &disk}
				err = ds.Shrink()
				if err != nil {
					log.Println("an error occurred while shrinking the disk, err: " + err.Error())
				}
			} else {
				log.Println("Disk Usage has checked the disk and no action has been taken")
			}
		}
	})
	d.scheduler.StartAsync()
}

func (d *DiskUsageChecker) StopScheduler() {
	d.scheduler.Stop()
	d.scheduler.Clear()
}

func InitDiskUsageChecker(factory *cmn.Factory, rb *reps.RepoBucket) *DiskUsageChecker {
	ret := &DiskUsageChecker{Factory: factory, Rb: rb}
	ret.StartScheduler()
	return ret
}
