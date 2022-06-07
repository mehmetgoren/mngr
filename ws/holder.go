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
)

type Holder struct {
	EventBus     *eb.EventBus
	Client       *Client
	EventHandler eb.EventHandler
}

type Holders struct {
	Rb *reps.RepoBucket

	WsStartStream  map[string]*Holder
	WsStopStream   map[string]*Holder
	WsEditor       map[string]*Holder
	WsFFmpegReader map[string]*Holder
	WsOnvif        map[string]*Holder
	WsVideoMerge   map[string]*Holder
	WsFrTrain      map[string]*Holder
}

func (h *Holders) Init() {
	closeWsConnection(h.WsStartStream)
	h.WsStartStream = make(map[string]*Holder)

	closeWsConnection(h.WsStopStream)
	h.WsStopStream = make(map[string]*Holder)

	closeWsConnection(h.WsEditor)
	h.WsEditor = make(map[string]*Holder)

	closeWsConnection(h.WsFFmpegReader)
	h.WsFFmpegReader = make(map[string]*Holder)

	closeWsConnection(h.WsOnvif)
	h.WsOnvif = make(map[string]*Holder)

	closeWsConnection(h.WsVideoMerge)
	h.WsVideoMerge = make(map[string]*Holder)

	closeWsConnection(h.WsFrTrain)
	h.WsFrTrain = make(map[string]*Holder)
}

func closeWsConnection(h map[string]*Holder) {
	if h == nil {
		return
	}
	for _, v := range h {
		err := v.Client.Close()
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
		default:
			panic("not supported")
		}
	}

	return true
}

func (h *Holders) createStartStreamEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "start_stream_response"}
	eh := &eb.StartStreamResponseEvent{Rb: h.Rb, Pusher: c}
	go e.Subscribe(eh)
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createStopStreamEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "stop_stream_response"}
	eh := &eb.StopStreamResponseEvent{Pusher: c}
	go e.Subscribe(eh)
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createEditorEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "editor_response"}
	eh := &eb.EditorResponseEvent{Pusher: c}
	go e.Subscribe(eh)
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createFFmpegReaderEvent(c *Client, sourceId string) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "ffrs" + sourceId}
	eh := &eb.FFmpegReaderResponseEvent{Pusher: c}
	go e.Subscribe(eh)
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createOnvifEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "onvif_response"}
	eh := &eb.OnvifResponseEvent{Pusher: c}
	go e.Subscribe(eh)
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createVideoMergeEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "vfm_response"}
	eh := &eb.VideMergeResponseEvent{Pusher: c}
	go e.Subscribe(eh)
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}

func (h *Holders) createFrTrainEvent(c *Client) *Holder {
	e := &eb.EventBus{PubSubConnection: h.Rb.PubSubConnection, Channel: "fr_train_response"}
	eh := &eb.FaceTrainResponseEvent{Pusher: c}
	go e.Subscribe(eh)
	return &Holder{
		EventBus:     e,
		Client:       c,
		EventHandler: eh,
	}
}
