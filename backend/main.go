package main

import (
	"coursework/websocketConnectionsMap"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

var MessengerWebsocketStruct = websocketConnectionsMap.CreateWebsocketConnections()

func main() {
	if err := Initializer(); err != nil {
		log.Fatal("Провалена инициализация: ", err)
		return
	}

	Router.NotFound = staticPageHandler

	Router.POST("/server/registration", RegistrationHandler)
	Router.POST("/server/login", loginHandler)

	Router.GET("/static/*filepath", fasthttp.FSHandler("../frontend/static", 1))

	Router.GET("/server/profile/get_posts", AuthMiddleware(GetPostsHandler))
	Router.GET("/server/profile/get_comments", AuthMiddleware(GetCommentsHandler))
	Router.POST("/server/profile/add_new_object", AuthMiddleware(AddNewObjectHandler))
	Router.POST("/server/profile/update_like", AuthMiddleware(UpdateLikeHandler))
	Router.GET("/server/profile/page_info", AuthMiddleware(GetProfilePageInfo))
	Router.POST("/server/profile/update_object_text", AuthMiddleware(UpdateObjectText))
	Router.POST("/server/profile/delete_object", AuthMiddleware(DeleteObject))

	Router.GET("/server/music/get_user_music", AuthMiddleware(GetUserMusicHandler))
	Router.GET("/server/music/get_combined_music", AuthMiddleware(GetCombinedMusicHandler))
	Router.POST("/server/music/post_music", AuthMiddleware(PostMusicHandler))
	Router.POST("/server/music/add_to_playlist", AuthMiddleware(AddMusicToPlayList))
	Router.GET("/server/music/remove_music", AuthMiddleware(DeleteMusicHandler))
	Router.GET("/server/music_storage/*filepath", fasthttp.FSHandler("../music_storage", 2))

	Router.GET("/server/settings/get_basic_info", AuthMiddleware(GetBasicInfoHandler))
	Router.GET("/server/settings/hobbies", AuthMiddleware(GetHobbiesHandler))
	Router.POST("/server/settings/update_hobbies", AuthMiddleware(UpdateHobbiesHandler))

	Router.GET("/server/settings/get_edu_emp", AuthMiddleware(GetEduEmpHandler))
	Router.POST("/server/settings/post_edu_emp", AuthMiddleware(UpdateEduEmpHandler))
	Router.POST("/server/settings/update_password", AuthMiddleware(UpdatePasswordHandler))

	Router.POST("/server/settings/update_basic_info/text_data", AuthMiddleware(UpdateBasicInfoTextHandler))
	Router.POST("/server/settings/update_basic_info/profile_avatar", AuthMiddleware(AuthMiddleware(UpdateProfileAvatar)))
	Router.POST("/server/settings/update_basic_info/profile_bg", AuthMiddleware(AuthMiddleware(UpdateProfileBg)))

	Router.POST("/server/relations/update_relationship", AuthMiddleware(UpdateRelationshipHandler))
	Router.GET("/server/relations/get_relations", AuthMiddleware(GetRelationshipsHandler))

	Router.GET("/server/search_people", AuthMiddleware(GetSearchedPeople))

	Router.GET("/server/video/get_user_video", AuthMiddleware(GetUserVideoHandler))
	Router.GET("/server/video/get_combined_video", AuthMiddleware(GetCombinedVideoHandler))
	Router.POST("/server/video/post_video", AuthMiddleware(PostVideoHandler))
	Router.GET("/server/video_storage/*filepath", fasthttp.FSHandler("../video_storage", 2))
	Router.POST("/server/video/add_to_playlist", AuthMiddleware(AddVideoToPlayList))
	Router.GET("/server/video/remove_video", AuthMiddleware(DeleteVideoHandler))

	Router.GET("/server/gallery/get_images", AuthMiddleware(GetUserImages))
	Router.POST("/server/gallery/post_image", AuthMiddleware(PostImageHandler))
	Router.GET("/server/gallery/delete_image", AuthMiddleware(DeleteImageHandler))
	Router.GET("/server/gallery_storage/*filepath", fasthttp.FSHandler("../gallery_storage", 2))

	Router.GET("/server/messenger/conversation_list", AuthMiddleware(SelectConversationsList))
	Router.GET("/server/messenger/conversation_messages", AuthMiddleware(SelectConversationMessages))
	Router.GET("/server/messenger/", AuthMiddleware(MessengerHandler))
	Router.GET("/server/messenger/get_short_profile_info", AuthMiddleware(MessengerGetShortProfileInfo))

	Router.GET("/server/about_me/select_extended_user_info", AuthMiddleware(SelectExtendedUserInfo))

	Router.GET("/server/profile_bgs/*filepath", fasthttp.FSHandler("../profile_bgs", 2))

	fmt.Println("LISTENING ON PORT " + ServePort)
	server := fasthttp.Server{MaxRequestBodySize: 1024 * 1024 * 1024, Handler: Router.Handler}

	if err := server.ListenAndServe(":" + ServePort); err != nil {
		log.Println("error when starting server: " + err.Error())
	}
	Postgres.Conn.Close()
	if err := Redis.Close(); err != nil {
		log.Println("error when closing Redis conn: " + err.Error())
	}
}
