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


func SelectConversationsList (ctx *fasthttp.RequestCtx){
	userId := 1
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))
	offset := functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))
	query := "select select_conversations_list($1,$2,$3)"
	var result json.RawMessage
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId,limit,offset).Scan(&result);
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

func SelectConversationMessages (ctx *fasthttp.RequestCtx){
	userId1 := 1
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
	messageFrom := 1
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
	MessengerWebsocketStruct.PushMessageToConnections(messageFrom, messageTo, result)
	ctx.SetStatusCode(400)
}

var upgrader = websocket.FastHTTPUpgrader{CheckOrigin: func(ctx *fasthttp.RequestCtx) bool { return true}}

func MessengerHandler (ctx *fasthttp.RequestCtx){
	userId := 1
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