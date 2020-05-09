package main

import (
	"context"
	"coursework/functools"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

func GetUserVideoHandler(ctx *fasthttp.RequestCtx){
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	startFrom:= functools.ByteSliceToString(ctx.QueryArgs().Peek("startFrom"))


	userVideo := make([]byte, 0, 1024)
	query:= "select json_agg(m) from (select video_id, name, adder_id from video where video_id in (select unnest(video_list) from users where user_id = $1 limit 8 offset $2)) m;"

	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, startFrom).Scan(&userVideo);
		err != nil {
		fmt.Println(err)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = ctx.WriteString(functools.ByteSliceToString(userVideo))


}

type CombinedVideoJSON struct {
	UserVideo string
	AllVideo string
}

func GetCombinedVideoHandler(ctx *fasthttp.RequestCtx){
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	startFrom:= functools.ByteSliceToString(ctx.QueryArgs().Peek("startFrom"))
	withVal := functools.ByteSliceToString(ctx.QueryArgs().Peek("withVal"))

	combVideo := CombinedVideoJSON{"[]","[]"}
	count:= 0

	if startFrom == "0" {
		userVideo := make([]byte, 0, 1024)
		query:= "select count(*), json_agg(m) from (select video_id, name, adder_id from video where video_id in (select unnest(video_list) from users where user_id = $1 and document @@ plainto_tsquery($2))) m;"
		if err := Postgres.Conn.QueryRow(context.Background(), query, userId, withVal).Scan(&count, &userVideo); err != nil {
			fmt.Println(err)
			return
		}
		combVideo.UserVideo = functools.ByteSliceToString(userVideo)
	}

	fmt.Println(count)
	limit:= 8
	if count < limit {
		allVideo := make([]byte, 0, 1024)
		query:= "select json_agg(m) from (select video_id, name, adder_id from video where document @@ plainto_tsquery($1) limit $2 offset $3) m; "
		if err := Postgres.Conn.QueryRow(context.Background(), query, withVal, limit - count, startFrom).Scan(&count, &allVideo); err != nil {
			fmt.Println(err)
			return
		}
		combVideo.AllVideo = functools.ByteSliceToString(allVideo)
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	jsonResult, _ := json.Marshal(combVideo);
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}