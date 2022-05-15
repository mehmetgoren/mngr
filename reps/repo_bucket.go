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
	NdRep     *NetworkDiscoveryRepository
	OvRep     *OnvifRepository
}

func (r *RepoBucket) Init() *RepoBucket {
	r.connMain = createRedisConnection(MAIN)
	r.PubSubConnection = createRedisConnection(EVENTBUS)

	r.ConfigRep = &ConfigRepository{Connection: r.connMain}
	r.SourceRep = &SourceRepository{Connection: r.connMain}
	r.StreamRep = &StreamRepository{Connection: r.connMain}
	r.OdRep = &OdRepository{Connection: r.connMain}
	r.NdRep = &NetworkDiscoveryRepository{Connection: r.connMain}
	r.OvRep = &OnvifRepository{Connection: r.connMain}

	return r
}

func (r *RepoBucket) GetMainConnection() *redis.Client {
	return r.connMain
}
