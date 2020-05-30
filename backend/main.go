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

	//Router.GET("/auth", CORSHandler(AuthMiddleware(authPageHandler))
	Router.POST("/registration", CORSHandler(RegistrationHandler))
	Router.POST("/login", CORSHandler(loginHandler))

	Router.GET("/static/*filepath", fasthttp.FSHandler("../frontend", 1))
	Router.GET("/frontend/*filepath", fasthttp.FSHandler("../frontend", 1))

	Router.GET("/profile/get_posts", CORSHandler(AuthMiddleware(GetPostsHandler)))
	Router.GET("/profile/get_comments", CORSHandler(AuthMiddleware(GetCommentsHandler)))
	Router.POST("/profile/add_new_object",CORSHandler(AuthMiddleware(AddNewObjectHandler)))
	Router.POST("/profile/update_like", CORSHandler(AuthMiddleware(UpdateLikeHandler)))
	Router.GET("/profile/page_info", CORSHandler(AuthMiddleware(GetProfilePageInfo)))
	Router.POST("/profile/update_object_text", CORSHandler(AuthMiddleware(UpdateObjectText)))
	Router.POST("/profile/delete_object", CORSHandler(AuthMiddleware(DeleteObject)))

	Router.GET("/music/get_user_music", CORSHandler(AuthMiddleware(GetUserMusicHandler)))
	Router.GET("/music/get_combined_music", CORSHandler(AuthMiddleware(GetCombinedMusicHandler)))
	Router.POST("/music/post_music", CORSHandler(AuthMiddleware(PostMusicHandler)))
	Router.POST("/music/add_to_playlist",CORSHandler(AuthMiddleware(AddMusicToPlayList)))
	Router.GET("/music/remove_music", CORSHandler(AuthMiddleware(DeleteMusicHandler)))
	Router.GET("/music_storage/*filepath", CORSHandler(fasthttp.FSHandler("../music_storage", 1)))


	Router.GET("/settings/get_basic_info", CORSHandler(AuthMiddleware(GetBasicInfoHandler)))
	Router.GET("/settings/hobbies", CORSHandler(AuthMiddleware(GetHobbiesHandler)))
	Router.POST("/settings/update_hobbies", CORSHandler(AuthMiddleware(UpdateHobbiesHandler)))

	Router.GET("/settings/get_edu_emp", CORSHandler(AuthMiddleware(GetEduEmpHandler)))
	Router.POST("/settings/post_edu_emp", CORSHandler(AuthMiddleware(UpdateEduEmpHandler)))
	Router.POST("/settings/update_password", CORSHandler(AuthMiddleware(UpdatePasswordHandler)))

	Router.POST("/settings/update_basic_info/text_data", CORSHandler(AuthMiddleware(UpdateBasicInfoTextHandler)))
	Router.POST("/settings/update_basic_info/profile_avatar", CORSHandler(AuthMiddleware(AuthMiddleware(UpdateProfileAvatar))))
	Router.POST("/settings/update_basic_info/profile_bg", CORSHandler(AuthMiddleware(AuthMiddleware(UpdateProfileBg))))


	Router.POST("/relations/update_relationship", CORSHandler(AuthMiddleware(UpdateRelationshipHandler)))
	Router.GET("/relations/get_relations",CORSHandler(AuthMiddleware(GetRelationshipsHandler)))

	Router.GET("/search_people",CORSHandler(AuthMiddleware(GetSearchedPeople)))

	Router.GET("/video/get_user_video", CORSHandler(AuthMiddleware(GetUserVideoHandler)))
	Router.GET("/video/get_combined_video", CORSHandler(AuthMiddleware(GetCombinedVideoHandler)))
	Router.POST("/video/post_video", CORSHandler(AuthMiddleware(PostVideoHandler)))
	Router.GET("/video_storage/*filepath", CORSHandler(fasthttp.FSHandler("../video_storage", 1)))
	Router.POST("/video/add_to_playlist",CORSHandler(AuthMiddleware(AddVideoToPlayList)))
	Router.GET("/video/remove_video", CORSHandler(AuthMiddleware(DeleteVideoHandler)))

	Router.GET ("/gallery/get_images", CORSHandler(AuthMiddleware(GetUserImages)))
	Router.POST("/gallery/post_image", CORSHandler(AuthMiddleware(PostImageHandler)))
	Router.GET("/gallery/delete_image", CORSHandler(AuthMiddleware(DeleteImageHandler)))
	Router.GET("/gallery_storage/*filepath", CORSHandler(fasthttp.FSHandler("../gallery_storage", 1)))

	Router.GET("/messenger/conversation_list", CORSHandler(AuthMiddleware(SelectConversationsList)))
	Router.GET("/messenger/conversation_messages", CORSHandler(AuthMiddleware(SelectConversationMessages)))
	Router.POST("/messenger/push_message", CORSHandler(AuthMiddleware(PushMessage)))
	Router.GET("/messenger/", CORSHandler(AuthMiddleware(MessengerHandler)))
	Router.GET("/messenger/get_short_profile_info", CORSHandler(AuthMiddleware(MessengerGetShortProfileInfo)))

	Router.GET("/about_me/select_extended_user_info",CORSHandler(AuthMiddleware(SelectExtendedUserInfo)))

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