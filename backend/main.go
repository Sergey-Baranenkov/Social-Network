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

	//Router.GET("/auth", CORSHandler(authPageHandler))
	Router.POST("/registration", CORSHandler(RegistrationHandler))
	Router.POST("/login", CORSHandler(loginHandler))

	Router.GET("/static/*filepath", fasthttp.FSHandler("../frontend", 1))
	Router.GET("/frontend/*filepath", fasthttp.FSHandler("../frontend", 1))

	Router.GET("/profile/posts", CORSHandler(GetPostsHandler))
	Router.GET("/profile/get_comments", CORSHandler(GetCommentsHandler))
	Router.POST("/profile/add_new_object",CORSHandler(AddNewObjectHandler))
	Router.POST("/profile/update_like", CORSHandler(UpdateLikeHandler))
	Router.GET("/profile/page_info", CORSHandler(GetProfilePageInfo))

	Router.GET("/music/get_user_music", CORSHandler(GetUserMusicHandler))
	Router.GET("/music/get_combined_music", CORSHandler(GetCombinedMusicHandler))
	Router.POST("/music/post_music", CORSHandler(PostMusicHandler))
	Router.GET("/music/remove_music", CORSHandler(DeleteMusicHandler))
	Router.GET("/music_storage/*filepath", CORSHandler(fasthttp.FSHandler("../music_storage", 1)))

	Router.GET("/settings/hobbies", CORSHandler(HobbiesHandler))
	Router.GET("/settings/privacy", CORSHandler(PrivacyHandler))
	Router.GET("/settings/edu_emp", CORSHandler(EduEmpHandler))
	Router.POST("/settings/update_basic_info/text_data", CORSHandler(UpdateBasicInfoTextHandler))
	Router.POST("/settings/update_basic_info/profile_avatar", CORSHandler(AuthMiddleware(UpdateProfileAvatar)))
	Router.POST("/settings/update_basic_info/profile_bg", CORSHandler(AuthMiddleware(UpdateProfileBg)))

	Router.POST("/relations/subscribe", CORSHandler(SubscribeHandler))
	Router.POST("/relations/unsubscribe", CORSHandler(UnsubscribeHandler))
	Router.POST("/relations/add_subscriber_to_friend", CORSHandler(AddSubscriberToFriendHandler))
	Router.POST("/relations/add_friend_to_subscriber", CORSHandler(AddFriendToSubscriberHandler))
	Router.GET("/relations/get_relations",CORSHandler(GetRelationshipsHandler))

	Router.GET("/search_people",CORSHandler(GetSearchedPeople))


	Router.GET("/video/get_user_video", CORSHandler(GetUserVideoHandler))
	Router.GET("/video/get_combined_video", CORSHandler(GetCombinedVideoHandler))
	Router.POST("/video/post_video", CORSHandler(PostVideoHandler))
	Router.GET("/video_storage/*filepath", CORSHandler(fasthttp.FSHandler("../video_storage", 1)))

	Router.GET ("/gallery/get_images", CORSHandler(GetUserImages))
	Router.POST("/gallery/post_image", CORSHandler(PostImageHandler))
	Router.GET("/gallery/delete_image", CORSHandler(DeleteImageHandler))
	Router.GET("/gallery_storage/*filepath", CORSHandler(fasthttp.FSHandler("../gallery_storage", 1)))

	Router.GET("/messenger/conversation_list", CORSHandler(SelectConversationsList))
	Router.GET("/messenger/conversation_messages", CORSHandler(SelectConversationMessages))
	Router.POST("/messenger/push_message", CORSHandler(PushMessage))
	Router.GET("/messenger/", CORSHandler(MessengerHandler))
	Router.GET("/messenger/get_short_profile_info", CORSHandler(MessengerGetShortProfileInfo))

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