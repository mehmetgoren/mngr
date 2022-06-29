package ws

import (
	"github.com/gin-gonic/gin"
	"log"
	"mngr/eb"
	"mngr/reps"
	"mngr/utils"
)

const (
	StartStreamEvent  = 0
	StopStreamEvent   = 1
	EditorEvent       = 2
	FFmpegReaderEvent = 3
	OnvifEvent        = 4
	VideoMergeEvent   = 5
	FrTrainEvent      = 6
	ProbeEvent        = 7
	NotifierEvent     = 8
)

type Holder struct {
	EventBus     *eb.EventBus
	Client       *Client
	EventHandler eb.EventHandler
}

type Holders struct {
	Rb *reps.RepoBucket

	userLogoutConnections map[string]*Client

	WsStartStream  map[string]*Holder
	WsStopStream   map[string]*Holder
	WsEditor       map[string]*Holder
	WsFFmpegReader map[string]*Holder
	WsOnvif        map[string]*Holder
	WsVideoMerge   map[string]*Holder
	WsFrTrain      map[string]*Holder
	Probe          map[string]*Holder
	Notifier       map[string]*Holder
}

func (h *Holders) Init() {
	h.userLogoutConnections = make(map[string]*Client)
	h.WsStartStream = make(map[string]*Holder)
	h.WsStopStream = make(map[string]*Holder)
	h.WsEditor = make(map[string]*Holder)
	h.WsFFmpegReader = make(map[string]*Holder)
	h.WsOnvif = make(map[string]*Holder)
	h.WsVideoMerge = make(map[string]*Holder)
	h.WsFrTrain = make(map[string]*Holder)
	h.Probe = make(map[string]*Holder)
	h.Notifier = make(map[string]*Holder)
}

func (h *Holders) UserLogin(token string, client *Client) {
	h.userLogoutConnections[token] = client
}

func (h *Holders) UserLogout(token string, triggerLogout bool) {
	closeWsConnection(h.WsStartStream, token)
	closeWsConnection(h.WsStopStream, token)
	closeWsConnection(h.WsEditor, token)
	closeWsConnection(h.WsFFmpegReader, token)
	closeWsConnection(h.WsOnvif, token)
	closeWsConnection(h.WsVideoMerge, token)
	closeWsConnection(h.WsFrTrain, token)
	closeWsConnection(h.Probe, token)
	closeWsConnection(h.Notifier, token)
	if val, ok := h.userLogoutConnections[token]; ok {
		if triggerLogout {
			err := val.Push(token)
			if err != nil {
				log.Println(err.Error())
			}
		}
		err := val.Close()
		if err != nil {
			log.Println(err.Error())
		}
		delete(h.userLogoutConnections, token)
	}
}

func closeWsConnection(h map[string]*Holder, token string) {
	if h == nil {
		return
	}
	if val, ok := h[token]; ok {
		err := val.Client.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (h *Holders) getDic(opType int) map[string]*Holder {
	switch opType {
	case StartStreamEvent:
		return h.WsStartStream
	case StopStreamEvent:
		return h.WsStopStream
	case EditorEvent:
		return h.WsEditor
	case FFmpegReaderEvent:
		return h.WsFFmpegReader
	case OnvifEvent:
		return h.WsOnvif
	case VideoMergeEvent:
		return h.WsVideoMerge
	case FrTrainEvent:
		return h.WsFrTrain
	case ProbeEvent:
		return h.Probe
	case NotifierEvent:
		return h.Notifier
	default:
		panic("not supported")
	}
}

// RegisterEndPoint keyExtended is optional. pass empty string if it isn't necessary
func (h *Holders) RegisterEndPoint(hub *Hub, ctx *gin.Context, opType int, keyExtended string) bool {
	token := utils.GetQsValue(ctx, "token")
	if len(token) == 0 {
		log.Println("invalid user token")
		return false
	}
	token += keyExtended
	client := CreateClient(hub, ctx.Writer, ctx.Request)
	dic := h.getDic(opType)
	if prev, ok := dic[token]; ok {
		err := prev.Client.Close()
		if err != nil {
			log.Println("Error while closing prev websockets connection for FFmpeg Reader. Err: ", err)
		}
		prev.Client = client
		prev.EventHandler.SetPusher(client)
		log.Println("holder's item has been already added,changing Ws Client for " + token)
	} else { // first timer
		switch opType {
		case StartStreamEvent:
			dic[token] = h.createStartStreamEvent(client)
			break
		case StopStreamEvent:
			dic[token] = h.createStopStreamEvent(client)
			break
		case EditorEvent:
			dic[token] = h.createEditorEvent(client)
			break
		case FFmpegReaderEvent:
			dic[token] = h.createFFmpegReaderEvent(client, keyExtended)
			break
		case OnvifEvent:
			dic[token] = h.createOnvifEvent(client)
			break
		case VideoMergeEvent:
			dic[token] = h.createVideoMergeEvent(client)
			break
		case FrTrainEvent:
			dic[token] = h.createFrTrainEvent(client)
			break
		case ProbeEvent:
			dic[token] = h.createProbeEvent(client)
			break
		case NotifierEvent:
			dic[token] = h.createNotifierEvent(client)
			break
		default:
			panic("not supported")
		}
	}

	return true
}

func (h *Holders) createStartStreamEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "start_stream_response"}
	eh := &eb.StartStreamResponseEvent{Rb: h.Rb, Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createStopStreamEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "stop_stream_response"}
	eh := &eb.StopStreamResponseEvent{Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createEditorEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "editor_response"}
	eh := &eb.EditorResponseEvent{Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createFFmpegReaderEvent(c *Client, sourceId string) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "ffrs" + sourceId}
	eh := &eb.FFmpegReaderResponseEvent{Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createOnvifEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "onvif_response"}
	eh := &eb.OnvifResponseEvent{Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createVideoMergeEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "vfm_response"}
	eh := &eb.VideMergeResponseEvent{Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createFrTrainEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "fr_train_response"}
	eh := &eb.FaceTrainResponseEvent{Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createProbeEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "probe_response"}
	eh := &eb.ProbeResponseEvent{Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createNotifierEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "notifier"}
	eh := &eb.NotifierResponseEvent{Pusher: c}
	go func() {
		err := e.Subscribe(eh)
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}
