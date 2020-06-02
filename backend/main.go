package main

import (
	"fmt"
	"github.com/lab259/cors"
	"github.com/valyala/fasthttp"
	"log"
	"coursework/websocketConnectionsMap"
)

var MessengerWebsocketStruct = websocketConnectionsMap.CreateWebsocketConnections()
	
func main() {
	if err := Initializer(); err != nil {
		log.Fatal("Провалена инициализация: ", err)
		return
	}

	Router.GET("/", CORSHandler(AuthMiddleware(authPageHandler)))

	Router.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.Redirect("/", 200)
	}

	Router.POST("/server/registration", CORSHandler(RegistrationHandler))
	Router.POST("/server/login", CORSHandler(loginHandler))


	Router.GET("/static/*filepath", fasthttp.FSHandler("../frontend/static", 1))

	Router.GET("/server/profile/get_posts", CORSHandler(AuthMiddleware(GetPostsHandler)))
	Router.GET("/server/profile/get_comments", CORSHandler(AuthMiddleware(GetCommentsHandler)))
	Router.POST("/server/profile/add_new_object",CORSHandler(AuthMiddleware(AddNewObjectHandler)))
	Router.POST("/server/profile/update_like", CORSHandler(AuthMiddleware(UpdateLikeHandler)))
	Router.GET("/server/profile/page_info", CORSHandler(AuthMiddleware(GetProfilePageInfo)))
	Router.POST("/server/profile/update_object_text", CORSHandler(AuthMiddleware(UpdateObjectText)))
	Router.POST("/server/profile/delete_object", CORSHandler(AuthMiddleware(DeleteObject)))

	Router.GET("/server/music/get_user_music", CORSHandler(AuthMiddleware(GetUserMusicHandler)))
	Router.GET("/server/music/get_combined_music", CORSHandler(AuthMiddleware(GetCombinedMusicHandler)))
	Router.POST("/server/music/post_music", CORSHandler(AuthMiddleware(PostMusicHandler)))
	Router.POST("/server/music/add_to_playlist",CORSHandler(AuthMiddleware(AddMusicToPlayList)))
	Router.GET("/server/music/remove_music", CORSHandler(AuthMiddleware(DeleteMusicHandler)))
	Router.GET("/server/music_storage/*filepath", CORSHandler(fasthttp.FSHandler("../music_storage", 2)))


	Router.GET("/server/settings/get_basic_info", CORSHandler(AuthMiddleware(GetBasicInfoHandler)))
	Router.GET("/server/settings/hobbies", CORSHandler(AuthMiddleware(GetHobbiesHandler)))
	Router.POST("/server/settings/update_hobbies", CORSHandler(AuthMiddleware(UpdateHobbiesHandler)))

	Router.GET("/server/settings/get_edu_emp", CORSHandler(AuthMiddleware(GetEduEmpHandler)))
	Router.POST("/server/settings/post_edu_emp", CORSHandler(AuthMiddleware(UpdateEduEmpHandler)))
	Router.POST("/server/settings/update_password", CORSHandler(AuthMiddleware(UpdatePasswordHandler)))

	Router.POST("/server/settings/update_basic_info/text_data", CORSHandler(AuthMiddleware(UpdateBasicInfoTextHandler)))
	Router.POST("/server/settings/update_basic_info/profile_avatar", CORSHandler(AuthMiddleware(AuthMiddleware(UpdateProfileAvatar))))
	Router.POST("/server/settings/update_basic_info/profile_bg", CORSHandler(AuthMiddleware(AuthMiddleware(UpdateProfileBg))))


	Router.POST("/server/relations/update_relationship", CORSHandler(AuthMiddleware(UpdateRelationshipHandler)))
	Router.GET("/server/relations/get_relations",CORSHandler(AuthMiddleware(GetRelationshipsHandler)))

	Router.GET("/server/search_people",CORSHandler(AuthMiddleware(GetSearchedPeople)))

	Router.GET("/server/video/get_user_video", CORSHandler(AuthMiddleware(GetUserVideoHandler)))
	Router.GET("/server/video/get_combined_video", CORSHandler(AuthMiddleware(GetCombinedVideoHandler)))
	Router.POST("/server/video/post_video", CORSHandler(AuthMiddleware(PostVideoHandler)))
	Router.GET("/server/video_storage/*filepath", CORSHandler(fasthttp.FSHandler("../video_storage", 2)))
	Router.POST("/server/video/add_to_playlist",CORSHandler(AuthMiddleware(AddVideoToPlayList)))
	Router.GET("/server/video/remove_video", CORSHandler(AuthMiddleware(DeleteVideoHandler)))

	Router.GET ("/server/gallery/get_images", CORSHandler(AuthMiddleware(GetUserImages)))
	Router.POST("/server/gallery/post_image", CORSHandler(AuthMiddleware(PostImageHandler)))
	Router.GET("/server/gallery/delete_image", CORSHandler(AuthMiddleware(DeleteImageHandler)))
	Router.GET("/server/gallery_storage/*filepath", CORSHandler(fasthttp.FSHandler("../gallery_storage", 2)))

	Router.GET("/server/messenger/conversation_list", CORSHandler(AuthMiddleware(SelectConversationsList)))
	Router.GET("/server/messenger/conversation_messages", CORSHandler(AuthMiddleware(SelectConversationMessages)))
	Router.POST("/server/messenger/push_message", CORSHandler(AuthMiddleware(PushMessage)))
	Router.GET("/server/messenger/", CORSHandler(AuthMiddleware(MessengerHandler)))
	Router.GET("/server/messenger/get_short_profile_info", CORSHandler(AuthMiddleware(MessengerGetShortProfileInfo)))

	Router.GET("/server/about_me/select_extended_user_info",CORSHandler(AuthMiddleware(SelectExtendedUserInfo)))

	Router.GET("/server/profile_bgs/*filepath", CORSHandler(fasthttp.FSHandler("../profile_bgs", 2)))

	fmt.Println("LISTENING ON PORT " + ServePort)
	server:=fasthttp.Server{MaxRequestBodySize: 1024*1024*1024, Handler: Router.Handler}

	if err := server.ListenAndServe(":"+ServePort); err != nil {
		log.Println("error when starting server: " + err.Error())
	}
	Postgres.Conn.Close()
	if err := Redis.Close(); err != nil {
		log.Println("error when closing Redis conn: " + err.Error())
	}
}

func CORSHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return cors.AllowAll().Handler(h)
}