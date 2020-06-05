package main

import (
	"bytes"
	"context"
	"coursework/functools"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"time"
)

func staticPageHandler(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("../frontend/index.html")
}

func AuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		userId := functools.ByteSliceToString(ctx.Request.Header.Cookie("userId"))
		numericUserId, err := strconv.Atoi(userId)
		if err != nil{
			fmt.Println("num parsing error")
			ctx.Response.Header.SetCookie(DeleteCookie("firstName"))
			ctx.Response.Header.SetCookie(DeleteCookie("userId"))
			ctx.Response.Header.SetCookie(DeleteCookie("lastName"))
			ctx.Response.Header.SetCookie(DeleteCookie("accessToken"))
			ctx.Redirect("/авторизация",307)
			return
		}
		ctx.SetUserValue("requestUserId", numericUserId)

		redisUserId, err := Redis.Get(functools.ByteSliceToString(ctx.Request.Header.Cookie("accessToken"))).Result()
		fmt.Println(err, redisUserId)

		if err == nil && redisUserId == userId{
			fmt.Println("без ошибок")
			next(ctx)
		}else{
			fmt.Println("ошибка нет в реедисе", redisUserId, userId, ctx.UserValue("requestUserId"))
			ctx.Response.Header.SetCookie(DeleteCookie("firstName"))
			ctx.Response.Header.SetCookie(DeleteCookie("userId"))
			ctx.Response.Header.SetCookie(DeleteCookie("lastName"))
			ctx.Response.Header.SetCookie(DeleteCookie("accessToken"))
			ctx.Redirect("/авторизация",307)
			return
		}
	}
}


func DeleteCookie(key string) *fasthttp.Cookie {
	c := fasthttp.Cookie{}
	c.SetKey(key)
	c.SetValue("")
	c.SetExpire(fasthttp.CookieExpireDelete)
	c.SetPath("/")
	return &c
}


type loginStruct struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func loginHandler(ctx *fasthttp.RequestCtx) {
	obj := &loginStruct{}

	if err := json.Unmarshal(ctx.PostBody(), obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	err := Validator.Struct(obj)

	if err != nil {
		ctx.Error("not validated", 403)
		return
	}

	var dbToken []byte
	var userId int
	var firstName string
	var lastName string

	_ = Postgres.Conn.QueryRow(context.Background(),
		"select user_id, first_name, last_name, token from users where email = $1 limit 1", obj.Email).Scan(
		&userId,
		&firstName,
		&lastName,
		&dbToken)

	if userToken := sha512.Sum512(append(functools.StringToByteSlice(obj.Password), Salt...)); !bytes.Equal(dbToken, userToken[:]) {
		ctx.Error("Incorrect email/pass combination", 402)
		return
	}
	successfulAuth(ctx, strconv.Itoa(userId), firstName, lastName)
	ctx.Redirect("/сообщения", 307)
}

type RegistrationStruct struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Sex       string `json:"sex" validate:"sex"`
	Password  string `json:"password" validate:"required"`
}

func RegistrationHandler(ctx *fasthttp.RequestCtx) {
	obj := &RegistrationStruct{}
	if err := json.Unmarshal(ctx.PostBody(), obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	if err := Postgres.Conn.QueryRow(context.Background(), "select 1 from users where email = $1 limit 1",
		obj.Email).Scan(); err != pgx.ErrNoRows {
		ctx.Error("User already exists", 402)
		return
	}

	token := sha512.Sum512(append(functools.StringToByteSlice(obj.Password), Salt...))

	var userId int
	if err := Postgres.Conn.QueryRow(context.Background(),
		"insert into users (first_name,last_name,email,sex, token) values($1,$2,$3,$4,$5) returning user_id",
		obj.FirstName,
		obj.LastName,
		obj.Email,
		obj.Sex,
		token[:]).Scan(&userId); err != nil {
		fmt.Println(err)
	} else {
		successfulAuth(ctx, strconv.Itoa(userId), obj.FirstName, obj.LastName)
		ctx.Redirect("/сообщения", 307)
	}
}

func CreateCookie(key string, value string, expire int) *fasthttp.Cookie {
	if strings.Compare(key, "") == 0 {
		key = "unhandled cookie"
	}
	authCookie := fasthttp.Cookie{}
	authCookie.SetKey(key)
	authCookie.SetValue(value)
	authCookie.SetMaxAge(expire)
	authCookie.SetPath("/")
	return &authCookie
}

func successfulAuth(ctx *fasthttp.RequestCtx, userId string, firstName string, lastName string) {
	var accessToken string
	for {
		accessToken = functools.RandomStringGenerator(128)
		if res, _ := Redis.Exists(accessToken).Result(); res == 0 {
			break
		}
	}

	accessTokenCookie := CreateCookie("accessToken", accessToken, 36000000)
	idCookie := CreateCookie("userId", userId, 36000000)
	firstNameCookie := CreateCookie("firstName", firstName, 36000000)
	lastNameCookie := CreateCookie("lastName", lastName, 36000000)
	Redis.Set(accessToken, userId, time.Hour * 24 * 365)

	ctx.Response.Header.SetCookie(accessTokenCookie)
	ctx.Response.Header.SetCookie(idCookie)
	ctx.Response.Header.SetCookie(firstNameCookie)
	ctx.Response.Header.SetCookie(lastNameCookie)
}
