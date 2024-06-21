package webpanel

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/base64"
)

const (
	GetArtisans = `SELECT entry_id, artisan_id, initials, nickname, gender, country_id, wii_number, mii_id, likes, perm_likes, mii_data FROM miis ORDER BY entry_id`
)

type Artisans struct {
	EntryId int
	ArtisanId int
	Initials string
	Nickname string
	Gender int
	CountryId int
	WiiNumber int
	MiiId []byte
	Likes int
	PermLikes int
	MiiData []byte
	MiiDataEncoded string
}

func (w *WebPanel) ViewArtisans(c *gin.Context) {
	rows, err := w.Pool.Query(w.Ctx, GetPlaza)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var plaza []Plaza
	for rows.Next() {
		plazadata := Plaza{}
		err = rows.Scan(&plazadata.EntryId, &plazadata.ArtisanId, &plazadata.Initials, &plazadata.Nickname, &plazadata.Gender, &plazadata.CountryId, &plazadata.WiiNumber, &plazadata.MiiId, &plazadata.Likes, &plazadata.PermLikes, &plazadata.MiiData)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		plazadata.MiiDataEncoded = base64.StdEncoding.EncodeToString(plazadata.MiiData)

		plaza = append(plaza, plazadata)
	}

	c.HTML(http.StatusOK, "view_artisans.html", gin.H{
		"numberOfMiis": len(plaza),
		"Plaza":         plaza,

	})
}
