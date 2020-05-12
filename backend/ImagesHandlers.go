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

type ImagesStruct struct {
	ImagesList json.RawMessage
	Done bool
}

func GetUserImages(ctx *fasthttp.RequestCtx){
	ctx.Response.Header.Set("Content-Type", "application/json")
	userId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	startFrom:= functools.ByteSliceToString(ctx.QueryArgs().Peek("startFrom"))
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))

	images:= ImagesStruct{emptyArray, false}

	query:= `select json_agg(i) from (select i.* from (select image_id, ordinality from  users, unnest(images_list) with ordinality image_id 
														where user_id = $1 limit $2 offset $3) as uil
    		 inner join images i on i.image_id = uil.image_id order by uil.ordinality) i;`
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, limit, startFrom).Scan(&images.ImagesList);
		err != nil {
		fmt.Println(err)
		return
	}

	if bytes.Equal(images.ImagesList,null){
		images.Done = true
		images.ImagesList = emptyArray
	}

	jsonResult, _ := json.Marshal(images)
	_, _ = ctx.WriteString(functools.ByteSliceToString(jsonResult))
}

func PostImageHandler(ctx *fasthttp.RequestCtx)  {
	f, err := ctx.FormFile("image")
	adderId:= 1

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	imageId:=0
	query:= "insert into images (adder_id) values ($1) returning image_id"
	if err := Postgres.Conn.QueryRow(context.Background(), query, adderId).Scan(&imageId);
		err != nil {
		fmt.Println(err)
		return
	}

	path := functools.PathFromIdGenerator(strconv.Itoa(imageId))

	sb := strings.Builder{}
	sb.WriteString("../gallery_storage")
	sb.WriteString(path)

	if err := os.MkdirAll(sb.String(), 0777); err != nil {
		fmt.Println(err)
		return
	}

	sb.WriteString("/img.jpg")

	if err := fasthttp.SaveMultipartFile(f, sb.String()); err != nil {
		fmt.Println(err)
		return
	}

	query = "update users set images_list = array_prepend($1, images_list) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, imageId, adderId);
		err != nil {
		fmt.Println(err)
		return
	}

	_, _ = ctx.WriteString(strconv.Itoa(imageId))
}

func DeleteImageHandler (ctx *fasthttp.RequestCtx){
	userId := 1
	imageId:= functools.ByteSliceToString(ctx.QueryArgs().Peek("imageId"))

	query := "update users set images_list = array_remove(images_list, $1) where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, imageId, userId);
		err != nil {
		ctx.Error("нет такого image_id", 400)
		return
	}
	ctx.SetStatusCode(200)
}