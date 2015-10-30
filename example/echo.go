package main

import (
	"fmt"
	"log"
	"regexp"

	"golang.org/x/net/context"

	"github.com/mix3/guiniol"
	"github.com/naoya/go-pit"
)

var pattern1 = regexp.MustCompile(`^([^:]+):\s+(.+)`)   // bot_name: *****
var pattern2 = regexp.MustCompile(`^<@(\w+)>:?\s+(.+)`) // @bot_name *****

func matcher(e *guiniol.EventCtx) (string, string) {
	t := e.MessageEvent().Text
	if m := pattern2.FindStringSubmatch(t); len(m) == 3 {
		return m[1], m[2]
	}
	if m := pattern1.FindStringSubmatch(t); len(m) == 3 {
		return m[1], m[2]
	}
	return "", ""
}

const (
	MessageKey = "MessageKey"
)

func Message(c context.Context) string {
	m, ok := c.Value(MessageKey).(string)
	if ok {
		return m
	}
	return ""
}

func Wrap(fn guiniol.CallbackFunc) guiniol.Callback {
	return guiniol.CallbackFunc(func(c context.Context, e *guiniol.EventCtx) {
		if e.MessageEvent().Subtype == "bot_message" {
			return
		}

		name, m := matcher(e)
		if name != e.UserName() && name != e.UserId() {
			return
		}

		c = context.WithValue(c, MessageKey, m)

		fn.Next(c, e)
	})
}

func EchoCallback(c context.Context, e *guiniol.EventCtx) {
	e.Reply(fmt.Sprintf(
		"@%s echo %s\n%s",
		e.UserIdToName(e.MessageEvent().User),
		Message(c),
		e.Permalink(),
	))
}

func main() {
	config, err := pit.Get("guiniol")
	if err != nil {
		log.Fatal(err)
	}

	conn := guiniol.NewConnection(config["slack.token"])
	conn.RegisterCb(Wrap(EchoCallback))
	conn.Loop()
}
