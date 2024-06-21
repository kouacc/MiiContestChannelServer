package webpanel

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"encoding/base64"
)

const (
	GetArtisans = `SELECT artisan_id, name, country_id, wii_number, mac_address, number_of_posts, total_likes, is_master, last_post, mii_data FROM artisans ORDER BY artisan_id`
	GetArtisanDetails = `SELECT artisan_id, name, country_id, wii_number, mac_address, number_of_posts, total_likes, is_master, last_post, mii_data FROM artisans WHERE artisan_id = $1`
	GetArtisansMiis = `SELECT entry_id, artisan_id, initials, nickname, gender, country_id, wii_number, mii_id, likes, perm_likes, mii_data FROM miis WHERE artisan_id = $1` 
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
	MiiDataEncoded string
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

		artisan.MiiDataEncoded = base64.StdEncoding.EncodeToString(artisan.MiiData)

		artisans = append(artisans, artisan)
	}

	c.HTML(http.StatusOK, "view_artisans.html", gin.H{
		"numberOfArtisans": len(artisans),
		"Artisans":         artisans,

	})
}

func (w *WebPanel) ViewArtisanDetails(c *gin.Context) {
	artisanId := c.Param("artisan_id")
	row := w.Pool.QueryRow(w.Ctx, GetArtisanDetails, artisanId)

	artisanDetails := Artisans{}
	//fetch artisan details
	err := row.Scan(&artisanDetails.ArtisanId, &artisanDetails.Name, &artisanDetails.CountryId, &artisanDetails.WiiNumber, &artisanDetails.MacAddress, &artisanDetails.NumberOfPosts, &artisanDetails.TotalLikes, &artisanDetails.IsMaster, &artisanDetails.LastPost, &artisanDetails.MiiData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	artisanDetails.MiiDataEncoded = base64.StdEncoding.EncodeToString(artisanDetails.MiiData)

	//fetch artisan miis
	rows, err := w.Pool.Query(w.Ctx, GetArtisansMiis, artisanId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var Miis []Plaza
	for rows.Next() {
		miidata := Plaza{}
		err = rows.Scan(&miidata.EntryId, &miidata.ArtisanId, &miidata.Initials, &miidata.Nickname, &miidata.Gender, &miidata.CountryId, &miidata.WiiNumber, &miidata.MiiId, &miidata.Likes, &miidata.PermLikes, &miidata.MiiData)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		miidata.MiiDataEncoded = base64.StdEncoding.EncodeToString(miidata.MiiData)

		Miis = append(Miis, miidata)
}

	c.HTML(http.StatusOK, "artisan_details.html", gin.H{
		"ArtisanDetails": artisanDetails,
		"numberOfMiis": len(Miis),
		"Miis":         Miis,
	})
}