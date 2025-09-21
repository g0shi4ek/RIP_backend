package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/database"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := database.NewPostgresClient()
	if err != nil {
		panic("failed to connect postgres")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&domain.ChargingTariff{},
		&domain.User{},
		&domain.ChargingApplication{},
		&domain.ChargingOrder{},
	)
	if err != nil {
		panic("cant migrate db")
	}

	if err = postgresDataMigrations(db); err != nil{
		panic("cant migrate postgres data")
	}
	

	mc, err := database.NewMinioClient()
	if err != nil{
		panic("failed to connect minio")
	}
	if err = minioDataMigrations(mc); err != nil{
		panic("cant migrate minio data")
	}
}

func postgresDataMigrations(db *gorm.DB) error{
	//migrations/init_000.sql
	var count int64
	db.Model(&domain.ChargingTariff{}).Count(&count)
	if count > 0 {
		log.Println("postgres already has data")
		return nil
	}

	sqlFile := "migrations/init_000.sql"
	if _, err := os.Stat(sqlFile); err == nil {
		sqlBytes, err := os.ReadFile(sqlFile)
		if err != nil {
			return fmt.Errorf("failed to read sql file: %v", err)
		}

		err = db.Exec(string(sqlBytes)).Error
		if err != nil {
			return fmt.Errorf("failed to execute sql migration: %v", err)
		}
		log.Printf("Executed sql migration from %s", sqlFile)
	}

	return nil

}

func minioDataMigrations(mc *database.MinioClient) error{
	ctx := context.Background()
	bucketName := mc.BucketName

	exists, err := mc.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = mc.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		log.Printf("created minio bucket: %s", bucketName)

	} else {
		log.Printf("minio bucket %s already exists", bucketName)
	}

	return nil
}
