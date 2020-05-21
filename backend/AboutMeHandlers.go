package main

import (
	"bytes"
	"context"
	"coursework/functools"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

func SelectExtendedUserInfo (ctx *fasthttp.RequestCtx){
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))

	query := "select get_extended_info($1)"
	var result json.RawMessage
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result); err != nil {
		ctx.SetStatusCode(400)
		return
	}
	if bytes.Equal(result,null){
		result = emptyArray
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}