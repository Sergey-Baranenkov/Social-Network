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
	offset:= functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))

	AudioStruct := AudioJson{emptyArray,emptyArray, false}

	query:= `select json_agg(m) from 
                    (select m.music_id, m.adder_id, m.name, m.author from 
                        (select music_id, ordinality from users, unnest(music_list) with ordinality music_id 
                             where user_id = $1 limit $2 offset $3) as uml
                        inner join music m on m.music_id = uml.music_id order by uml.ordinality
                    ) m ;`
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, limit, offset).Scan(&AudioStruct.UserMusic);
		err != nil {
		fmt.Println(err)
		return
	}

	if bytes.Equal(AudioStruct.UserMusic,null){
		AudioStruct.Done = true
		AudioStruct.UserMusic = emptyArray
	}

	jsonResult, _ := json.Marshal(AudioStruct)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

func GetCombinedMusicHandler(ctx *fasthttp.RequestCtx){
	fmt.Println("vovovoo")
	ctx.Response.Header.Set("Content-Type", "application/json")
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	offset:= functools.ByteSliceToString(ctx.QueryArgs().Peek("offset"))
	withVal := functools.ByteSliceToString(ctx.QueryArgs().Peek("withValue"))
	limit, err := strconv.Atoi(functools.ByteSliceToString(ctx.QueryArgs().Peek("limit")))

	if err!=nil{
		ctx.SetStatusCode(400)
		return
	}

	AudioStruct := AudioJson{emptyArray,emptyArray, false}
	count:= 0

	if offset == "0" {
		query:= `select count(*), json_agg(m) from
                    (select m.music_id, m.adder_id, m.name, m.author from
                        (select music_id, ordinality from users, unnest(music_list) with ordinality music_id where user_id = $1) as uml
                        inner join music m on m.music_id = uml.music_id where document @@ plainto_tsquery($2) order by uml.ordinality
                    ) m ;
				`

		if err := Postgres.Conn.QueryRow(context.Background(), query, userId, withVal).Scan(&count, &AudioStruct.UserMusic); err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(400)
			return
		}
		if bytes.Equal(AudioStruct.UserMusic, null){
			AudioStruct.UserMusic = emptyArray
		}
	}


	if count < limit {
		query:= `select json_agg(m) from (select music_id, name, author, adder_id from music 
				 where document @@ plainto_tsquery($1) limit $2 offset $3) m;`

		if err := Postgres.Conn.QueryRow(context.Background(), query, withVal, limit - count, offset).Scan(&AudioStruct.AllMusic); err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(400)
			return
		}
		if bytes.Equal(AudioStruct.AllMusic,null){
			AudioStruct.Done = true
			AudioStruct.AllMusic = emptyArray
		}
	}

	jsonResult, _ := json.Marshal(AudioStruct)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

func PostMusicHandler(ctx *fasthttp.RequestCtx) {
	f, err := ctx.FormFile("audio")
	adderId:= 1
	author:= functools.ByteSliceToString(ctx.QueryArgs().Peek("author"))
	title:= functools.ByteSliceToString(ctx.QueryArgs().Peek("title"))

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	musicId := 0
	query:= "insert into music (adder_id, name, author) values ($1,$2,$3) returning music_id"
	if err := Postgres.Conn.QueryRow(context.Background(), query, adderId, title, author).Scan(&musicId);
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

	query = "update users set music_list = array_prepend($1, music_list) where user_id = $2"
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

func AddMusicToPlayList(ctx * fasthttp.RequestCtx){
	userId := 1
	musicId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("musicId"))
	fmt.Println(musicId)
	query := "update users set music_list = array_prepend($1, music_list) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, musicId, userId);
		err != nil {
			fmt.Println(err)
		ctx.Error("нет такого music id", 400)
		return
	}
	ctx.SetStatusCode(200)
}