package main

import (
	//"net/http"
	"snapsta/controllers"

	"github.com/gin-gonic/gin"
	//"gorm.io/driver/sqlite"
	//"gorm.io/gorm"
	//"snapsta/models"
)

func main() {
	router := gin.Default()
	router.Static("/media", "./media")
	router.LoadHTMLGlob("templates/*.html")
	api := router.Group("/api/")
	api.GET("/post/:id", controllers.PostDetails)
	api.GET("/follows", controllers.GetallFollows)
	api.GET("/getcomments", controllers.GetCommentsofapost)
	api.GET("/all", controllers.AllUser)
	api.GET("/userpost", controllers.Postfromuser)
	api.GET("/getloggeduser", controllers.Getloggeduserinfo)

	router.GET("/", controllers.Home)
	router.GET("/auth/google/login", controllers.GetGoogleLoginURL)
	router.GET("/auth/google/callback", controllers.CallBackReceiver)
	router.GET("/signup", controllers.SignupForm)
	router.POST("/signup/submit", controllers.SignUpFormSubmit)
	router.POST("/auth/login", controllers.LoginSubmit)
	router.GET("/logout", controllers.Logout)
	router.GET("/createpost", controllers.CreatePost)
	router.POST("/createpost", controllers.CreatePost)
	router.GET("/allposts", controllers.ViewPosts)
	router.GET("/post/editpost", controllers.EditPost)
	router.GET("/post/", controllers.SearchPosts)
	router.POST("/post/editpost", controllers.EditPost)
	router.GET("/post/deletepost", controllers.DeletePost)
	router.GET("/post/comment/:id", controllers.CreateComment)
	router.POST("/post/comment/:id", controllers.CreateComment)
	router.GET("/post/like/:id", controllers.LikeUnlike)
	router.GET("/profile", controllers.ProfileView)
	router.POST("/profile", controllers.ProfileView)
	router.GET("/deleteaccount", controllers.DeleteAccount)
	router.GET("/profile/follow", controllers.FollowUnfolllow)
	router.GET("/profile/followers", controllers.ViewFollowers)
	router.GET("/profile/:username", controllers.ViewOthersProfile)

	router.Run("127.0.0.1:8000")

}
