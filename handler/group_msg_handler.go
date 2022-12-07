package handler

import (
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/wechatgpt/wechatbot/config"
	"github.com/wechatgpt/wechatbot/openai"
	"log"
	"strings"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

type GroupMessageHandler struct {
}

func (gmh *GroupMessageHandler) handle(msg *openwechat.Message) error {
	if !msg.IsText() {
		return nil
	}
	return gmh.ReplyText(msg)
}

func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

func (gmh *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {

	sender, err := msg.Sender()
	group := openwechat.Group{User: sender}
	log.Printf("Received Group %v Text Msg : %v", group.NickName, msg.Content)

	appConfig := config.GetConfig()
	if appConfig == nil {
		return errors.New("can not get appConfig file,please check")
	}

	if !strings.Contains(msg.Content, appConfig.ChatGpt.Keyword) {
		return nil
	}
	splitItems := strings.Split(msg.Content, appConfig.ChatGpt.Keyword)
	if len(splitItems) < 2 {
		return nil
	}
	requestText := strings.TrimSpace(splitItems[1])
	reply, err := openai.Completions(requestText, appConfig.ChatGpt.Token)
	if err != nil {
		log.Println(err)
		text, err := msg.ReplyText(fmt.Sprintf("bot error: %s", err.Error()))
		if err != nil {
			return err
		}
		fmt.Println(text)
		return err
	}

	if reply != nil {
		_, err = msg.ReplyText(*reply)
		if err != nil {
			log.Println(err)
		}
		return err
	}

	return nil
}
