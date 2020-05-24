package main

import (
	"context"
	"coursework/functools"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
)

func GetPostsHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	var posts string
	if err := Postgres.Conn.QueryRow(context.Background(), "select get_posts($1, $2)", 1, 1).Scan(&posts);
	err != nil {
		fmt.Println(err)
		return
	}
	_, _ = ctx.WriteString(posts)
}

func AddNewObjectHandler(ctx * fasthttp.RequestCtx){
	path := functools.ByteSliceToString(ctx.QueryArgs().Peek("path"))
	text := functools.ByteSliceToString(ctx.QueryArgs().Peek("text"))
	authId := ctx.UserValue("requestUserId").(int)

	var result json.RawMessage
	if err := Postgres.Conn.QueryRow(context.Background(),
		"select push_object($1,$2,$3)", authId, path, text).Scan(&result);
	err != nil {
		fmt.Println("Error:", err)
		ctx.SetStatusCode(400)
		return
	}

	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}

func GetCommentsHandler(ctx *fasthttp.RequestCtx) {
	path := functools.ByteSliceToString(ctx.QueryArgs().Peek("path"))
	ctx.Response.Header.Set("Content-Type", "application/json")
	var comments = make([]byte, 0, 1024)
	if err := Postgres.Conn.QueryRow(context.Background(),
		"select get_comments($1)", path).Scan(&comments); err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(comments))
}

func UpdateLikeHandler(ctx *fasthttp.RequestCtx) {
	authId := ctx.UserValue("requestUserId").(int)
	path := functools.ByteSliceToString(ctx.QueryArgs().Peek("path"))
	option := functools.ByteSliceToString(ctx.QueryArgs().Peek("meLiked"))
	if option == "false" {
		if _, err := Postgres.Conn.Exec(context.Background(),
			"insert into likes(path,auth_id) values($1,$2)", path, authId); err != nil {
			fmt.Println("Error:", err)
			ctx.SetStatusCode(400)
			return
		}
	} else if option == "true" {
		if _, err := Postgres.Conn.Exec(context.Background(),
			"delete from likes where path = $1 and auth_id = $2", path, authId); err != nil {
			fmt.Println("Error:", err)
			ctx.SetStatusCode(400)
			return
		}
	}
	ctx.SetStatusCode(200)

}

func GetProfilePageInfo(ctx *fasthttp.RequestCtx){
	requesterId := ctx.UserValue("requestUserId").(int)
	userId := functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))

	var query string
	var err error
	var result string

	if strconv.Itoa(requesterId) == userId {
		query = `select to_json(k) from 
				(select first_name, last_name, avatar_ref, bg_ref, tel, country, city, birthday, images_list[:9]
					from users where user_id = $1) k;`
		err = Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result)
	}else{
		query = `select to_json(k) from 
				(select first_name, last_name, avatar_ref, bg_ref, tel, country, city, birthday, images_list[:9], get_relationship($2, $1) as rel
					from users where user_id = $1) k;`
		err = Postgres.Conn.QueryRow(context.Background(), query, userId, requesterId).Scan(&result)
	}

	if err != nil {
		ctx.SetStatusCode(400)
		fmt.Println(err)
		return
	}

	_, _ = ctx.WriteString(result)
}

type ObjectStruct struct {
	Text string
	Path string
}

func UpdateObjectText(ctx *fasthttp.RequestCtx){
	userId := ctx.UserValue("requestUserId").(int)

	os := ObjectStruct{}
	if err := json.Unmarshal(ctx.PostBody(), &os); err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	query := `update objects set text = $1 where path = $2 and auth_id = $3`
	if _, err := Postgres.Conn.Exec(context.Background(), query, os.Text, os.Path, userId);
		err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(400)
		return
	}

	ctx.SetStatusCode(200)
}

func DeleteObject(ctx *fasthttp.RequestCtx){
	userId := ctx.UserValue("requestUserId").(int)

	path := functools.ByteSliceToString(ctx.PostBody())
	fmt.Println(path)
	query := `delete from objects where path = $1 and auth_id = $2`
	if _, err := Postgres.Conn.Exec(context.Background(), query, path, userId);
		err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(400)
		return
	}

	ctx.SetStatusCode(200)
}