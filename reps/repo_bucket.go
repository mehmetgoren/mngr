package reps

import (
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"strconv"
	"time"
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
	CloudRep        *CloudRepository

	users map[string]*models.UserSession
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
	r.CloudRep = &CloudRepository{Connection: r.connMain}

	r.initUsers()
	r.initSessionClearScheduler()

	return r
}

func (r *RepoBucket) initUsers() {
	users, err := r.UserRep.GetUsers()
	if err != nil {
		log.Println(err.Error())
	} else {
		r.users = make(map[string]*models.UserSession)
		for _, user := range users {
			r.users[user.Token] = &models.UserSession{}
			r.users[user.Token].User = user
		}
	}
}

func (r *RepoBucket) AddUser(user *models.User) {
	u, found := r.users[user.Token]
	if !found {
		u = &models.UserSession{}
		r.users[user.Token] = u
		u.User = user
	}
	u.LastVisitAt = time.Now()
}

func (r *RepoBucket) RemoveUser(token string) {
	delete(r.users, token)
	log.Println("a user with a token " + token + " has been deleted")
}

func (r *RepoBucket) IsUserAuthenticated(token string) (*models.UserSession, bool) {
	user, found := r.users[token]
	if found {
		user.LastVisitAt = time.Now()
	}
	return user, found
}

func (r *RepoBucket) GetMainConnection() *redis.Client {
	return r.connMain
}

func (r *RepoBucket) initSessionClearScheduler() {
	sessionTimeout := 1200. // 20 minutes as seconds
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(func() {
		now := time.Now()
		removeTokenList := make([]string, 0)
		for token, v := range r.users {
			diff := now.Sub(v.LastVisitAt)
			log.Println("session timeout check time diff as seconds: " + strconv.Itoa(int(diff.Seconds())))
			if diff.Seconds() >= sessionTimeout {
				log.Println("An inactive user has been detected and will be removed: "+v.Username, " token: "+v.Token)
				removeTokenList = append(removeTokenList, token)
			}
		}
		for _, token := range removeTokenList {
			r.RemoveUser(token)
		}
		log.Println("Session timeout has been checked")
	})
	s.StartAsync()
}
