package webpanel

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	GetArtisans = `SELECT artisan_id, name, country_id, wii_number, mac_address, number_of_posts, total_likes, is_master, last_post, mii_data FROM artisans ORDER BY artisan_id`
)

type Artisans struct {
	ArtisanId int
	Name string
	CountryId int
	WiiNumber int
	MacAddress string
	NumberOfPosts int
	TotalLikes int
	IsMaster bool
	LastPost time.Time
	MiiData []byte
}

func (w *WebPanel) ViewArtisans(c *gin.Context) {
	rows, err := w.Pool.Query(w.Ctx, GetArtisans)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var artisans []Artisans
	for rows.Next() {
		artisan := Artisans{}
		err = rows.Scan(&artisan.ArtisanId, &artisan.Name, &artisan.CountryId, &artisan.WiiNumber, &artisan.MacAddress, &artisan.NumberOfPosts, &artisan.TotalLikes, &artisan.IsMaster, &artisan.LastPost, &artisan.MiiData)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		artisans = append(artisans, artisan)
	}

	c.HTML(http.StatusOK, "view_artisans.html", gin.H{
		"numberOfArtisans": len(artisans),
		"Artisans":         artisans,

	})
}
