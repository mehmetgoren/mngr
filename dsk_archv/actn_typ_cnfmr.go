package dsk_archv

import (
	"mngr/models"
	"strings"
)

const (
	Delete            = 0
	MoveToNewLocation = 1
)

type ActionTypeConfirmer struct {
	Config *models.Config
}

func (a *ActionTypeConfirmer) GetActionType() int {
	c := a.Config
	if c.Archive.ActionType == Delete {
		return Delete
	}
	if len(c.Archive.MoveLocation) == 0 {
		return Delete
	}
	if strings.Contains(c.Archive.MoveLocation, c.General.RootFolderPath) { //means move location is on the same physical disk
		return Delete
	}

	return MoveToNewLocation
}
