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

func GetUserVideoHandler(ctx *fasthttp.RequestCtx){
	ctx.Response.Header.Set("Content-Type", "application/json")

	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	offset:= functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))

	VideoStruct := VideoJSON{emptyArray,emptyArray, false}

	query:= `
			select json_agg(v) from 
                   (select v.video_id, v.adder_id, v.name from 
						(select video_id, ordinality from  users, unnest(video_list) with ordinality video_id 
							where user_id = $1 limit $2 offset $3) as uvl
    					inner join video v on v.video_id = uvl.video_id order by uvl.ordinality
					) v;`
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, limit, offset).Scan(&VideoStruct.UserVideos);
		err != nil {
		fmt.Println(err)
		return
	}
	if bytes.Equal(VideoStruct.UserVideos,null){
		VideoStruct.Done = true
		VideoStruct.UserVideos = emptyArray
	}

	jsonResult, _ := json.Marshal(VideoStruct)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

func GetCombinedVideoHandler(ctx *fasthttp.RequestCtx){
	ctx.Response.Header.Set("Content-Type", "application/json")
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	startFrom:= functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))
	withVal := functools.ByteSliceToString(ctx.QueryArgs().Peek("withValue"))
	limit, err := strconv.Atoi(functools.ByteSliceToString(ctx.QueryArgs().Peek("limit")))

	if err!=nil{
		ctx.SetStatusCode(400)
		return
	}

	VideoStruct := VideoJSON{emptyArray,emptyArray, false}
	count:= 0

	if startFrom == "0" {
		query:= `
				select count(*), json_agg(v) from
                    (select v.video_id, v.adder_id, v.name from
                        (select video_id, ordinality from users, unnest(video_list) with ordinality video_id where user_id = $1) as uvl
                        inner join video v on v.video_id = uvl.video_id where document @@ plainto_tsquery($2) order by uvl.ordinality
                    ) v ;
			    `
		if err := Postgres.Conn.QueryRow(context.Background(), query, userId, withVal).Scan(&count, &VideoStruct.UserVideos); err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(400)
			return
		}
		if bytes.Equal(VideoStruct.UserVideos,null){
			VideoStruct.UserVideos = emptyArray
		}
	}

	if count < limit {
		query:= "select json_agg(m) from (select video_id, name, adder_id from video " +
				"where document @@ plainto_tsquery($1) limit $2 offset $3) m; "
		if err := Postgres.Conn.QueryRow(context.Background(), query, withVal, limit - count, startFrom).Scan(&VideoStruct.AllVideos); err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(400)
			return
		}
		if bytes.Equal(VideoStruct.AllVideos, null){
			VideoStruct.Done = true
			VideoStruct.AllVideos = emptyArray
		}
	}

	jsonResult, _ := json.Marshal(VideoStruct)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

func PostVideoHandler(ctx *fasthttp.RequestCtx) {
	f, err := ctx.FormFile("video")
	adderId:= ctx.UserValue("requestUserId").(int)
	title:= functools.ByteSliceToString(ctx.QueryArgs().Peek("title"))
	fmt.Println(title)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	videoId := 0
	query:= "insert into video (adder_id, name) values ($1, $2) returning video_id"
	if err := Postgres.Conn.QueryRow(context.Background(), query, adderId, title).Scan(&videoId);
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

	query = "update users set video_list = array_prepend($1, video_list) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, videoId, adderId);
		err != nil {
		fmt.Println("err??", err)
		return
	}

	_, _ = ctx.WriteString(strconv.Itoa(videoId))

}

func DeleteVideoHandler (ctx *fasthttp.RequestCtx){
	userId := ctx.UserValue("requestUserId").(int)
	videoId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("videoId"))

	query := "update users set video_list = array_remove(video_list, $1) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, videoId, userId);
		err != nil {
		ctx.Error("нет такого video id", 400)
		return
	}
	ctx.SetStatusCode(200)
}

func AddVideoToPlayList(ctx * fasthttp.RequestCtx){
	userId := ctx.UserValue("requestUserId").(int)
	videoId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("videoId"))
	fmt.Println(videoId)
	query := "update users set video_list = array_prepend($1, video_list) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, videoId, userId);
		err != nil {
		fmt.Println(err)
		ctx.Error("нет такого video id", 400)
		return
	}
	ctx.SetStatusCode(200)
}



