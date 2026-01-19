package handler

import (
	"log"
	"net/http"
	"time"
	"url/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WORKER: remove stale data
// 2 am, get all the url info whose ttl > 10 years (time.now - ttl ~ >= 10 yrs)
// Then we'll lookup last clicked_at from user_clicks table, then we'll check if clicked_at's time falls around 2 months, if yes, then no need to remove it now, and we also need to update ttl for this short_id,
// New ttl :=  time.Now().Add(time.Hour * 24 *30 *12 * 10)
// We need to clear/delete that row/url

func CornJob(c *gin.Context) {
	// Calculate the date 10 years ago
	tenYearsAgo := time.Now().AddDate(-10, 0, 0)
	twoMonthsAgo := time.Now().AddDate(0, -2, 0)

	// Find all URLs with TTL older than 10 years
	var urls []model.URL
	result := dbCon.Where("ttl IS NOT NULL AND ttl <= ?", tenYearsAgo).Find(&urls)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("Error fetching stale URLs: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch stale URLs",
		})
		return
	}

	deletedCount := 0
	updatedCount := 0

	// Process each URL
	for _, url := range urls {
		// Get the last click for this URL
		var lastClick model.URLClicks
		clickResult := dbCon.Where("url_id = ?", url.ID).
			Order("clicked_at DESC").
			First(&lastClick)

		// If no clicks found, delete the URL
		if clickResult.Error == gorm.ErrRecordNotFound {
			if err := dbCon.Delete(&url).Error; err != nil {
				log.Printf("Error deleting URL with ID %d: %v", url.ID, err)
				continue
			}
			deletedCount++
			continue
		}

		// If there's an error other than record not found, skip this URL
		if clickResult.Error != nil {
			log.Printf("Error fetching clicks for URL ID %d: %v", url.ID, clickResult.Error)
			continue
		}

		// Check if last click is within the last 2 months
		if lastClick.ClickedAt.After(twoMonthsAgo) {
			// Last click is recent, update TTL to 10 years from now
			newTTL := time.Now().AddDate(10, 0, 0)
			if err := dbCon.Model(&url).Update("ttl", newTTL).Error; err != nil {
				log.Printf("Error updating TTL for URL ID %d: %v", url.ID, err)
				continue
			}
			updatedCount++
		} else {
			// Last click is older than 2 months, delete the URL
			if err := dbCon.Delete(&url).Error; err != nil {
				log.Printf("Error deleting URL with ID %d: %v", url.ID, err)
				continue
			}
			deletedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Cron job completed successfully",
		"deleted_count": deletedCount,
		"updated_count": updatedCount,
	})
}
