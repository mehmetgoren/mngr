package reps

import (
	"github.com/go-redis/redis/v8"
)

type RepoBucket struct {
	connMain         *redis.Client
	PubSubConnection *redis.Client

	ConfigRep *ConfigRepository
	SourceRep *SourceRepository
	StreamRep *StreamRepository
	OdRep     *OdRepository
}

func (r *RepoBucket) Init() *RepoBucket {
	r.connMain = CreateRedisConnection(MAIN)
	r.PubSubConnection = CreateRedisConnection(EVENTBUS)

	r.ConfigRep = &ConfigRepository{Connection: r.connMain}
	r.SourceRep = &SourceRepository{Connection: r.connMain}
	r.StreamRep = &StreamRepository{Connection: r.connMain}
	r.OdRep = &OdRepository{Connection: r.connMain}

	return r
}
