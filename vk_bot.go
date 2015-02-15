package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ZurgInq/vk_bot/notes"
	"github.com/ZurgInq/vk_bot/sysstat"
	"github.com/parnurzeal/gorequest"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const resultPrefix = "~"

var (
	baseUrl      = "https://api.vk.com/method/"
	access_token = flag.String("token", "", "vk access token")
	vk_id        = flag.Int64("vk_id", 0, "user vk id")
	request      = gorequest.New()
)

type MessagesResponse struct {
	Response Response `json:"response,omitempty"`
}

type Response struct {
	Count int64     `json:"count,omitempty"`
	Items []MsgItem `json:"items,omitempty"`
}

type MsgItem struct {
	Id        int64  `json:"id,omitempty"`
	Date      int64  `json:"date,omitempty"`
	Out       int64  `json:"out,omitempty"`
	UserId    int64  `json:"user_id,omitempty"`
	ReadState int64  `json:"read_state,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
}

func sendMsg(userId int64, msg string) {
	userIdStr := strconv.FormatInt(userId, 10)
	_, body, errs := request.Get(baseUrl + "messages.send").
		Query("v=5.28").
		Query("user_id=" + userIdStr).
		Query("message=" + resultPrefix + msg).
		Query("access_token=" + *access_token).
		End()
	fmt.Println(errs)
	fmt.Println(body)
}

func getMsgs(lastMessageId int64) (MessagesResponse, error) {
	prepareRequest := request.Get(baseUrl + "messages.get")

	if lastMessageId != 0 {
		lastMsgId := strconv.FormatInt(lastMessageId, 10)
		prepareRequest = prepareRequest.Query("last_message_id=" + lastMsgId)
	}
	_, body, errs := prepareRequest.
		Query("v=5.28").
		Query("out=0").
		Query("count=20").
		Query("access_token=" + *access_token).
		End()

	if len(errs) > 0 {
		fmt.Println(errs)
		return MessagesResponse{}, errs[0]
	}

	result := &MessagesResponse{}
	err := json.Unmarshal([]byte(body), result)

	if err != nil {
		fmt.Println(err)
		return MessagesResponse{}, err
	}

	return *result, nil
}

func sigintHandle() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			os.Exit(0)
		}
	}()
}

func checkResultPrefix(msg string) bool {
	return !strings.HasPrefix(msg, resultPrefix)
}

func checkVkId(msgUserId int64) bool {
	return msgUserId == *vk_id
}

func checkTime(ts int64) bool {
	return (ts+10 > time.Now().Unix())
}

func doCmd(msg string) {
	switch true {
	case strings.HasPrefix(msg, "@"):
		doSysCmd(msg)
	case strings.HasPrefix(msg, "!"):
		doNoteCmd(msg)
	}
}

func doSysCmd(msg string) {
	switch true {
	case strings.HasPrefix(msg, "@sys/host"):
		sendMsg(*vk_id, fmt.Sprintf("%+v", sysstat.GetHost()))
	case strings.HasPrefix(msg, "@sys/disk"):
		sendMsg(*vk_id, fmt.Sprintf("%+v", sysstat.GetDisk("/")))
	case strings.HasPrefix(msg, "@sys/load"):
		sendMsg(*vk_id, fmt.Sprintf("%+v", sysstat.GetLoad()))
	case strings.HasPrefix(msg, "@sys/ram"):
		sendMsg(*vk_id, fmt.Sprintf("%+v", sysstat.GetRam()))
	case strings.HasPrefix(msg, "@sys"):
		sendMsg(*vk_id, fmt.Sprintf("\n%+v", sysstat.GetSystem("/")))
	case strings.HasPrefix(msg, "@sh"):
		args := strings.SplitN(msg, " ", 2)
		if len(args) < 2 {
			return
		}
		fmt.Println("exec: ", args[1])
		sendMsg(*vk_id, fmt.Sprintf(" %+v", sysstat.ExecShell(args[1])))
	}
}

func doNoteCmd(msg string) {
	switch true {
	case strings.HasPrefix(msg, "!add"):
		args := strings.SplitN(msg, " ", 2)
		if len(args) < 2 {
			return
		}
		notes.Add(args[1])
		sendMsg(*vk_id, "ok")
	case strings.HasPrefix(msg, "!del"):
		args := strings.SplitN(msg, " ", 2)
		if len(args) < 2 {
			return
		}
		if id, err := strconv.Atoi(args[1]); err != nil {
			fmt.Println("Error: ", err)
			sendMsg(*vk_id, "Error: "+err.Error())
		} else {
			notes.Del(id)
		}
		sendMsg(*vk_id, "ok")
	case strings.HasPrefix(msg, "!list"):
		sendMsg(*vk_id, "\n"+notes.List())
	}
}

func main() {
	flag.Parse()
	if *access_token == "" {
		fmt.Println("access_token empty")
		return
	}

	if *vk_id == 0 {
		fmt.Println("vk_id empty")
		return
	}

	fmt.Println("start")
	sigintHandle()

	c := time.Tick(2 * time.Second)
	lastMsgId := int64(0)
	for _ = range c {
		msgs, err := getMsgs(lastMsgId)
		if err != nil {
			fmt.Println(err)
		} else {
			if len(msgs.Response.Items) > 0 {
				lastMsgId = msgs.Response.Items[0].Id
			}
			for _, msg := range msgs.Response.Items {
				msgBody := strings.Trim(msg.Body, "")
				if checkVkId(msg.UserId) && checkTime(msg.Date) && checkResultPrefix(msgBody) {
					go doCmd(msgBody)
				}
			}
		}
	}
}
