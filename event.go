package guiniol

var eventMapping = map[string]interface{}{
	"hello":           HelloEvent{},
	"message":         MessageEvent{},
	"channel_created": ChannelCreatedEvent{},
	"channel_deleted": ChannelDeletedEvent{},
	"channel_rename":  ChannelRenameEvent{},
	"user_change":     UserChangeEvent{},
}

type Type struct {
	Type string `json:"type"`
}

type HelloEvent struct {
	Type string `json:"type"`
}

type MessageEvent struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
	Ts      string `json:"ts"`
}

type ChannelCreatedEvent struct {
	Type    string `json:"type"`
	Channel struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
}

type ChannelDeletedEvent struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

type ChannelRenameEvent struct {
	Type    string `json:"type"`
	Channel struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
}

type UserChangeEvent struct {
	Type string `json:"type"`
	User struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}
