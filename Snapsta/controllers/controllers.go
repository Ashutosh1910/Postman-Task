package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	//"strings"
	//"errors"
	"log"
	"net/http"
	"path/filepath"
	"snapsta/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var logged bool
var email string
var loggedUser models.User
var PFP_DIR = "/media/profile_pics"
var (
	gmailLogger = &oauth2.Config{
		RedirectURL:  "http://127.0.0.1:8000/auth/google/callback/",
		ClientID:     "250331160860-su1n7pgvjb8apdp9jt0iehbqjvlot08v.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "SLogger" // Replace with a random string or use a session-based state
)

type FullPost struct {
	Post         models.Post
	BY           models.User
	Islogged     bool
	Liked        bool
	Follows_user bool
}

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("snap.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
}
func Home(ctx *gin.Context) {
	if logged {
		ctx.Redirect(http.StatusSeeOther, "/allposts")
	} else {
		db.AutoMigrate(&models.Comment{}, &models.Post{}, &models.User{}, &models.Like{}, &models.Image{}, &models.Follow{})
		ctx.HTML(200, "base.html", gin.H{"title": "homepage",
			"logged": logged})
	}
}

func GetGoogleLoginURL(ctx *gin.Context) {
	url := gmailLogger.AuthCodeURL(oauthStateString)
	//fmt.Println(url)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func CallBackReceiver(ctx *gin.Context) {
	state := ctx.Query("state")
	if state != oauthStateString {
		log.Fatalln(state)
	}

	code := ctx.Query("code")
	token, err := gmailLogger.Exchange(ctx, code)
	if err != nil {
		log.Fatalln(err)
	}
	user, err := getUserInfo(token)
	if err != nil {
		log.Println(err)
	}
	//ctx.Set("loggedUser",user["email"])
	logged = true
	email = fmt.Sprintf("%v", user["email"])
	old_user := models.User{}
	exists := db.First(&old_user, "email_id=?", fmt.Sprintf("%v", user["email"]))
	if exists.Error != nil {
		ctx.Redirect(http.StatusSeeOther, "/signup")
	}
	loggedUser = old_user
	ctx.Redirect(http.StatusSeeOther, "/")

}

func getUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo/", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
func SignupForm(ctx *gin.Context) {

	ctx.HTML(200, "base.html", gin.H{"signup": true})
}
func SignUpFormSubmit(ctx *gin.Context) {
	var user models.User
	username := ctx.PostForm("username")
	var old models.User
	unique := db.Find(&old, "username=?", username)

	fmt.Println(unique.RowsAffected)
	if unique.RowsAffected == 0 { // to check unique username
		user.Username = username
		//email,logged:=ctx.Get("loggedUser")
		//fmt.Println(logged)
		if logged {
			if email != ctx.PostForm("email") {
				ctx.HTML(200, "base.html", gin.H{"signup": true, "notvalidemail": true, "email": email})
			} else {
				user.EmailID = ctx.PostForm("email")
			}
		}
		user.EmailID = ctx.PostForm("email")
		user.Password = ctx.PostForm("password")
		user.Name = ctx.PostForm("name")
		user.Bio = ctx.PostForm("bio")
		pic, err := ctx.FormFile("profile_pic")
		if err != nil {
			user.ProfileImage = ""
		} else {
			upload_to := "media/profile_pics/"
			picerr := ctx.SaveUploadedFile(pic, filepath.Join(upload_to, pic.Filename))
			if picerr != nil {
				panic(err)
			}
			user.ProfileImage = pic.Filename
		}
		logged = true
		email = user.EmailID
		db.Save(&user)
		loggedUser = user
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		ctx.HTML(200, "base.html", gin.H{"signup": true, "enteruniqueusername": true})
	}
}

func LoginSubmit(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	var user models.User
	exists := db.Find(&user, "username=?", username)
	//fmt.Println(username)
	//fmt.Println(exists.RowsAffected)
	if exists.RowsAffected == 0 || user.Password != password {

		logged = false
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		logged = true
		email = loggedUser.EmailID
		loggedUser = user

		ctx.Redirect(http.StatusSeeOther, "/")
	}
}

func AllUser(ctx *gin.Context) {
	var all []models.User
	db.Find(&all)
	ctx.JSON(200, gin.H{"all": all})
}
func Postfromuser(ctx *gin.Context) {

	var user models.User
	db.Preload("Posts").Find(&user, loggedUser.ID)
	ctx.JSON(200, gin.H{"userwithpost": user})

}
func Getloggeduserinfo(ctx *gin.Context) {
	if logged {

		db.Preload("Followers").Preload("Following").Find(&loggedUser, loggedUser.ID)
		ctx.JSON(200, gin.H{"loggedUser": loggedUser,
			"logged": logged})
	} else {
		ctx.JSON(200, gin.H{"loggedUser": nil})

	}
}
func Logout(ctx *gin.Context) {
	logged = false
	ctx.Redirect(http.StatusSeeOther, "/")
}
func CreatePost(ctx *gin.Context) {

	if logged {
		if ctx.Request.Method == "POST" {
			post_title := ctx.PostForm("title")
			fmt.Println(post_title)
			post_content := ctx.PostForm("content")

			imagesfiles, err := ctx.MultipartForm()
			if err != nil {
				panic(err)
			}
			images := imagesfiles.File["images"]
			fmt.Println(post_title)

			post := models.Post{
				Title:       post_title,
				Description: post_content,
				Author:      uint8(loggedUser.ID),
			}
			db.Create(&post)
			POST_DIR := "media/post_pics/"
			// var post_images []models.Image
			for _, image := range images {
				err := ctx.SaveUploadedFile(image, filepath.Join(POST_DIR, image.Filename))
				if err != nil {
					panic(err)
				}
				post_image := models.Image{Under_post: post.ID, Url: image.Filename}
				db.Create(&post_image)
			}
			ctx.Redirect(http.StatusSeeOther, "/")

		} else {
			ctx.HTML(200, "base.html", gin.H{"CreatePost": true, "logged": logged})
		}
	} else {
		ctx.Redirect(http.StatusSeeOther, "/")
	}
}

func Getallposts(ctx *gin.Context) {
	var userposts []models.Post
	db.Preload("Images").Find(&userposts, "author=?", loggedUser.ID)
	ctx.JSON(200, gin.H{"alluserposts": userposts})
}

func makefullposts(allposts []models.Post) []FullPost {
	var allfullpost []FullPost

	for _, post := range allposts {
		var like models.Like
		var postUser models.User
		db.Preload("Followers").Preload("Following").Find(&postUser, "ID=?", post.Author)
		fullpost := FullPost{
			Post: post, BY: postUser}
		if post.Author == uint8(loggedUser.ID) {
			fullpost.Islogged = true
		}
		liked := db.Where("liker=? AND liked_post=?", loggedUser.ID, post.ID).Find(&like)
		if liked.RowsAffected != 0 {
			fullpost.Liked = true
		}
		for _, follow := range postUser.Followers {
			if follow.Wer == loggedUser.ID {
				fullpost.Follows_user = true
				break
			}
		}

		allfullpost = append(allfullpost, fullpost)
	}
	return allfullpost
}
func ViewPosts(ctx *gin.Context) {
	if logged {

		var allposts []models.Post
		var allfullpost []FullPost
		db.Preload("Likes").Preload("Images").Find(&allposts)
		allfullpost = makefullposts(allposts)
		ctx.HTML(200, "base.html", gin.H{"logged": logged, "allposts": allfullpost})
	} else {
		ctx.Redirect(303, "/")
	}
}

func EditPost(ctx *gin.Context) {

	var post_to_edit models.Post
	id_string := ctx.Query("id")
	id, _ := strconv.Atoi(id_string)
	db.Preload("Images").Find(&post_to_edit, "id=?", (id))
	if logged && post_to_edit.Author == uint8(loggedUser.ID) {
		//fmt.Println(post_to_edit)

		if ctx.Request.Method == "POST" {

			fmt.Println(post_to_edit)
			new_title := ctx.PostForm("title")
			new_content := ctx.PostForm("content")
			imagesfiles, err := ctx.MultipartForm()
			if err != nil {
				panic(err)
			}
			images := imagesfiles.File["images"]
			post_to_edit.Title = new_title
			post_to_edit.Description = new_content
			if len(images) != 0 {
				POST_DIR := "media/post_pics/"
				for _, image := range images {
					var post_image models.Image
					err := ctx.SaveUploadedFile(image, filepath.Join(POST_DIR, image.Filename))
					if err != nil {
						panic(err)
					}
					post_image = models.Image{Under_post: post_to_edit.ID, Url: image.Filename}
					db.Create(&post_image)
				}
			}
			//fmt.Println(post_to_edit)

			db.Save(&post_to_edit)
			//var check models.Post
			//db.Preload("Images").Find(&check,"id=?",uint(id))
			//fmt.Println(check)

			ctx.Redirect(303, "/")
		} else if ctx.Request.Method == "GET" {
			ctx.HTML(200, "base.html", gin.H{"editpostform": true, "logged": logged, "post": post_to_edit})
		}

	} else {
		ctx.Redirect(http.StatusSeeOther, "/")
	}

}

func DeletePost(ctx *gin.Context) {

	var post_to_delete models.Post
	id, _ := strconv.Atoi(ctx.Query("id"))
	db.Find(&post_to_delete, id)
	if logged && post_to_delete.Author == uint8(loggedUser.ID) {
		db.Delete(&post_to_delete, id)
		ctx.Redirect(303, "/")
	} else {
		ctx.Redirect(303, "/")
	}

}
func CreateComment(ctx *gin.Context) {
	if logged {
		var new_comment models.Comment
		var post models.Post
		id, _ := strconv.Atoi(ctx.Param("id"))
		//fmt.Println(id)
		db.Preload("Comments").Find(&post, id)
		if ctx.Request.Method == "POST" {
			comment_content := ctx.PostForm("commenttext")
			new_comment.Content = comment_content
			new_comment.Author_ID = uint8(loggedUser.ID)
			new_comment.Under_post = uint8(post.ID)
			db.Create(&new_comment)
			ctx.Redirect(303, "/")
		} else {
			type fullcomment struct {
				Comment  models.Comment
				Username string
			}
			comments := post.Comments
			var fullcomments []fullcomment
			for _, comment := range comments {
				var user models.User
				db.Find(&user, comment.Author_ID)
				fullcomment := fullcomment{
					Comment:  comment,
					Username: user.Username,
				}
				fullcomments = append(fullcomments, fullcomment)
			}

			ctx.HTML(200, "base.html", gin.H{"Commentform": true, "logged": logged, "post": post, "Comments": fullcomments})
		}
	} else {
		ctx.Redirect(303, "/")
	}
}

func LikeUnlike(ctx *gin.Context) {
	if logged {
		id, _ := strconv.Atoi(ctx.Param("id"))
		var post models.Post
		db.Find(&post, id)
		var like_checker models.Like
		liked := db.Find(&like_checker, "liker=? AND liked_post=?", loggedUser.ID, id)
		if liked.RowsAffected == 0 {
			like := models.Like{
				Liker:      loggedUser.ID,
				Liked_post: uint(id),
			}
			post.No_of_Likes++
			db.Save(&like)
		} else {
			db.Delete(&like_checker, like_checker.ID)
			post.No_of_Likes--
		}
		db.Save(&post)

	}
	ctx.Redirect(303, "/")

}

func PostDetails(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var post models.Post
	db.Preload("Comments").Preload("Likes").Preload("Images").Find(&post, id)
	ctx.JSON(200, post)

}
func GetCommentsofapost(ctx *gin.Context) {

	var post models.Post
	db.Preload("Comments").Find(&post, 1)
	ctx.JSON(200, post.Comments)
}
func ProfileViewJson(ctx *gin.Context) {
	if logged {
		var fulluser models.User
		db.Preload("Posts").Preload("Followers").Preload("Following").Find(&fulluser, loggedUser.ID)
		ctx.JSON(200, gin.H{"user": fulluser})

	} else {
		ctx.Redirect(303, "/")
	}
}
func ProfileView(ctx *gin.Context) {
	if logged {
		var fulluser models.User
		db.Preload("Posts").Preload("Followers").Preload("Following").Find(&fulluser, loggedUser.ID)
		if ctx.Request.Method == "GET" {

			ctx.HTML(200, "base.html", gin.H{"user": fulluser, "logged": logged, "ownprofile": true})
		} else {
			username := ctx.PostForm("username")
			var old models.User
			unique := db.Find(&old, "username=?", username)
			fmt.Println(unique.RowsAffected)
			if unique.RowsAffected == 0 || loggedUser.Username == username { // to check unique username
				loggedUser.Username = username
				loggedUser.Name = ctx.PostForm("name")
				loggedUser.Bio = ctx.PostForm("bio")
				loggedUser.EmailID = ctx.PostForm("email")

				new_pfp, err := ctx.FormFile("profile-pic")
				if err == nil {
					err := ctx.SaveUploadedFile(new_pfp, "media/profile_pics/"+new_pfp.Filename)
					if err != nil {
						panic(err)
					}

					loggedUser.ProfileImage = new_pfp.Filename
				}
				db.Save(&loggedUser)
				ctx.Redirect(303, "/profile")

			} else {
				ctx.Redirect(303, "/profile")
			}
		}

	} else {
		ctx.Redirect(303, "/")
	}
}

func DeleteAccount(ctx *gin.Context) {

	if logged {
		db.Delete(&loggedUser, loggedUser.ID)
		ctx.Redirect(303, "/logout")
	} else {
		ctx.Redirect(303, "/")
	}
}
func FollowUnfolllow(ctx *gin.Context) {
	if logged {
		user_id, _ := strconv.Atoi(ctx.Query("id"))
		var user_to_follow models.User
		db.Find(&user_to_follow, user_id)

		var follow_check models.Follow
		follows := db.Find(&follow_check, "wer=? AND wed=?", loggedUser.ID, user_to_follow.ID)
		fmt.Println(follows.RowsAffected)
		if follows.RowsAffected == 0 {
			var follow models.Follow
			follow.Wer = loggedUser.ID
			follow.Wed = user_to_follow.ID
			db.Create(&follow)
			user_to_follow.No_of_Followers++
			loggedUser.No_of_Following++
		} else {
			db.Delete(&follow_check, follow_check.ID)
			user_to_follow.No_of_Followers--
			loggedUser.No_of_Following--
		}

		db.Save(&user_to_follow)
		db.Save(&loggedUser)
		ctx.Redirect(303, "/")

	} else {
		ctx.Redirect(303, "/")
	}
}
func ViewFollowers(ctx *gin.Context) {
	if logged {
		db.Preload("Followers").Preload("Following").Find(&loggedUser, loggedUser.ID)
		//TO VIEW FOLLOWERS
		var followerlist []models.User
		var followinglist []models.User
		for _, follower := range loggedUser.Followers {
			var user models.User
			db.Find(&user, follower.Wer)
			followerlist = append(followerlist, user)
		}
		for _, following := range loggedUser.Following {
			var user models.User
			db.Find(&user, following.Wed)
			followinglist = append(followinglist, user)
		}
		ctx.HTML(200, "base.html", gin.H{"followers": followerlist, "following": followinglist, "logged": logged, "followdata": true})

	} else {
		ctx.Redirect(303, "/")
	}
}

func GetallFollows(ctx *gin.Context) {

	var f []models.Follow
	db.Find(&f)
	ctx.JSON(200, f)
}

func SearchPosts(ctx *gin.Context) {
	filter := ctx.Query("search")
	if logged && filter != "" {
		var FilteredPosts []models.Post

		filter_tag := "%" + "#" + filter + "%"
		filter = "%" + filter + "%"
		db.Preload("Images").Find(&FilteredPosts, "title LIKE ? OR description LIKE ?", filter, filter_tag)
		fullposts := makefullposts(FilteredPosts)
		// fmt.Println(fullposts)

		ctx.HTML(200, "base.html", gin.H{"logged": logged, "allposts": fullposts})

	} else {
		ctx.Redirect(303, "/")
	}
}

func ViewOthersProfile(ctx *gin.Context) {

	username := ctx.Param("username")
	var user models.User
	var Islogged bool
	var followsuser bool
	exists := db.Preload("Posts").Preload("Followers").Preload("Following").Find(&user, "username=?", username)
	if exists.RowsAffected != 0 && logged {
		if user.Username == loggedUser.Username {
			Islogged = true
		}
		if follows := db.Find(&models.Follow{}, "wer=? AND wed=?", loggedUser.ID, user.ID).RowsAffected; follows != 0 {
			followsuser = true
		}
		ctx.HTML(200, "base.html", gin.H{"user": user, "logged": logged, "otherprofile": true, "Islogged": Islogged, "FollowsUser": followsuser})
	} else {
		ctx.Redirect(303, "/")
	}
}
