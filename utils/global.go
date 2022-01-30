package utils

import (
	"mngr/reps"
)

var connMain = reps.CreateRedisConnection(reps.MAIN)
var ConfigRep = reps.ConfigRepository{Connection: connMain}

var ConnSources = reps.CreateRedisConnection(reps.SOURCES)
var SourceRep = reps.SourceRepository{Connection: ConnSources}

var ConnPubSub = reps.CreateRedisConnection(reps.EVENTBUS)
