package main

import (
	"bytes"
	// "fmt"
	"path/filepath"
	// "regexp"
	"strings"

	"log"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Format struct {
	ID         string `json:"id"`
	Extension  string `json:"ext"`
	Resolution string `json:"resolution"`
	FPS        string `json:"fps"`
}

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
		log.Printf("Running command: %s\n", cmd.String())

		output, err := cmd.Output()
		if err != nil {
			log.Printf("Command error: %s\n", output)
			panic(err)
		}

		formats := parseYtDlpOutput(string(output))
		// jsonOutput, _ := json.MarshalIndent(formats, "", "  ")

		// fmt.Println(string(jsonOutput))
		// var out bytes.Buffer
		// var stderr bytes.Buffer
		// cmd.Stdout = &out
		// cmd.Stderr = &stderr

		// err = cmd.Run()
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 		"error":   "Failed to execute yt-dlp",
		// 		"details": stderr.String(),
		// 	})
		// }

		// log.Printf("Command output: %s\n", out.String())

		// Return the command output
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			// "formats": out.String(),
			"formats": formats,
		})
	})

	app.Get("/youtube/download", func(c *fiber.Ctx) error {

		URL := c.Query("url")
		startTime := c.Query("start")
		endTime := c.Query("end")
        // video := c.Query("video")
        // audio := c.Query("audio")

		// if URL == "" || startTime == "" || endTime == "" || video == "" || audio == "" {
        if URL == "" || startTime == "" || endTime == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bad Request, Missing query parameters ( url, start, end, video, audio )",
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

		c.Set("Content-Type", "video/mp4")
		c.Set("Content-Disposition", "attachment; filename=\"ClipIt.mp4\"")
		return c.SendFile(absPath, true)

		// return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// 	"message": "video downloaded successfully",
		// })
	})

	log.Fatal(app.Listen(":8000"))
}

func parseYtDlpOutput(output string) []Format {
	lines := strings.Split(output, "\n")
	var formats []Format
	startParsing := false

	for _, line := range lines {
		// Start parsing only after the table separator
		if strings.HasPrefix(line, "ID") {
			startParsing = true
			continue
		}
		if !startParsing || strings.TrimSpace(line) == "" || strings.HasPrefix(line, "--") {
			continue
		}

		// Split line into fields
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		format := Format{
			ID:         fields[0],
			Extension:  fields[1],
			Resolution: fields[2],
			FPS:        fields[3],
		}

		formats = append(formats, format)
	}

	return formats
}

// yt-dlp \\n  -f 231+234 \\n  --merge-output-format mp4 \\n  --download-sections "*20:30-23:40" \\n  --force-keyframes-at-cuts \\n  --no-playlist \\n  -o "clip2.%(ext)s" \\n  "https://www.youtube.com/watch?v=3-ELBiUkUWc"\n
//  700  yt-dlp \\n  -f 231+234 \\n  --merge-output-format mp4 \\n  --download-sections "*45:41-48:26" \\n  --force-keyframes-at-cuts \\n  --no-playlist \\n  -o "clip3.%(ext)s" \\n  "https://www.youtube.com/watch?v=3-ELBiUkUWc"\n
