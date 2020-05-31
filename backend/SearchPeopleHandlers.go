package main

import (
	"bytes"
	"context"
	"coursework/functools"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type SearchedPeopleStruct struct {
	People json.RawMessage
	Done bool
}

func GetSearchedPeople(ctx * fasthttp.RequestCtx){
	value:= functools.ByteSliceToString(ctx.QueryArgs().Peek("value"))
	offset:= functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))
	limit:= functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))

	sps := SearchedPeopleStruct{null, false}

	query := `select json_agg(m) from (select user_id, first_name, last_name from users where full_name @@ plainto_tsquery($1) offset $2 limit $3) m;`
	if err := Postgres.Conn.QueryRow(context.Background(), query, value, offset, limit).Scan(&sps.People);
		err != nil {
		ctx.SetStatusCode(400)
		return
	}
	
	if bytes.Equal(sps.People, null){
		sps.People = emptyArray
		sps.Done = true
	}

	result, _ := json.Marshal(sps)

	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}