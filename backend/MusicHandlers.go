package main

import (
	"context"
	"coursework/functools"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetMusicHandler(ctx *fasthttp.RequestCtx) {
	request:= ctx.QueryArgs().Peek("request")
	ctx.Response.Header.Set("Content-Type", "application/json")
	music := make([]byte, 0, 1024)
	query:= "select json_agg(m) from (select music_id, name, author from music where document @@ plainto_tsquery($1)) m"
	if err := Postgres.Conn.QueryRow(context.Background(), query, request).Scan(&music);
		err != nil {
		fmt.Println(err)
		return
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(music))
}

func PostMusicHandler(ctx *fasthttp.RequestCtx) {
	f, err := ctx.FormFile("audio")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	path, err := functools.StringToPath(functools.RandomStringGenerator(16), 2)
	if err != nil {
		log.Fatal(err)
	}

	sb := strings.Builder{}
	sb.WriteString("./music/")
	sb.WriteString(path)

	if err := os.MkdirAll(sb.String(), 0777); err != nil {
		fmt.Println(err)
		return
	}

	sb.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
	sb.WriteString(".mp3")

	if err := fasthttp.SaveMultipartFile(f, sb.String()); err != nil {
		fmt.Println(err)
		return
	}
}