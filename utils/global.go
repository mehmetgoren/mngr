package utils

import (
	"mngr/reps"
)

var connMain = CreateRedisConnection(MAIN)
var ConfigRep = reps.ConfigRepository{Connection: connMain}

var ConnSources = CreateRedisConnection(SOURCES)
var SourceRep = reps.SourceRepository{Connection: ConnSources}

var ConnPubSub = CreateRedisConnection(EVENTBUS)
