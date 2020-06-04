package main

import (
	"bytes"
	"context"
	"coursework/functools"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/valyala/fasthttp"
	"image/jpeg"
	"os"
	"strconv"
	"strings"
)

type BasicInfoStruct struct {
	Sex      string `json:"sex" validate:"sex"`
	Status   uint   `json:"status" validate:"min=0,max=5"`
	Birthday string `json:"birthday"`
	Tel      uint   `json:"tel"`
	Country  string `json:"country"`
	City     string `json:"city"`
}

func UpdateBasicInfoTextHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	bis := BasicInfoStruct{}

	if err := json.Unmarshal(ctx.PostBody(), &bis); err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	if err := Validator.Struct(bis); err != nil {
		fmt.Println(err)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	query := "update users set sex = $1, status = $2, birthday = $3, tel = $4, " +
		"country = $5, city = $6 where user_id = $7"
	if _, err := Postgres.Conn.Exec(
		context.Background(),
		query,
		bis.Sex,
		bis.Status,
		bis.Birthday,
		bis.Tel,
		bis.Country,
		bis.City,
		userId); err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetStatusCode(200)
}

func UpdateProfileAvatar(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)

	f, err := ctx.FormFile("photo")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	path := functools.PathFromIdGenerator(strconv.Itoa(userId))

	sb := strings.Builder{}
	sb.WriteString("../profile_bgs")
	sb.WriteString(path)
	if err := os.MkdirAll(sb.String(), 0777); err != nil {
		fmt.Println(err)
		return
	}
	sb.WriteString("/profile_avatar.jpg")
	if err := fasthttp.SaveMultipartFile(f, sb.String()); err != nil {
		fmt.Println(err)
		return
	}

	ff, err := f.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	decodedJpeg, err := jpeg.Decode(ff)
	if err != nil {
		fmt.Println(err)
		return
	}
	avatar200x200Pic := resize.Resize(200, 200, decodedJpeg, resize.Lanczos3)

	w, err := os.Create(sb.String())

	if err != nil {
		fmt.Println(err)
		return
	}

	if err := jpeg.Encode(w, avatar200x200Pic, nil); err != nil {
		fmt.Println(err)
		return
	}

}

func UpdateProfileBg(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)

	f, err := ctx.FormFile("photo")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	path := functools.PathFromIdGenerator(strconv.Itoa(userId))

	sb := strings.Builder{}
	sb.WriteString("../profile_bgs")
	sb.WriteString(path)
	if err := os.MkdirAll(sb.String(), 0777); err != nil {
		fmt.Println(err)
		return
	}
	sb.WriteString("/profile_bg.jpg")
	if err := fasthttp.SaveMultipartFile(f, sb.String()); err != nil {
		fmt.Println(err)
		return
	}
}

func GetEduEmpHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	query := "select edu_and_emp_info from users where user_id = $1"
	var result json.RawMessage
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result); err != nil {
		ctx.Error("", 400)
		return
	}
	if bytes.Equal(result, null) {
		result = emptyArray
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}

func GetHobbiesHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	query := "select to_json(s) from (select hobby, fav_music, fav_films, fav_books, fav_games, other_interests from users where user_id = $1) s"
	var result json.RawMessage
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result); err != nil {
		ctx.Error("", 400)
		return
	}
	if bytes.Equal(result, null) {
		result = emptyArray
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}

func UpdateEduEmpHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	reqJson := ctx.Request.Body()

	query := "update users set edu_and_emp_info = $1 where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, reqJson, userId); err != nil {
		ctx.Error("", 400)
		return
	}
	ctx.SetStatusCode(400)
}

type UpdatePasswordStruct struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" validate:"min=10,required"`
}

func UpdatePasswordHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	reqJson := ctx.Request.Body()

	ups := UpdatePasswordStruct{}

	if err := json.Unmarshal(reqJson, &ups); err != nil {
		ctx.SetStatusCode(400)
		return
	}

	err := Validator.Struct(ups)
	if err != nil {
		fmt.Println("not valid struct")
		ctx.Error("not valid", 403)
		return
	}

	var oldToken []byte
	query := "select token from users where user_id = $1"
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&oldToken); err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(400)
		return
	}

	usersOldToken := sha512.Sum512(append(functools.StringToByteSlice(ups.OldPassword), Salt...))

	if !bytes.Equal(usersOldToken[:], oldToken) {
		fmt.Println("Старый пароль не верный")
		ctx.Error("Старый пароль не верный", 400)
		return
	}

	newToken := sha512.Sum512(append(functools.StringToByteSlice(ups.NewPassword), Salt...))
	query = "update users set token = $1 where user_id = $2"
	if _, err := Postgres.Conn.Exec(context.Background(), query, ups.NewPassword, newToken); err != nil {
		ctx.Error("Unhandled error", 400)
		return
	}
	ctx.SetStatusCode(200)
}

type UpdateHobbiesStruct struct {
	Hobby          string `json:"hobby"`
	FavMusic       string `json:"fav_music"`
	FavFilms       string `json:"fav_films"`
	FavBooks       string `json:"fav_books"`
	FavGames       string `json:"fav_games"`
	OtherInterests string `json:"other_interests"`
}

func UpdateHobbiesHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	reqJson := ctx.Request.Body()

	uhs := UpdateHobbiesStruct{}

	if err := json.Unmarshal(reqJson, &uhs); err != nil {
		ctx.SetStatusCode(400)
		return
	}

	query := "update users set hobby = $1, fav_music = $2, fav_films = $3, " +
		"fav_books = $4, fav_games = $5, other_interests = $6 where user_id = $7"
	if _, err := Postgres.Conn.Exec(
		context.Background(),
		query,
		uhs.Hobby,
		uhs.FavMusic,
		uhs.FavFilms,
		uhs.FavBooks,
		uhs.FavGames,
		uhs.OtherInterests,
		userId); err != nil {
		ctx.Error("", 400)
		return
	}
	ctx.SetStatusCode(400)
}

func GetBasicInfoHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	query := "select to_json(s) from (select sex, status, birthday, tel, country, city from users where user_id = $1) s "
	var result json.RawMessage
	if err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result); err != nil {
		fmt.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	if bytes.Equal(result, null) {
		result = emptyArray
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}
