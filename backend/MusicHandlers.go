package main

import (
	"context"
	"coursework/functools"
	"fmt"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
	"strings"
)

func GetUserMusicHandler(ctx *fasthttp.RequestCtx) {
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))

	ctx.Response.Header.Set("Content-Type", "application/json")
	music := make([]byte, 0, 1024)
	query:= `select json_agg(m) from (select music_id, name, author, adder_id from music where
				music_id in (select unnest(music_list) from users where user_id = $1)) m`

	if err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&music);
		err != nil {
		fmt.Println(err)
		return
	}

	_, _ = ctx.WriteString(functools.ByteSliceToString(music))
}

func GetAllMusicHandler(ctx *fasthttp.RequestCtx){
	request:= ctx.QueryArgs().Peek("request")
	limit:= functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))

	ctx.Response.Header.Set("Content-Type", "application/json")
	music := make([]byte, 0, 1024)
	query:= "select json_agg(m) from (select music_id, name, author, adder_id from music where document @@ plainto_tsquery($1) limit $2) m"
	if err := Postgres.Conn.QueryRow(context.Background(), query, request, limit).Scan(&music);
		err != nil {
		fmt.Println(err)
		return
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(music))
}

func PostMusicHandler(ctx *fasthttp.RequestCtx) {
	f, err := ctx.FormFile("audio")

	adderId:= 1

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	musicId:=0
	query:= "insert into music (adder_id) values ($1) returning music_id"
	if err := Postgres.Conn.QueryRow(context.Background(), query, adderId).Scan(&musicId);
		err != nil {
		fmt.Println(err)
		return
	}

	path := functools.PathFromIdGenerator(strconv.Itoa(musicId))

	sb := strings.Builder{}
	sb.WriteString("./music")
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

	_, _ = ctx.WriteString(path)

}