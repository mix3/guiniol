package guiniol

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type EventCtx struct {
	connection   *Connection
	messageEvent MessageEvent
}

func (ctx *EventCtx) Connection() *Connection {
	return ctx.connection
}

func (ctx *EventCtx) MessageEvent() MessageEvent {
	return ctx.messageEvent
}

func (ctx *EventCtx) Token() string {
	return ctx.connection.token
}

func (ctx *EventCtx) UserId() string {
	return ctx.connection.userId
}

func (ctx *EventCtx) UserName() string {
	return ctx.connection.userName
}

func (ctx *EventCtx) Domain() string {
	return ctx.connection.domain
}

func (ctx *EventCtx) Permalink() string {
	tss := strings.Split(ctx.messageEvent.Ts, ".")
	if len(tss) != 2 {
		return ""
	}
	return fmt.Sprintf(
		"https://%s.slack.com/archives/%s/p%s%s",
		ctx.Domain(),
		ctx.ChannelIdToName(ctx.messageEvent.Channel),
		tss[0],
		tss[1],
	)
}

func (ctx *EventCtx) UserIdToName(userId string) string {
	v, ok := ctx.connection.userMap[userId]
	if ok {
		return v
	}
	return ""
}

func (ctx *EventCtx) ChannelIdToName(channelId string) string {
	v, ok := ctx.connection.channelMap[channelId]
	if ok {
		return v
	}
	return ""
}

func (ctx *EventCtx) Reply(message string) error {
	resp, err := http.PostForm("https://slack.com/api/chat.postMessage", url.Values{
		"token":   {ctx.connection.token},
		"channel": {ctx.messageEvent.Channel},
		"text":    {message},
		"as_user": {"true"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (ctx *EventCtx) SendMessage(channel, message string) error {
	resp, err := http.PostForm("https://slack.com/api/chat.postMessage", url.Values{
		"token":   {ctx.connection.token},
		"channel": {channel},
		"text":    {message},
		"as_user": {"true"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
