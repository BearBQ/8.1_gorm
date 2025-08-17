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

	if err = CreatePosts(db); err != nil {
		fmt.Println(err)
	}

	fmt.Println("работа продолжена")

}

func CreatePosts(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	user := Users{
		Name:  "Dima",
		Email: "asdd@as1dd",
		Posts: []Post{
			{Title: "пост номер 1",
				Content: "содержимое поста"},
			{Title: "пост номер 2",
				Content: "содержимое поста"},
		},
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err

	}
	log.Println("Транзакция успешна")
	return tx.Commit().Error
}
