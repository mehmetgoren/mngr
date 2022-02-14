package utils

import (
	"mngr/reps"
)

var connMain = reps.CreateRedisConnection(reps.MAIN)
var ConfigRep = reps.ConfigRepository{Connection: connMain}

var SourceRep = reps.SourceRepository{Connection: connMain}
var StreamRep = reps.StreamRepository{Connection: connMain}

var ConnPubSub = reps.CreateRedisConnection(reps.EVENTBUS)
