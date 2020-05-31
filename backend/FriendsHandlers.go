package main

import (
	"context"
	"coursework/functools"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)


func UpdateRelationshipHandler(ctx *fasthttp.RequestCtx) {
	requesterId := ctx.UserValue("requestUserId").(int)
	relationWith := functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	prevRelationType := functools.ByteSliceToString(ctx.QueryArgs().Peek("prevRelType"))

	fmt.Println(prevRelationType)
	var query string
	switch prevRelationType {
	case "0": // subscribe
		query = "insert into relations__subscribers (subscriber_id, subscribed_id) values ($1, $2)"
	case "1": // unsubscribe
		query = "delete from relations__subscribers where subscriber_id = $1 and subscribed_id = $2"
	case "2": // subscriber to friend
		query = "select add_subscriber_to_friend($1, $2)"
	case "3":
		query = "select add_friend_to_subscriber($1, $2)"
	}

	if _, err := Postgres.Conn.Exec(context.Background(), query, requesterId, relationWith);
		err != nil {
			fmt.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetStatusCode(200)
}


func GetRelationshipsHandler(ctx *fasthttp.RequestCtx)  {
	userId := functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))
	mode := functools.ByteSliceToString(ctx.QueryArgs().Peek("mode"))
	fmt.Println(userId, limit, mode)
	var query string
	switch mode {
	case "3":
		query = `
				with friends as (
					select (case when user_id1 = $1 then user_id2 else user_id1 end) as uid
					from relations__friends where user_id1 = $1 or user_id2 = $1
				) select json_agg(row) from (select user_id, first_name, last_name from friends f inner join users on user_id = f.uid limit $2) row;
				`
	case "2":
		query = `
				with subscribers as (
					select subscriber_id from relations__subscribers where subscribed_id = $1
				) select json_agg(row) from (select user_id, first_name, last_name from subscribers s inner join users on user_id = s.subscriber_id limit $2) row;
				`
	case "1":
		query = `
				with subscribed as (
					select subscribed_id from relations__subscribers where subscriber_id = $1
				) select json_agg(row) from (select user_id, first_name, last_name from subscribed s inner join users on user_id = s.subscribed_id limit $2) row;
				`
	default:
		ctx.Error("Не указан тип relationships", 400)
		return
	}

	var result json.RawMessage

	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, limit).Scan(&result);
		err != nil {
			ctx.SetStatusCode(400)
		return
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}