package main

import (
	"bytes"
	"context"
	"coursework/functools"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"log"
)

type conversationList struct {
	Conversations json.RawMessage
	Done          bool
}

func SelectConversationsList(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))
	offset := functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))
	cList := conversationList{emptyArray, false}

	query := "select select_conversations_list($1,$2,$3)"
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, limit, offset).Scan(&cList.Conversations); err != nil {
		ctx.Error("параметры не верны", 400)
		return
	}
	if bytes.Equal(cList.Conversations, null) {
		cList.Conversations = emptyArray
		cList.Done = true
	}

	jsonResult, _ := json.Marshal(cList)

	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

type ConversationMessagesStruct struct {
	MessagesList   json.RawMessage
	ConversationId int
	Done           bool
}

func SelectConversationMessages(ctx *fasthttp.RequestCtx) {
	userId1 := ctx.UserValue("requestUserId").(int)
	userId2 := functools.ByteSliceToString(ctx.QueryArgs().Peek("userId2"))
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))
	offset := functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))
	cms := ConversationMessagesStruct{emptyArray, 0, false}

	query := "select * from select_conversation_messages($1,$2,$3,$4)"
	_ = Postgres.Conn.QueryRow(context.Background(), query, userId1, userId2, limit, offset).Scan(&cms.MessagesList, &cms.ConversationId)

	if bytes.Equal(cms.MessagesList, null) {
		cms.MessagesList = emptyArray
		cms.Done = true
	}

	jsonResult, _ := json.Marshal(cms)

	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}



var upgrader = websocket.FastHTTPUpgrader{CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true }}

type MessageStruct struct {
	MessageTo   int    `json:"messageTo"`
	MessageText string `json:"messageText"`
}

func MessengerHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	err := upgrader.Upgrade(ctx, func(wconn *websocket.Conn) {
		MessengerWebsocketStruct.AddConn(userId, wconn)
		ms := &MessageStruct{}
		var jsonRes json.RawMessage
		for {
			err := wconn.ReadJSON(&jsonRes)
			if err != nil {
				MessengerWebsocketStruct.RemoveConn(userId, wconn)
				fmt.Println("closed" + string(userId) + err.Error() )
				break
			}
			err = json.Unmarshal(jsonRes, &ms)
			if err != nil{
				log.Println("Неправильная структура сообщения от пользователя", userId, functools.ByteSliceToString(jsonRes))
				return
			}
			PushMessage(userId, ms)
		}
	})
	if err != nil {
		fmt.Println("cannot establish upgrade connection")
	}
}

func PushMessage(messageFrom int, messageStruct *MessageStruct) {
	var result json.RawMessage
	query := "select push_message ($1,$2,$3)"
	if err := Postgres.Conn.QueryRow(context.Background(), query, messageFrom, messageStruct.MessageTo, messageStruct.MessageText).Scan(&result); err != nil {
		return
	}
	MessengerWebsocketStruct.PushMessageToConnections(messageStruct.MessageTo, result)
	if messageStruct.MessageTo != messageFrom {
		MessengerWebsocketStruct.PushMessageToConnections(messageFrom, result)
	}
}


func MessengerGetShortProfileInfo(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	conversationId := functools.ByteSliceToString(ctx.QueryArgs().Peek("conversationId"))

	var result json.RawMessage
	query := "select get_short_profile_info($1, $2)"
	if err := Postgres.Conn.QueryRow(context.Background(), query, conversationId, userId).Scan(&result); err != nil {
		ctx.Error("параметры не верны", 400)
		return
	}

	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}
