package main

import (
	"bytes"
	"path/filepath"
	// "fmt"
	// "context"
	// "encoding/json"
	"log"
	// "net/url"
	// "os"
	"os/exec"
	// "time"

	// "github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// YOUTUBE_API_KEY := os.Getenv("YOUTUBE_API_KEY")

	app := fiber.New()
	// client := resty.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/youtube/formats", func(c *fiber.Ctx) error {
		// Get the 'url' query parameter

		URL := c.Query("url")
		if URL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing 'url' query parameter",
			})
		}

		cmd := exec.Command("yt-dlp", "--list-formats", URL)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		log.Printf("Running command: %s\n", cmd.String())
		log.Printf("Command error: %s\n", stderr.String())

		err = cmd.Run()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to execute yt-dlp",
				"details": stderr.String(),
			})
		}

		log.Printf("Command output: %s\n", out.String())

		// Return the command output
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": out.String(),
		})
	})

	app.Get("/youtube/download", func(c *fiber.Ctx) error {

		URL := c.Query("url")
		startTime := c.Query("start")
		endTime := c.Query("end")

		if URL == "" || startTime == "" || endTime == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bad Request, Missing query parameters ( url, start, end )",
			})
		}

		outputFile := "clip.mp4"

		cmd := exec.Command("yt-dlp", "-f", "231+234", "--merge-output-format", "mp4", "--download-sections", "*"+startTime+"-"+endTime+"", "--force-keyframes-at-cuts", "--no-playlist", "-o", outputFile, URL)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		log.Printf("Running command: %s\n", cmd.String())

		// Execute the command and wait for it to finish
		err = cmd.Run()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to execute yt-dlp to download",
				"details": stderr.String(),
			})
		}

		absPath, err := filepath.Abs(outputFile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to get absolute path",
				"details": err.Error(),
			})
		}

		defer func() {
			delCmd := exec.Command("rm", outputFile)
			err = delCmd.Run()

			if err != nil {
				log.Printf("Error deleting file: %s\n", err.Error())
			}
		}()

		return c.SendFile(absPath, true)

		// return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// 	"message": "video downloaded successfully",
		// })
	})

	log.Fatal(app.Listen(":8001"))
}

// yt-dlp \\n  -f 231+234 \\n  --merge-output-format mp4 \\n  --download-sections "*20:30-23:40" \\n  --force-keyframes-at-cuts \\n  --no-playlist \\n  -o "clip2.%(ext)s" \\n  "https://www.youtube.com/watch?v=3-ELBiUkUWc"\n
//  700  yt-dlp \\n  -f 231+234 \\n  --merge-output-format mp4 \\n  --download-sections "*45:41-48:26" \\n  --force-keyframes-at-cuts \\n  --no-playlist \\n  -o "clip3.%(ext)s" \\n  "https://www.youtube.com/watch?v=3-ELBiUkUWc"\n
