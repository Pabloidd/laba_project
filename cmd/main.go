package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"mini_project/internal/config"
	"mini_project/internal/model"
	"mini_project/internal/repository"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := sql.Open("pgx", cfg.SQLConnectionString)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}
	fmt.Println("подключение к PostgreSQL успешно")

	userRepo := repository.NewUserRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	user := &model.User{Nickname: "pavel"}
	if err := userRepo.Create(ctx, user); err != nil {
		log.Fatalf("create user: %v", err)
	}
	fmt.Printf("создан пользователь: id=%d nickname=%s\n", user.Id, user.Nickname)

	review, err := model.NewReview(0, *user, "хороший товар", 5)
	if err != nil {
		log.Fatalf("new review: %v", err)
	}
	if err := reviewRepo.Create(ctx, review); err != nil {
		log.Fatalf("create review: %v", err)
	}
	fmt.Printf("создан отзыв: id=%d stars=%d text=%q\n", review.Id, review.StarCount, review.Text)

	got, err := reviewRepo.GetByID(ctx, review.Id)
	if err != nil {
		log.Fatalf("get review: %v", err)
	}
	fmt.Printf(
		"прочитан отзыв: id=%d user=%s(%d) stars=%d text=%q\n",
		got.Id, got.User.Nickname, got.User.Id, got.StarCount, got.Text,
	)

	reviews, err := reviewRepo.GetAll(ctx)
	if err != nil {
		log.Fatalf("get all reviews: %v", err)
	}
	fmt.Printf("всего отзывов в БД: %d\n", len(reviews))
}
