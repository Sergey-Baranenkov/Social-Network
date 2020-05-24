package main

import (
	"bytes"
	"context"
	"coursework/functools"
	"encoding/json"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"strconv"
)

type conversationList struct {
	Conversations json.RawMessage
	Done bool
}
func SelectConversationsList (ctx *fasthttp.RequestCtx){
	userId := ctx.UserValue("requestUserId").(int)
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))
	offset := functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))

	cList := conversationList{emptyArray, false}

	query := "select select_conversations_list($1,$2,$3)"
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId,limit,offset).Scan(&cList.Conversations);
		err != nil {
			fmt.Println(err)
			ctx.Error("параметры не верны", 400)
			return
	}
	if bytes.Equal(cList.Conversations, null){
		cList.Conversations = emptyArray
		cList.Done = true
	}

	jsonResult, _ := json.Marshal(cList)

	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}


func SelectConversationMessages (ctx *fasthttp.RequestCtx){
	userId1 := ctx.UserValue("requestUserId").(int)
	userId2 := functools.ByteSliceToString(ctx.QueryArgs().Peek("userId2"))
	query := "select select_conversation_messages($1,$2)"
	var result json.RawMessage
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId1,userId2).Scan(&result);
		err != nil {
			fmt.Println(err)
		ctx.Error("параметры не верны", 400)
		return
	}
	if bytes.Equal(result, null){
		result = emptyArray
	}

	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}

func PushMessage (ctx *fasthttp.RequestCtx){
	messageFrom := ctx.UserValue("requestUserId").(int)
	fmt.Println(functools.ByteSliceToString(ctx.QueryArgs().Peek("messageTo")))
	messageTo, err := strconv.Atoi(functools.ByteSliceToString(ctx.QueryArgs().Peek("messageTo")))
	if err != nil{
		ctx.Error("параметры не верны", 400)
		return
	}
	messageText := functools.ByteSliceToString(ctx.QueryArgs().Peek("messageText"))

	var result json.RawMessage
	query := "select push_message ($1,$2,$3)"
	if err := Postgres.Conn.QueryRow(context.Background(), query, messageFrom, messageTo, messageText).Scan(&result);
		err != nil {
		ctx.Error("параметры не верны", 400)
		return
	}
	MessengerWebsocketStruct.PushMessageToConnections(messageTo, result)
	if  messageTo != messageFrom{
		MessengerWebsocketStruct.PushMessageToConnections(messageFrom, result)
	}

	ctx.SetStatusCode(400)
}

var upgrader = websocket.FastHTTPUpgrader{CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true}}
func MessengerHandler (ctx *fasthttp.RequestCtx){
	userId := ctx.UserValue("requestUserId").(int)
	err := upgrader.Upgrade(ctx, func(wconn *websocket.Conn){
		MessengerWebsocketStruct.AddConn(userId, wconn)
		fmt.Println("connected")
		var response json.RawMessage
		for {
			_, _, err := wconn.NextReader()
			if err != nil{
				MessengerWebsocketStruct.RemoveConn(userId,wconn)
				fmt.Println("closed" + string(userId))
				break
			}
			fmt.Println(functools.ByteSliceToString(response))
		}
	})
	if err != nil{
		fmt.Println("cannot establish upgrade connection")
	}
}

func MessengerGetShortProfileInfo(ctx *fasthttp.RequestCtx){
	userId := ctx.UserValue("requestUserId").(int)
	conversationId := functools.ByteSliceToString(ctx.QueryArgs().Peek("conversationId"))

	var result json.RawMessage
	query := "select get_short_profile_info($1, $2)"
	if err := Postgres.Conn.QueryRow(context.Background(), query, conversationId, userId).Scan(&result);
		err != nil {
		ctx.Error("параметры не верны", 400)
		return
	}

	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}