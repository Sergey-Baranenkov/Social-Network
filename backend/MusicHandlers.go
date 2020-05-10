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

func GetUserMusicHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	startFrom:= functools.ByteSliceToString(ctx.QueryArgs().Peek("startFrom"))

	AudioStruct := AudioJson{json.RawMessage("[]"),json.RawMessage("[]")}
	query:= `select json_agg(m) from (select music_id, name, author, adder_id from music 
			 where music_id in (select unnest(music_list) from users where user_id = $1 limit 10 offset $2) order by created_at desc) m`
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, startFrom).Scan(&AudioStruct.UserMusic);
		err != nil {
		fmt.Println(err)
		return
	}

	if bytes.Equal(AudioStruct.UserMusic,null){
		AudioStruct.UserMusic = emptyArray
	}
	jsonResult, _ := json.Marshal(AudioStruct)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))

}

func GetCombinedMusicHandler(ctx *fasthttp.RequestCtx){
	ctx.Response.Header.Set("Content-Type", "application/json")
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	startFrom:= functools.ByteSliceToString(ctx.QueryArgs().Peek("startFrom"))
	withVal := functools.ByteSliceToString(ctx.QueryArgs().Peek("withVal"))
	AudioStruct := AudioJson{json.RawMessage("[]"),json.RawMessage("[]")}
	count:= 0

	if startFrom == "0" {
		query:= `select count(*), json_agg(m) from (select music_id, name, author, adder_id from music
													where music_id in (select unnest(music_list) from users 
																		where user_id = $1 and document @@ plainto_tsquery($2)) 
													order by created_at desc ) m;`

		if err := Postgres.Conn.QueryRow(context.Background(), query, userId, withVal).Scan(&count, &AudioStruct.UserMusic); err != nil {
			fmt.Println(err)
			return
		}
		if bytes.Equal(AudioStruct.UserMusic,null){
			AudioStruct.UserMusic = emptyArray
		}
	}

	limit:= 10
	if count < limit {
		query:= `select json_agg(m) from (select music_id, name, author, adder_id from music 
				 where document @@ plainto_tsquery($1) limit $2 offset $3) m;`

		if err := Postgres.Conn.QueryRow(context.Background(), query, withVal, limit - count, startFrom).Scan(&AudioStruct.AllMusic); err != nil {
			fmt.Println(err)
			return
		}
		if bytes.Equal(AudioStruct.AllMusic,null){
			AudioStruct.AllMusic = emptyArray
		}
	}

	jsonResult, _ := json.Marshal(AudioStruct)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

func PostMusicHandler(ctx *fasthttp.RequestCtx) {
	f, err := ctx.FormFile("audio")
	adderId:= 1

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	musicId := 0
	query:= "insert into music (adder_id) values ($1) returning music_id"
	if err := Postgres.Conn.QueryRow(context.Background(), query, adderId).Scan(&musicId);
		err != nil {
		fmt.Println(err)
		return
	}

	path := functools.PathFromIdGenerator(strconv.Itoa(musicId))

	sb := strings.Builder{}
	sb.WriteString("../music_storage")
	sb.WriteString(path)

	if err := os.MkdirAll(sb.String(), 0777); err != nil {
		fmt.Println(err)
		return
	}

	sb.WriteString("/audio.mp3")

	if err := fasthttp.SaveMultipartFile(f, sb.String()); err != nil {
		fmt.Println(err)
		return
	}

	query = "update users set music_list = array_append(music_list, $1) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, musicId, adderId);
		err != nil {
		fmt.Println(err)
		return
	}

	_, _ = ctx.WriteString(strconv.Itoa(musicId))

}

func DeleteMusicHandler (ctx *fasthttp.RequestCtx){
	userId := 1
	musicId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("musicId"))

	query := "update users set music_list = array_remove(music_list, $1) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, musicId, userId);
		err != nil {
		ctx.Error("нет такого music id", 400)
		return
	}
	ctx.SetStatusCode(200)
}