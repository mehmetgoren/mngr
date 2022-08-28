package ws

import (
	"github.com/gin-gonic/gin"
	"log"
	"mngr/eb"
	"mngr/reps"
	"mngr/utils"
	"strconv"
)

type Holder struct {
	EventBus     *eb.EventBus
	Client       *Client
	EventHandler eb.EventHandler
}

type Event interface {
	GetOp() int
	GetChannelName(keyExtended string) string
	CreateEventHandler() eb.EventHandler
}

func getKey(op int, token string) string {
	return strconv.Itoa(op) + token
}

type Holders struct {
	Rb     *reps.RepoBucket
	Events []Event

	dic map[string]*Holder
}

func (h *Holders) Init() {
	h.Events = make([]Event, 0)
	h.Events = append(h.Events, &UserEvent{})
	h.Events = append(h.Events, &StartStreamEvent{})
	h.Events = append(h.Events, &StopStreamEvent{})
	h.Events = append(h.Events, &EditorEvent{})
	h.Events = append(h.Events, &FFmpegReaderEvent{})
	h.Events = append(h.Events, &OnvifEvent{})
	h.Events = append(h.Events, &VideoMergeEvent{})
	h.Events = append(h.Events, &FrTrainEvent{})
	h.Events = append(h.Events, &ProbeEvent{})
	h.Events = append(h.Events, &NotifierEvent{})

	h.dic = make(map[string]*Holder)
}

func (h *Holders) UserLogin(token string, client *Client) {
	h.dic[getKey(User, token)] = h.CreateEvent(User, client, "")
}
func (h *Holders) UserLogout(token string, triggerLogout bool) {
	for j := StartStream; j < len(h.Events); j++ {
		if val, ok := h.dic[token]; ok {
			err := val.Client.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	if val, ok := h.dic[getKey(User, token)]; ok {
		if triggerLogout {
			err := val.Client.Push(token)
			if err != nil {
				log.Println(err.Error())
			}
		}
		err := val.Client.Close()
		if err != nil {
			log.Println(err.Error())
		}
		delete(h.dic, getKey(User, token))
	}
}

func (h *Holders) CreateEvent(op int, c *Client, keyExtended string) *Holder {
	evt := h.Events[op]
	var e *eb.EventBus
	eh := evt.CreateEventHandler()
	if eh != nil { // like logout events
		eh.SetPusher(c)
		e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: evt.GetChannelName(keyExtended)}
		go func() {
			err := e.Subscribe(eh)
			if err != nil {
				log.Println(err.Error())
			}
		}()
	}
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) RegisterEndPoint(hub *Hub, ctx *gin.Context, opType int, extendedKey string) bool {
	token := utils.GetQsValue(ctx, "token")
	if len(token) == 0 {
		log.Println("invalid user token")
		return false
	}
	token += extendedKey
	client := CreateClient(hub, ctx.Writer, ctx.Request)
	if prev, ok := h.dic[getKey(opType, token)]; ok {
		//err := prev.Client.Close()
		//if err != nil {
		//	log.Println("Error while closing prev websockets connection. Err: ", err)
		//}
		prev.Client = client
		prev.EventHandler.SetPusher(client)
		log.Println("holder's item has been already added,changing Ws Client for " + token)
	} else {
		h.dic[getKey(opType, token)] = h.CreateEvent(opType, client, extendedKey)
	}
	return true
}
