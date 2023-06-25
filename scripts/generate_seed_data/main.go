package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ebookstore/internal/core/auth"
	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/internal/migrator"
	"github.com/ebookstore/internal/platform/config"
	"github.com/ebookstore/internal/platform/hash"
	"github.com/ebookstore/internal/platform/payment"
	"github.com/ebookstore/internal/platform/storage"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	viper.AutomaticEnv()

	databaseURI := config.NewMigrationDatabaseURI()
	source := config.NewMigrationSource()
	migratorConfig := migrator.Config{
		DatabaseURI: databaseURI,
		Source:      source,
	}

	migrator := migrator.New(migratorConfig)
	migrator.Sync()

	db = config.NewConnection()
}

func main() {
	cleanDatabase()
	users := createUsers()
	log.Printf("Created %d Users", len(users))

	books := createBooks()
	log.Printf("Created %d Books", len(books))

	orders := createOrders()
	log.Printf("Created %d Orders", len(orders))
}

func cleanDatabase() {
	log.Println("Cleaning database...")

	if err := db.Delete(&shop.Order{}, "1 = 1").Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Delete(&catalog.Book{}, "1 = 1").Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Delete(&auth.User{}, "1 = 1").Error; err != nil {
		log.Fatal(err)
	}
}

func createUsers() []auth.User {
	log.Println("Creating users...")

	bcryptWrapper := hash.NewBcryptWrapper()

	hashedPassword, err := bcryptWrapper.HashPassword("password")
	if err != nil {
		log.Fatal(err)
	}

	customer1 := auth.User{
		ID:        "customer-id1",
		FirstName: "Raphael",
		LastName:  "Collin",
		Email:     "raphael@test.com",
		Role:      auth.Customer,
		Password:  hashedPassword,
	}
	if err = db.Create(&customer1).Error; err != nil {
		log.Fatal(err)
	}

	customer2 := auth.User{
		ID:        "customer-id2",
		FirstName: "Joe",
		LastName:  "Trump",
		Email:     "joe@test.com",
		Role:      auth.Customer,
		Password:  hashedPassword,
	}
	if err = db.Create(&customer2).Error; err != nil {
		log.Fatal(err)
	}

	admin := auth.User{
		ID:        "admin-id",
		FirstName: "Joe",
		LastName:  "Pratt",
		Email:     "admin@test.com",
		Role:      auth.Admin,
		Password:  hashedPassword,
	}
	if err = db.Create(&admin).Error; err != nil {
		log.Fatal(err)
	}

	return []auth.User{customer1, customer2, admin}
}

func createBooks() []catalog.Book {
	log.Println("Creating books...")

	storageClient := storage.NewStorage(storage.Config{
		S3Client: config.NewS3Client(config.NewAWSConfig()),
		Bucket:   config.NewBucket(),
	})

	book1Poster, err := os.Open("./book1_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	if err = storageClient.SaveFile(context.TODO(), "book1-poster", "image/jpeg", book1Poster); err != nil {
		log.Fatal(err)
	}

	book1Content, err := os.Open("./book1_content.pdf")
	if err != nil {
		log.Fatal(err)
	}
	if err = storageClient.SaveFile(context.TODO(), "book1-content", "application/pdf", book1Content); err != nil {
		log.Fatal(err)
	}

	book1 := catalog.Book{
		ID:          "book-id1",
		Title:       "Clean Code",
		Description: "Craftsman Guide",
		AuthorName:  "Robert c. Martin",
		ContentID:   "book1-content",
		Images: []catalog.Image{
			{
				ID:          "book1-poster",
				Description: "poster",
			},
		},
		Price:       0,
		ReleaseDate: time.Time{},
	}
	if err = db.Create(&book1).Error; err != nil {
		log.Fatal(err)
	}

	book2Poster, err := os.Open("./book2_image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	if err = storageClient.SaveFile(context.TODO(), "book2-poster", "image/jpeg", book2Poster); err != nil {
		log.Fatal(err)
	}

	book2Content, err := os.Open("./book2_content.pdf")
	if err != nil {
		log.Fatal(err)
	}
	if err = storageClient.SaveFile(context.TODO(), "book2-content", "application/pdf", book2Content); err != nil {
		log.Fatal(err)
	}

	book2 := catalog.Book{
		ID:          "book-id2",
		Title:       "Domain Driver Design",
		Description: "Tackling Complexity",
		AuthorName:  "Eric Evans",
		ContentID:   "book2-content",
		Images: []catalog.Image{
			{
				ID:          "book2-poster",
				Description: "poster",
			},
		},
		Price:       0,
		ReleaseDate: time.Time{},
	}
	if err = db.Create(&book2).Error; err != nil {
		log.Fatal(err)
	}

	return []catalog.Book{book1, book2}
}

func createOrders() []shop.Order {
	log.Println("Creating orders...")

	stripeClient := payment.NewStripePaymentService()

	order1 := shop.Order{
		ID:     "order-id1",
		Status: shop.Pending,
		Total:  40000,
		BookID: "book-id1",
		UserID: "customer-id1",
	}
	if err := stripeClient.CreatePaymentIntentForOrder(context.TODO(), &order1); err != nil {
		log.Fatal(err)
	}
	if err := db.Create(&order1).Error; err != nil {
		log.Fatal(err)
	}

	order2 := shop.Order{
		ID:     "order-id2",
		Status: shop.Pending,
		Total:  40000,
		BookID: "book-id2",
		UserID: "customer-id1",
	}
	if err := stripeClient.CreatePaymentIntentForOrder(context.TODO(), &order2); err != nil {
		log.Fatal(err)
	}
	order2.Status = shop.Paid

	if err := db.Create(&order2).Error; err != nil {
		log.Fatal(err)
	}

	order3 := shop.Order{
		ID:     "order-id3",
		Status: shop.Paid,
		Total:  40000,
		BookID: "book-id2",
		UserID: "customer-id2",
	}
	if err := stripeClient.CreatePaymentIntentForOrder(context.TODO(), &order3); err != nil {
		log.Fatal(err)
	}
	if err := db.Create(&order3).Error; err != nil {
		log.Fatal(err)
	}

	return []shop.Order{order1, order2, order3}
}
