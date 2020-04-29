package main

import (
	"context"
	"coursework/functools"
	"fmt"
	"github.com/valyala/fasthttp"
)

func GetPostsHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	var posts = make([]byte, 0, 1024)
	if err := Postgres.Conn.QueryRow(context.Background(), "select get_posts($1,$2)", 1, 1).Scan(&posts);
	err != nil {
		fmt.Println(err)
		return
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(posts))
}

func CommentsTestHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.QueryArgs().Peek("path")
	limit := ctx.QueryArgs().Peek("lim")
	fmt.Println(limit, path)
	ctx.Response.Header.Set("Content-Type", "application/json")
	var comments = make([]byte, 0, 1024)
	if err := Postgres.Conn.QueryRow(context.Background(),
		"select get_comments($1)", path).Scan(&comments); err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(comments))
}

var revokeLike = "0"
var setLike = "1"

func LikeHandler(ctx *fasthttp.RequestCtx) {
	authId := "1" //fix
	path := functools.ByteSliceToString(ctx.QueryArgs().Peek("path"))
	option := functools.ByteSliceToString(ctx.QueryArgs().Peek("opt"))
	if option == setLike {
		if _, err := Postgres.Conn.Exec(context.Background(),
			"insert into likes(path,auth_id) values($1,$2)", path, authId); err != nil {
			fmt.Println("Error:", err)
			return
		}
	} else if option == revokeLike {
		if _, err := Postgres.Conn.Exec(context.Background(),
			"delete from likes(path,auth_id) values($1,$2)", path, authId); err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

}

func AddCommentHandler(ctx *fasthttp.RequestCtx) {
	path := functools.ByteSliceToString(ctx.QueryArgs().Peek("path"))
	authId := "1"
	message := "Hello world"
	if _, err := Postgres.Conn.Exec(context.Background(),
		"insert into objects (auth_id, text, path) values ($1, $2, $3);", authId, message, path); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
