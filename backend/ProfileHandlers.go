package main

import (
	"context"
	"coursework/functools"
	"fmt"
	"github.com/valyala/fasthttp"
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
	authId := "1"

	if err := Postgres.Conn.QueryRow(context.Background(),
		"insert into objects (auth_id, path, text) values ($1, $2, $3) returning path", authId, path, text).Scan(&path);
	err != nil {
		fmt.Println("Error:", err)
		ctx.SetStatusCode(400)
		return
	}
	_, _ = ctx.WriteString(path)
}

func CommentsTestHandler(ctx *fasthttp.RequestCtx) {
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
	authId := "1" //fix
	path := functools.ByteSliceToString(ctx.QueryArgs().Peek("path"))
	option := functools.ByteSliceToString(ctx.QueryArgs().Peek("meLiked"))
	fmt.Println(path)
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
	userId := functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))

	var result string

	query := `select to_json(k) from 
				(select first_name, last_name, avatar_ref, bg_ref, tel, country, city, birthday, images_list[:9] 
					from users where user_id = $1) k;`

	if err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result);
		err != nil {
		fmt.Println(err)
		return
	}

	_, _ = ctx.WriteString(result)
}