package guiniol

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"

	"github.com/k0kubun/pp"
)

type Connection struct {
	token      string
	userId     string
	userName   string
	domain     string
	userMap    map[string]string
	channelMap map[string]string
	cb         Callback
}

func NewConnection(token string) *Connection {
	return &Connection{
		token:      token,
		userMap:    map[string]string{},
		channelMap: map[string]string{},
		cb:         defaultCb,
	}
}

type RTMStart struct {
	Ok    bool   `json:"ok"`
	Url   string `json:"url"`
	Error string `json:"error"`
	Self  struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"self"`
	Team struct {
		Domain string `json:"domain"`
	} `json:"team"`
	Users []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"users"`
	Channels []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"channels"`
}

func (conn *Connection) newWSConnection() (*websocket.Conn, error) {
	resp, err := http.PostForm(
		"https://slack.com/api/rtm.start",
		url.Values{"token": {conn.token}},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r RTMStart
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	conn.userId = r.Self.Id
	conn.userName = r.Self.Name
	conn.domain = r.Team.Domain

	conn.userMap = map[string]string{}
	conn.channelMap = map[string]string{}
	for _, v := range r.Users {
		conn.userMap[v.Id] = v.Name
	}
	for _, v := range r.Channels {
		conn.channelMap[v.Id] = v.Name
	}

	return websocket.Dial(r.Url, "", "https://slack.com/")
}

func (conn *Connection) Loop() {
	for {
		ws, err := conn.newWSConnection()
		if err != nil {
			log.Fatal(err)
		}
		defer ws.Close()
		ws.SetDeadline(time.Now().Add(10 * time.Minute))

		func() {
			for {
				data := json.RawMessage{}
				if err := websocket.JSON.Receive(ws, &data); err != nil {
					log.Printf("failed websocket json receive: %v", err)
					return
				}

				event := &Type{}
				if err := json.Unmarshal(data, event); err != nil {
					log.Printf("failed json unmarshal: %v", err)
					continue
				}

				v, ok := eventMapping[event.Type]
				if !ok {
					continue
				}

				typeOf := reflect.TypeOf(v)
				ep := reflect.New(typeOf).Interface()
				if err := json.Unmarshal(data, ep); err != nil {
					log.Printf("failed json unmarshal for type: %v", err)
					continue
				}

				switch e := ep.(type) {
				case *HelloEvent:
					// ...
				case *MessageEvent:
					conn.CallCb(*e)
				case *ChannelCreatedEvent:
					conn.channelMap[e.Channel.Id] = e.Channel.Name
					pp.Println(conn.channelMap)
				case *ChannelDeletedEvent:
					delete(conn.channelMap, e.Channel)
					pp.Println(conn.channelMap)
				case *ChannelRenameEvent:
					conn.channelMap[e.Channel.Id] = e.Channel.Name
					pp.Println(conn.channelMap)
				case *UserChangeEvent:
					conn.userMap[e.User.Id] = e.User.Name
					pp.Println(conn.userMap)
				default:
				}
			}
		}()
	}

}

func (conn *Connection) CallCb(e MessageEvent) {
	conn.cb.Next(
		context.Background(),
		&EventCtx{
			connection:   conn,
			messageEvent: e,
		},
	)
}

func (conn *Connection) RegisterCb(cb Callback) {
	conn.cb = cb
}
