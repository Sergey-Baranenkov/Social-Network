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
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
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
	if  _, err := Postgres.Conn.Exec(
		context.Background(),
		query,
		bis.Sex,
		bis.Status,
		bis.Birthday,
		bis.Tel,
		bis.Country,
		bis.City,
		userId);
		err != nil {
			fmt.Println(err)
			ctx.SetStatusCode(400)
			return
	}
	ctx.SetStatusCode(200)
}

func UpdateProfileAvatar(ctx *fasthttp.RequestCtx) {
	f, err := ctx.FormFile("photo")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	path50, err := makeAvatarPath(50, 50, f)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}

	path150, err := makeAvatarPath(150, 150, f)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
	/*
		if _, err := Postgres.Conn.Exec(context.Background(), "erererer"); err!=nil{
			log.Fatal("Cannot update avatar of user")
		}
	*/
	fmt.Println(ctx.UserValue("a"), path50, path150) //TODO
}

func makeAvatarPath(width uint, height uint, file *multipart.FileHeader) (string, error) {
	path, err := functools.StringToPath(functools.RandomStringGenerator(16), 2)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	sb.WriteString("./profile_avatars/")
	sb.WriteString(strconv.Itoa(int(width)))
	sb.WriteString("x")
	sb.WriteString(strconv.Itoa(int(height)))
	sb.WriteString("/")
	sb.WriteString(path)

	if err := os.MkdirAll(sb.String(), 0777); err != nil {
		return "", err
	}

	sb.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
	sb.WriteString(".jpg")

	ff, err := file.Open()
	if err != nil {
		return "", err
	}

	decodedJpeg, err := jpeg.Decode(ff)
	if err != nil{
		return "", err
	}
	avatar100xPic := resize.Resize(width, height, decodedJpeg, resize.Lanczos3)

	w, err := os.Create(sb.String())

	if err != nil {
		return "", err
	}

	if err := jpeg.Encode(w, avatar100xPic, nil); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func UpdateProfileBg(ctx *fasthttp.RequestCtx) {
	f, err := ctx.FormFile("photo")
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	path, err := functools.StringToPath(functools.RandomStringGenerator(16), 2)
	if err != nil {
		log.Fatal(err)
		return
	}

	sb := strings.Builder{}
	sb.WriteString("./profile_bgs/")
	sb.WriteString(path)

	if err := os.MkdirAll(sb.String(), 0777); err != nil {
		fmt.Println(err)
		return
	}

	sb.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
	sb.WriteString(".jpg")

	if err := fasthttp.SaveMultipartFile(f, sb.String()); err != nil {
		fmt.Println(err)
		return
	}
}

func GetEduEmpHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	query := "select edu_and_emp_info from users where user_id = $1"
	var result json.RawMessage
	if  err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result);
		err != nil {
		ctx.Error("", 400)
		return
	}
	if bytes.Equal(result, null){
		result = emptyArray
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}

func GetHobbiesHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	query := "select to_json(s) from (select hobby, fav_music, fav_films, fav_books, fav_games, other_interests from users where user_id = $1) s"
	var result json.RawMessage
	if  err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result);
		err != nil {
		ctx.Error("", 400)
		return
	}
	if bytes.Equal(result, null){
		result = emptyArray
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}

func UpdateEduEmpHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	reqJson := ctx.Request.Body()

	query := "update users set edu_and_emp_info = $1 where user_id = $2"
	if  _, err := Postgres.Conn.Exec(context.Background(), query, reqJson, userId);
		err != nil {
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

	if err := json.Unmarshal(reqJson, &ups); err!= nil{
		ctx.SetStatusCode(400);
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
	if  err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&oldToken);
		err != nil {
			fmt.Println(err)
		ctx.SetStatusCode( 400)
		return
	}

	usersOldToken := sha512.Sum512(append(functools.StringToByteSlice(ups.OldPassword), Salt...))

	if !bytes.Equal(usersOldToken[:], oldToken){
		fmt.Println("Старый пароль не верный")
		ctx.Error("Старый пароль не верный", 400)
		return
	}

	newToken := sha512.Sum512(append(functools.StringToByteSlice(ups.NewPassword), Salt...))
	query = "update users set token = $1 where user_id = $2"
	if  _, err := Postgres.Conn.Exec(context.Background(), query, ups.NewPassword, newToken);
		err != nil {
		ctx.Error("Unhandled error", 400)
		return
	}
	ctx.SetStatusCode(200)
}


type UpdateHobbiesStruct struct {
	Hobby string `json:"hobby"`
	FavMusic string `json:"fav_music"`
	FavFilms string `json:"fav_films"`
	FavBooks string `json:"fav_books"`
	FavGames string `json:"fav_games"`
	OtherInterests string `json:"other_interests"`
}

func UpdateHobbiesHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	reqJson := ctx.Request.Body()

	uhs := UpdateHobbiesStruct{}

	if err := json.Unmarshal(reqJson, &uhs); err!= nil{
		ctx.SetStatusCode(400)
		return
	}

	query := "update users set hobby = $1, fav_music = $2, fav_films = $3, " +
			 "fav_books = $4, fav_games = $5, other_interests = $6 where user_id = $7"
	if  _, err := Postgres.Conn.Exec(
		context.Background(),
			query,
			uhs.Hobby,
			uhs.FavMusic,
			uhs.FavFilms,
			uhs.FavBooks,
			uhs.FavGames,
			uhs.OtherInterests,
			userId);
		err != nil {
			ctx.Error("", 400)
			return
	}
	ctx.SetStatusCode(400)
}

func GetBasicInfoHandler(ctx *fasthttp.RequestCtx) {
	userId := ctx.UserValue("requestUserId").(int)
	query := "select to_json(s) from (select sex, status, birthday, tel, country, city from users where user_id = $1) s "
	var result json.RawMessage
	if  err := Postgres.Conn.QueryRow(context.Background(), query, userId).Scan(&result);
		err != nil {
			fmt.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	if bytes.Equal(result, null){
		result = emptyArray
	}
	_, _ = ctx.WriteString(functools.ByteSliceToString(result))
}