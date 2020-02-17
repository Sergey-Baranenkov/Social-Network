package main

import (
	"bytes"
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/go-redis/redis/v7"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"html/template"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
	"coursework/postgres"
)



var (
	rc = postgres.RegistrationConn{}
	authHTML *template.Template
	rdb *redis.Client
	err error
	salt = []byte("Ilya Bychkov")
	r  = router.New()
)

func main() {
	if err:= rc.CreateConnection(os.Getenv("REGISTER_DB"));err!=nil{
		log.Fatal(err)
	}

	if err := rc.CreateRegTable();err!=nil{
		log.Fatal(err)
	}

	if authHTML, err = template.ParseFiles("frontend/auth.html"); err != nil{
		log.Fatal(err)
	}

       
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	if err = rdb.Ping().Err();err!=nil{
		log.Fatal(err)
		return
	}

	r.GET("/frontend/auth.jsx", authReactHandler)
	r.GET("/auth",authPageHandler)
	r.POST("/registration",registrationHandler)
	r.POST("/login",loginHandler)
	r.GET("/frontend/css/*filepath",fasthttp.FSHandler("./frontend/css",2))
	r.GET("/frontend/images/*filepath",fasthttp.FSHandler("./frontend/images",2))
	r.GET("/secretpage",AccessMiddleware(secretPageHandler))
	r.NotFound = redirectHandler
        
        clientPort:="8080"
        fmt.Println("LISTENING ON PORT " + clientPort)
	if err:= fasthttp.ListenAndServe(":" + clientPort, r.Handler); err!= nil{
		log.Println("error when starting server: " + err.Error())
	}
	if err := rc.Close(); err!=nil{
		log.Println("error when closing regDb conn: " + err.Error())
	}
	if err := rdb.Close(); err!=nil{
		log.Println("error when closing regRedis conn: " + err.Error())
	}
}

func redirectHandler(ctx *fasthttp.RequestCtx){
	ctx.Redirect("/auth",2020)
}

func secretPageHandler(ctx *fasthttp.RequestCtx){
	fmt.Fprint(ctx,"Hello"+ByteSliceToString(ctx.Request.Header.Cookie("userId"))+"!")
}
func authPageHandler(ctx *fasthttp.RequestCtx){
		ctx.SetContentType("text/html")
		authHTML.Execute(ctx,authHTML)
}

func AccessMiddleware(next fasthttp.RequestHandler)fasthttp.RequestHandler{
	return func(ctx *fasthttp.RequestCtx){
		redisAccessToken, err := rdb.Get(ByteSliceToString(ctx.Request.Header.Cookie("userId"))).Result()
		if err == nil && bytes.Compare(ctx.Request.Header.Cookie("accessToken"),StringToByteSlice(redisAccessToken)) == 0{
			next(ctx)
		}else{
			fmt.Fprint(ctx,"U have no permission there")
		}
	}
}


func authReactHandler(ctx *fasthttp.RequestCtx){
	ctx.SendFile("frontend/auth.jsx")
}

func loginHandler(ctx *fasthttp.RequestCtx){
	//тут должна быть валидация
	email:= ctx.FormValue("email")
	password:= ctx.FormValue("password")
	if len(email)==0 || len(password)==0{
		ctx.Error("Поля не заполнены",402)
	}

	var dbToken []byte
	var userId int
	var firstName string
	var lastName string
	rc.Conn.QueryRow(context.Background(),"select user_id, first_name, last_name, token from registration where email = $1 limit 1", email).Scan(
		&userId,
		&firstName,
		&lastName,
		&dbToken)
	if userToken:=sha512.Sum512(append(password,salt...)); bytes.Compare(dbToken,userToken[:]) != 0{
		ctx.Error("Incorrect email/pass combination",402)
		return
	}
	successfulAuth(ctx,strconv.Itoa(userId))
	ctx.Redirect("/secretpage",200)

}

func registrationHandler(ctx *fasthttp.RequestCtx){
	firstName:=ctx.FormValue("first_name")
	lastName:=ctx.FormValue("last_name")
	email:= ctx.FormValue("email")
	password:= ctx.FormValue("password")

	if len(email)==0 || len(password)==0{
		ctx.Error("Поля не заполнены",402)
	}

	if err:= rc.Conn.QueryRow(context.Background(),"select 1 from registration where email = $1 limit 1", email).Scan();err != pgx.ErrNoRows{
		ctx.Error("User already exists",402)
	}

	token:=sha512.Sum512(append(password,salt...))
	var userId int
	if err := rc.Conn.QueryRow(context.Background(), "insert into registration (first_name,last_name,email,token) values($1,$2,$3,$4) returning user_id",
		ByteSliceToString(firstName),
		ByteSliceToString(lastName),
		ByteSliceToString(email),
		token[:]).Scan(&userId);
	err!=nil{
		log.Println(err)
		ctx.Error("Unhandled error",404)
	}else{
		successfulAuth(ctx,strconv.Itoa(userId))
		ctx.Redirect("/secretpage",200)
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

func successfulAuth(ctx *fasthttp.RequestCtx, userId string){
	var access_token string
	for {
		access_token = Hasher(128)
		if _,err:= rdb.Get(access_token).Result();err!=nil{
			break
		}
	}

	accessTokenCookie :=CreateCookie("accessToken",access_token,36000000000)
	idCookie :=CreateCookie("userId",userId,36000000000)
	rdb.Set(userId,access_token,360000000000)
	ctx.Response.Header.SetCookie(accessTokenCookie)
	ctx.Response.Header.SetCookie(idCookie)
}

func ByteSliceToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func StringToByteSlice(str string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}
