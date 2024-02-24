package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"unique"`
	Name      string `gorm:"default:''"`
	Password  string
	No_of_Followers int64  `gorm:"default:0"`
	No_of_Following int64  `gorm:"default:0"`
	Bio       string `gorm:"default:''"`
	Posts     []Post `gorm:"foreignkey:Author"`
	EmailID	string        `gorm:"unique"`
	ProfileImage string   `gorm:"default:''"`
	Followers []Follow     `gorm:"foreignkey:Wed"`
	Following []Follow     `gorm:"foreignkey:Wer"`
	Comments []Comment     `gorm:"foreignkey:Author_ID"`
	Active   bool
	
}

type Post struct {
	gorm.Model
	Title       string   `gorm:"default:''"`
	Description string   `gorm:"default:''"`
	Images      []Image `gorm:"foreignkey:Under_post"`
	Author      uint8
	Comments    []Comment `gorm:"foreignkey:Under_post"`
	Likes       []Like     `gorm:"foreignkey:Liked_post"`
	No_of_Likes uint8    `gorm:"default:0"`
}


type Comment struct {
	Under_post uint8  `gorm:"constraint:OnDelete:CASCADE"`
	Content    string `gorm:"default''"`
	Author_ID    uint8`gorm:"constraint:OnDelete:CASCADE"`
}
type Like struct{
	gorm.Model
	Liker uint 
	Liked_post uint 
	
}
type Image struct{
	gorm.Model
 Url string `gorm:"default:''"`
 Under_post uint `gorm:"constraint:OnDelete:CASCADE"`
}


type Follow struct{
	gorm.Model
	Wer uint      `gorm:"constraint:OnDelete:CASCADE"`
	Wed uint      `gorm:"constraint:OnDelete:CASCADE"`
	
}
