package main

import (
	"bytes"
	"context"
	"coursework/functools"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
	"strings"
)

type VideoJSON struct {
	UserVideos json.RawMessage
	AllVideos json.RawMessage
}

var null = []byte ("null")
var emptyArray = json.RawMessage("[]")

func GetUserVideoHandler(ctx *fasthttp.RequestCtx){
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	startFrom:= functools.ByteSliceToString(ctx.QueryArgs().Peek("startFrom"))

	VideoStruct := VideoJSON{json.RawMessage("[]"),json.RawMessage("[]")}
	query:= "select json_agg(m) from (select video_id, name, adder_id from video where video_id in (select unnest(video_list) from users where user_id = $1 limit 8 offset $2) order by created_at desc ) m;"
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, startFrom).Scan(&VideoStruct.UserVideos);
		err != nil {
		fmt.Println(err)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	if bytes.Equal(VideoStruct.UserVideos,null){
		VideoStruct.UserVideos = emptyArray
	}

	jsonResult, _ := json.Marshal(VideoStruct)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

func GetCombinedVideoHandler(ctx *fasthttp.RequestCtx){
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	startFrom:= functools.ByteSliceToString(ctx.QueryArgs().Peek("startFrom"))
	withVal := functools.ByteSliceToString(ctx.QueryArgs().Peek("withVal"))
	VideoStruct := VideoJSON{json.RawMessage("[]"),json.RawMessage("[]")}
	count:= 0

	if startFrom == "0" {
		query:= "select count(*), json_agg(m) from (select video_id, name, adder_id from video where video_id in (select unnest(video_list) from users where user_id = $1 and document @@ plainto_tsquery($2)) order by created_at desc ) m;"
		if err := Postgres.Conn.QueryRow(context.Background(), query, userId, withVal).Scan(&count, &VideoStruct.UserVideos); err != nil {
			fmt.Println(err)
			return
		}
		if bytes.Equal(VideoStruct.UserVideos,null){
			VideoStruct.UserVideos = emptyArray
		}
	}

	limit:= 8
	if count < limit {
		query:= "select json_agg(m) from (select video_id, name, adder_id from video where document @@ plainto_tsquery($1) limit $2 offset $3) m; "
		if err := Postgres.Conn.QueryRow(context.Background(), query, withVal, limit - count, startFrom).Scan(&VideoStruct.AllVideos); err != nil {
			fmt.Println(err)
			return
		}
		if bytes.Equal(VideoStruct.AllVideos,null){
			VideoStruct.AllVideos = emptyArray
		}
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	jsonResult, _ := json.Marshal(VideoStruct)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

func PostVideoHandler(ctx *fasthttp.RequestCtx) {
	f, err := ctx.FormFile("video")
	adderId:= 1
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	videoId := 0
	query:= "insert into video (adder_id) values ($1) returning video_id"
	if err := Postgres.Conn.QueryRow(context.Background(), query, adderId).Scan(&videoId);
		err != nil {
		fmt.Println(err)
		return
	}

	path := functools.PathFromIdGenerator(strconv.Itoa(videoId))

	sb := strings.Builder{}
	sb.WriteString("../video_storage")
	sb.WriteString(path)

	if err := os.MkdirAll(sb.String(), 0777); err != nil {
		fmt.Println(err)
		return
	}

	sb.WriteString("/video.mp4")

	if err := fasthttp.SaveMultipartFile(f, sb.String()); err != nil {
		fmt.Println(err)
		return
	}

	query = "update users set video_list = array_append(video_list, $1) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, videoId, adderId);
		err != nil {
		fmt.Println(err)
		return
	}

	_, _ = ctx.WriteString(strconv.Itoa(videoId))

}