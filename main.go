package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/abe27/api/crypto/configs"
	"github.com/abe27/api/crypto/routers"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	configs.Store, err = gorm.Open(sqlite.Open(os.Getenv("DBNAME")), &gorm.Config{
		DisableAutomaticPing:                     true,
		DisableForeignKeyConstraintWhenMigrating: false,
		SkipDefaultTransaction:                   true,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tbt_", // table name prefix, table for `User` would be `t_users`
			SingularTable: false,  // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   false,  // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("CID", "Cid"),
		},
	})

	if err != nil {
		panic("failed to connect database")
	}

	// Seed database
	configs.SeedDB()

	// Create Session Store
	configs.Session = session.New()
}

func main() {
	// Create config variable
	config := fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Web API Service", // add custom server header
		AppName:       "API Version 1.0",
		BodyLimit:     10 * 1024 * 1024, // this is the default limit of 10MB
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	}

	app := fiber.New(config)
	app.Use(cors.New())
	app.Use(requestid.New())

	// Custom File Writer
	dte := time.Now()
	file, err := os.OpenFile(fmt.Sprintf("./%s.log", dte.Format("20060102")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	app.Use(logger.New(logger.Config{
		Output: file,
	}))
	// app.Use(logger.New())
	app.Static("/", "./public")
	routers.Routers(app)
	app.Listen(":4000")
}
