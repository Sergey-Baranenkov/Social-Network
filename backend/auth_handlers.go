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
)

/*func authPageHandler(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("frontend/html/auth.html")
}*/

func AuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		/*redisAccessToken, err := Redis.Get(functools.ByteSliceToString(ctx.Request.Header.Cookie("userId"))).Result()
		if err == nil && bytes.Compare(ctx.Request.Header.Cookie("accessToken"),
			functools.StringToByteSlice(redisAccessToken)) == 0{
			ctx.SetUserValue("a","b")
			next(ctx)
		}else{
			ctx.Redirect("/auth",401)
		}
		*/
		ctx.SetUserValue("a", "b")
		next(ctx)
	}
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

	if userToken := sha512.Sum512(append(functools.StringToByteSlice(obj.Password), Salt...));
	!bytes.Equal(dbToken, userToken[:]) {
		ctx.Error("Incorrect email/pass combination", 402)
		return
	}
	successfulAuth(ctx, strconv.Itoa(userId))
	ctx.Redirect("/posts", 200)

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
		successfulAuth(ctx, strconv.Itoa(userId))
		ctx.Redirect("/posts", 200)
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
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	return &authCookie
}

func successfulAuth(ctx *fasthttp.RequestCtx, userId string) {
	var accessToken string
	for {
		accessToken = functools.RandomStringGenerator(128)
		if _, err := Redis.Get(accessToken).Result(); err != nil {
			break
		}
	}

	accessTokenCookie := CreateCookie("accessToken", accessToken, 36000000000)
	idCookie := CreateCookie("userId", userId, 36000000000)
	Redis.Set(userId, accessToken, 360000000000)
	ctx.Response.Header.SetCookie(accessTokenCookie)
	ctx.Response.Header.SetCookie(idCookie)
}
