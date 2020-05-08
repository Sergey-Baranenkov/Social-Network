package main

import (
	"context"
	"coursework/functools"
	"github.com/valyala/fasthttp"
)

func GetSearchedPeople(ctx * fasthttp.RequestCtx){
	request:= functools.ByteSliceToString(ctx.QueryArgs().Peek("request"))
	result := make([]byte, 1024)
	query := `select json_agg(m) from (select user_id, first_name, last_name, avatar_ref from users where full_name @@ plainto_tsquery($1)) m;`
	if err := Postgres.Conn.QueryRow(context.Background(), query, request).Scan(&result);
		err != nil {
		_, _ = ctx.WriteString("[{}]")
		return
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}