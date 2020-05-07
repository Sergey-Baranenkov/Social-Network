package main

import (
	"context"
	"coursework/functools"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

func SubscribeHandler(ctx *fasthttp.RequestCtx){
	subscriberId := 1
	subscribedId := 2

	query := "insert into relations__subscribers (subscriber, subscribed) values ($1, $2)"
	if _, err := Postgres.Conn.Exec(context.Background(), query, subscriberId, subscribedId);
		err != nil {
		ctx.Error("Вы уже подписаны!", 400)
		return
	}
	ctx.SetStatusCode(200)
}

func UnsubscribeHandler(ctx *fasthttp.RequestCtx){
	subscriberId := 1
	subscribedId := 2

	query := "delete from relations__subscribers where subscriber = $1 and subscribed = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, subscriberId, subscribedId);
		err != nil {
		ctx.Error("Вы не подписаны!", 400)
		return
	}
	ctx.SetStatusCode(200)
}

func AddSubscriberToFriendHandler(ctx *fasthttp.RequestCtx){
	subscriberId := 1
	subscribedId := 2
	query := "select add_subscriber_to_friend($1, $2)"
	if _, err := Postgres.Conn.Exec(context.Background(), query, subscribedId, subscriberId);
		err != nil {
		fmt.Println(err)
		ctx.Error("Пользователя нет в подписчиках!", 400)
		return
	}
	ctx.SetStatusCode(200)
}

func AddFriendToSubscriberHandler(ctx *fasthttp.RequestCtx){
	subscriberId := 1
	subscribedId := 2
	query := "select add_friend_to_subscriber($1, $2)"

	if _, err := Postgres.Conn.Exec(context.Background(), query, subscribedId, subscriberId);
		err != nil {
			fmt.Println(err)
		ctx.Error("Пользователя нет в друзьях!", 400)
		return
	}
	ctx.SetStatusCode(200)
}

type RelationshipsInfo struct {
	UserId int
	FirstName  string
	LastName  string
	AvatarRef string
}

func GetRelationshipsHandler(ctx *fasthttp.RequestCtx)  {
	userId := functools.ByteSliceToString(ctx.QueryArgs().Peek("userId"))
	limit := functools.ByteSliceToString(ctx.QueryArgs().Peek("limit"))
	mode := functools.ByteSliceToString(ctx.QueryArgs().Peek("mode"))

	var query string
	switch mode {
	case "friends":
		query = `
				with friends as (
					select user_id1 from relations__friends where user_id2 = $1
					union
					select user_id2 from relations__friends where user_id1 = $1
				) select user_id1, first_name, last_name, avatar_ref from friends f inner join users on user_id = f.user_id1 limit $2;
				`
	case "subscribers":
		query = `
				with subscribers as (
					select subscriber_id from relations__subscribers where subscribed_id = $1
				) select subscriber_id, first_name, last_name, avatar_ref from subscribers s inner join users on user_id = s.subscriber_id limit $2;
				`
	case "subscribed":
		query = `
				with subscribed as (
					select subscribed_id from relations__subscribers where subscriber_id = 1
				) select subscribed_id, first_name, last_name, avatar_ref from subscribed s inner join users on user_id = s.subscribed_id;
				`
	default:
		ctx.Error("Не указан тип relationships", 400)
		return
	}
	result := &RelationshipsInfo{}

	if err := Postgres.Conn.QueryRow(context.Background(), query, userId, limit).Scan(
																						&result.UserId,
																						&result.FirstName,
																						&result.LastName,
																						&result.AvatarRef);
		err != nil {
		_, _ = ctx.WriteString("[{}]")
		return
	}
	outputJson, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(outputJson))
}