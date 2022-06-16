package reps

import (
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type RepoBucket struct {
	connMain         *redis.Client
	PubSubConnection *redis.Client

	ConfigRep       *ConfigRepository
	SourceRep       *SourceRepository
	StreamRep       *StreamRepository
	OdRep           *OdRepository
	NdRep           *NetworkDiscoveryRepository
	OvRep           *OnvifRepository
	UserRep         *UserRepository
	ServiceRep      *ServiceRepository
	RtspTemplateRep *RtspTemplateRepository
	FailedStreamRep *FailedStreamRepository
	RecStuckRep     *RecStuckRepository
	VariousRep      *VariousInfosRepository

	users map[string]*models.User
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
	r.UserRep = &UserRepository{Connection: r.connMain}
	r.ServiceRep = &ServiceRepository{Connection: r.connMain}
	r.RtspTemplateRep = &RtspTemplateRepository{Connection: r.connMain}
	r.FailedStreamRep = &FailedStreamRepository{Connection: r.connMain}
	r.RecStuckRep = &RecStuckRepository{Connection: r.connMain}
	r.VariousRep = &VariousInfosRepository{Connection: r.connMain}

	r.initUsers()

	return r
}

func (r *RepoBucket) initUsers() {
	r.users = make(map[string]*models.User)
	users, err := r.UserRep.GetUsers()
	if err != nil {
		log.Println(err.Error())
	} else {
		for _, user := range users {
			r.users[user.Token] = user
		}
	}
}

func (r *RepoBucket) AddUser(user *models.User) {
	r.users[user.Token] = user
}

func (r *RepoBucket) RemoveUser(token string) {
	delete(r.users, token)
}

func (r *RepoBucket) IsUserAuthenticated(token string) (*models.User, bool) {
	user, found := r.users[token]
	return user, found
}

func (r *RepoBucket) GetMainConnection() *redis.Client {
	return r.connMain
}
