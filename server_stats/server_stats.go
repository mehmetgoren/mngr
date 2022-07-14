package server_stats

import (
	"errors"
	human "github.com/dustin/go-humanize"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
	"github.com/shirou/gopsutil/disk"
	"mngr/models"
	"mngr/utils"
	"strings"
	"time"
)

type ServerStats struct {
	Cpu     CpuInfo     `json:"cpu"`
	Memory  MemoryInfo  `json:"memory"`
	Disk    DiskInfo    `json:"disk"`
	Network NetworkInfo `json:"network"`
}

func (s *ServerStats) InitCpuInfos() error {
	c := CpuInfo{}
	before, err := cpu.Get()
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		return err
	}
	total := float64(after.Total - before.Total)
	c.Count = before.CPUCount
	c.UserUsage = utils.Round(float64(after.User-before.User) / total * 100)
	c.UserUsageHuman = human.CommafWithDigits(c.UserUsage, 2) + " %"
	c.SystemUsage = utils.Round(float64(after.System-before.System) / total * 100)
	c.SystemUsageHuman = human.CommafWithDigits(c.SystemUsage, 2) + " %"
	c.Idle = utils.Round(float64(after.Idle-before.Idle) / total * 100)
	c.IdleHuman = human.CommafWithDigits(c.Idle, 2) + " %"
	c.UsagePercent = utils.Round(100.0 - c.Idle)
	c.UsagePercentHuman = human.CommafWithDigits(c.UsagePercent, 2) + " %"

	s.Cpu = c
	return nil
}

func (s *ServerStats) InitMemInfos() error {
	mem, err := memory.Get()
	if err != nil {
		return err
	}

	m := MemoryInfo{}
	m.Total = mem.Total / 1024 / 1024
	m.TotalHuman = human.IBytes(mem.Total)
	m.Used = mem.Used / 1024 / 1024
	m.UsedHuman = human.IBytes(mem.Used)
	m.Cached = mem.Cached / 1024 / 1024
	m.CachedHuman = human.IBytes(mem.Cached)
	m.Free = mem.Free / 1024 / 1024
	m.FreeHuman = human.IBytes(mem.Free)
	m.UsagePercent = utils.Round(float64(m.Used) / float64(m.Total) * 100.0)
	m.UsagePercentHuman = human.CommafWithDigits(m.UsagePercent, 2) + " %"

	s.Memory = m
	return nil
}

func (s *ServerStats) InitDiskInfos(config *models.Config) error {
	parts, err := disk.Partitions(true)
	if err != nil {
		return err
	}
	for _, p := range parts {
		if p.Mountpoint == "/" {
			continue
		}
		if strings.HasPrefix(config.General.RootFolderPath, p.Mountpoint) {
			device := p.Mountpoint
			u, err := disk.Usage(device)
			if err != nil {
				return err
			}
			if u == nil || u.Total == 0 {
				return errors.New("null reference exceptions")
			}

			d := DiskInfo{}
			d.MountPoint = p.Mountpoint
			d.Fstype = u.Fstype
			d.Total = u.Total / 1024 / 1024
			d.TotalHuman = human.Bytes(u.Total)
			d.Used = u.Used / 1024 / 1024
			d.UsedHuman = human.Bytes(u.Used)
			d.Free = u.Free / 1024 / 1024
			d.FreeHuman = human.Bytes(u.Free)
			d.UsagePercent = utils.Round(u.UsedPercent)
			d.UsagePercentHuman = human.CommafWithDigits(d.UsagePercent, 2) + " %"
			s.Disk = d

			break
		}
	}
	return nil
}

func (s *ServerStats) InitNetworkInfos() error {
	current, err := network.Get()
	if err != nil {
		return err
	}

	var down uint64 = 0
	var up uint64 = 0
	for _, c := range current {
		down += c.RxBytes
		up += c.TxBytes
	}

	n := NetworkInfo{}
	n.Download = down / 1024 / 1024
	n.DownloadHuman = human.Bytes(down)
	n.Upload = up / 1024 / 1024
	n.UploadHuman = human.Bytes(up)
	s.Network = n

	return nil
}
