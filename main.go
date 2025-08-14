package main

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name  string
	Email string
}

func main() {
	dsn := "host=postgres.orbstack-pg.orb.local user=admin password=di1mon11421 dbname=gorm_learning port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Println("All fine", db, err)
	ctx := context.Background()
	db.AutoMigrate(&Users{})

	err = gorm.G[Users](db).Create(ctx, &Users{Name: "second", Email: "asd@asd.ru"})
	_ = err
}
