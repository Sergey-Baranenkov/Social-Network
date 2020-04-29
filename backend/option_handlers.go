package main

import (
	"coursework/functools"
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
	Status   uint   `json:"status" validate:"required,min=0,max=5"`
	Birthday string `json:"birthday"`
	Tel      uint   `json:"tel"`
	Country  string `json:"country"`
	City     string `json:"city"`
}

func UpdateBasicInfoTextHandler(ctx *fasthttp.RequestCtx) {
	obj := &BasicInfoStruct{}

	if err := json.Unmarshal(ctx.PostBody(), obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}

	if err := Validator.Struct(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
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

func HobbiesHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = ctx.Write([]byte(`{"hobby":"Хобби",
							  "music": "Паша техник",
						      "films":"Интерстеллар", 
							  "books":"Преступление и наказание", 
 							  "games": "TF",
							  "others":""
							}`))
}

func PrivacyHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = ctx.Write([]byte(`{"can_m":"Все",
							  "has_access": "Только друзья",
						      "sound_n": false
							}`))
}

func EduEmpHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	_, _ = ctx.Write([]byte(`{"data": [
										{"title": "a", "period": "b", "description": "c"},
										{"title": "d", "period": "e", "description": "f"}
]}`))
}
