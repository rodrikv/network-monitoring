package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"github.com/rodrikv/network-monitoring/internal/models"
	"github.com/rodrikv/network-monitoring/internal/ping"
)

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var statics embed.FS

// MonitoringResult represents the result of monitoring a source
type MonitoringResult struct {
	ID           int     `json:"id"`
	SeqID        int     `json:"seq_id"`
	Source       string  `json:"source"`
	Status       string  `json:"status"`
	ResponseTime float64 `json:"response_time"`
	TimeStamp    string  `json:"timestamp"`
}

// Database connection
func main() {
	// Initialize the database
	models.InitDB()

	args := os.Args

	m, err := ping.NewMonitor(
		args[1:],
		models.Db,
	)

	if err != nil {
		log.Fatal(err)
	}

	m.Start()

	// Start the monitoring job
	c := cron.New()
	c.AddFunc("@every 2d", dumpData)
	c.Start()

	// Create a new Gin router
	router := gin.Default()

	staticFS, err := fs.Sub(statics, "static")
	if err != nil {
		panic("Failed to create sub filesystem for static files: " + err.Error())
	}

	// Serve static files from the embedded filesystem
	router.StaticFS("/static", http.FS(staticFS))

	// Define routes
	temps, err := template.ParseFS(templates, "templates/*")

	if err != nil {
		log.Fatal(err)
	}
	router.SetHTMLTemplate(temps)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/chart", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chart.html", nil)
	})
	router.GET("/chart-data", chartData)
	router.GET("/historical-data", latestData)
	router.GET("/range-data", rangeData)

	// Run the application
	router.Run(":8080")
}

func isValidInput(input string) bool {
	// Allow only alphanumeric characters
	// You can customize this regular expression based on your specific requirements
	// For example, you might allow certain special characters if needed
	validInputRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return validInputRegex.MatchString(input)
}

func isValidTimestamp(input string) bool {
	_, err := strconv.ParseInt(input, 10, 64)
	return err == nil
}

func rangeData(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")

	// validate start and end for sql injection
	if !isValidInput(start) || !isValidInput(end) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !isValidTimestamp(start) || !isValidTimestamp(end) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start timestamp"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end timestamp"})
		return
	}

	rows, err := models.Db.Query("SELECT seq_id, source, status, response_time, timestamp FROM monitoring WHERE timestamp BETWEEN ? AND ? ORDER BY timestamp ASC", startTime, endTime)
	if err != nil {
		log.Println("Error fetching data:", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var results []MonitoringResult

	for rows.Next() {
		var result MonitoringResult
		err := rows.Scan(&result.SeqID, &result.Source, &result.Status, &result.ResponseTime, &result.TimeStamp)
		if err != nil {
			log.Println("Error scanning data:", err)
			return
		}
		results = append(results, result)
	}

	log.Print(results)

	c.JSON(http.StatusOK, results)
}

func fetchData() ([]MonitoringResult, error) {
	rows, err := models.Db.Query("SELECT * FROM monitoring")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MonitoringResult

	for rows.Next() {
		var result MonitoringResult
		err := rows.Scan(&result.ID, &result.Source, &result.Status, &result.ResponseTime, &result.TimeStamp)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func dumpData() {
	// Get data from the database
	results, err := fetchData()
	if err != nil {
		log.Println("Error fetching data:", err)
		return
	}

	// Marshal data to JSON
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Println("Error marshaling data to JSON:", err)
		return
	}

	// Create a new file for dumping data
	fileName := fmt.Sprintf("dump_%s.json", time.Now().Format("2006-01-02"))
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating dump file:", err)
		return
	}
	defer file.Close()

	// Write JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("Error writing data to dump file:", err)
		return
	}

	log.Printf("Data dumped to %s\n", fileName)
}

func chartData(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		results, err := fetchChartData()
		if err != nil {
			log.Println("Error fetching chart data:", err)
			continue
		}

		for _, result := range results {
			_, err = c.Writer.Write([]byte(fmt.Sprintf("data:%s\n\n", result)))

			if err != nil {
				log.Println("Error writing to stream:", err)
				return
			}
		}

		c.Writer.Flush()
	}
}

func fetchChartData() ([]string, error) {
	rows, err := models.Db.Query("SELECT seq_id,  source, status, response_time, timestamp FROM monitoring WHERE timestamp >= datetime('now', '-5 seconds') ORDER BY timestamp ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string

	for rows.Next() {
		var result MonitoringResult
		err := rows.Scan(&result.SeqID, &result.Source, &result.Status, &result.ResponseTime, &result.TimeStamp)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		results = append(results, string(data))
	}

	return results, nil
}

func latestData(c *gin.Context) {
	rows, err := models.Db.Query("SELECT seq_id, source, status, response_time, timestamp FROM monitoring ORDER BY id DESC LIMIT 1000")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	groupedResults := make(map[string][]MonitoringResult)

	for rows.Next() {
		var result MonitoringResult
		err := rows.Scan(&result.SeqID, &result.Source, &result.Status, &result.ResponseTime, &result.TimeStamp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Check if the source is already in the map
		if _, ok := groupedResults[result.Source]; !ok {
			groupedResults[result.Source] = make([]MonitoringResult, 0)
		}

		// Append the result to the corresponding source in the map
		groupedResults[result.Source] = append(groupedResults[result.Source], result)
	}

	// Create a list containing the length and records for each source
	var finalResults []map[string]interface{}
	for source, records := range groupedResults {
		result := make(map[string]interface{})
		result["source"] = source
		result["count"] = len(records)
		result["data"] = records
		finalResults = append(finalResults, result)
	}

	c.JSON(http.StatusOK, finalResults)
}
