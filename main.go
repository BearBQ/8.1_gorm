package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name  string
	Email string `gorm:"type:varchar(40);unique"`
	Posts []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint
	User    Users `gorm:"foreignKey:UserID"`
}

func main() {

	//содержимое .env:
	//DATABASE_URL="host=postgres.orbstack-pg.orb.local user=admin password=di1mon11421 dbname=gorm_learning port=5432"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL environment variable is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to base")
	//ctx := context.Background()
	db.AutoMigrate(&Users{}, &Post{})

	//err = gorm.G[Users](db).Create(ctx, &Users{Name: "no admin", Email: "asd@asd.ru"})
	_ = err

	newPosts := []Post{
		{Title: "1",
			Content: "123123",
			UserID:  1,
		},
		{Title: "2",
			Content: "2222222",
			UserID:  1},
	}
	for _, post := range newPosts {
		err = AddPost(db, post)
		if err != nil {
			log.Println(err)
		}
	}

	resultUser, err := GetUser(db, 2) //Получаем запись по номеру ID
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resultUser)

	resultUser, err = GetUserWithEmail(db, "asd@asd.ru")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resultUser)

	resultUser, err = GetUsersWithPosts(db, 1)
	if err != nil {
		log.Println(err)
	}

	for _, post := range resultUser.Posts {
		fmt.Println(post.Title, post.Content)
	}

	err = DeleteUserAndPosts(db, 1)
	if err != nil {
		log.Println(err)
	}
}

func AddPost(db *gorm.DB, post Post) error {

	result := db.Create(&post)
	return result.Error
}

func GetUser(db *gorm.DB, id uint) (Users, error) {
	var user Users
	result := db.First(&user, id)
	return user, result.Error
}

func GetUserWithEmail(db *gorm.DB, email string) (Users, error) {
	var user Users
	result := db.Where("Email=?", email).First(&user)
	return user, result.Error
}

func GetUsersWithPosts(db *gorm.DB, id uint) (Users, error) { //выборка пользователя с постами
	var user Users
	result := db.Preload("Posts").First(&user, id)
	return user, result.Error
}

func GetAllUsersWithPosts(db *gorm.DB) (Users, error) { //выборка пользователя с постами
	var user Users
	result := db.Preload("Posts").Find(&user)
	return user, result.Error
}

func DeleteUserAndPosts(db *gorm.DB, id uint) error {
	return db.Select("Posts").Delete(&Users{}, id).Error
}
