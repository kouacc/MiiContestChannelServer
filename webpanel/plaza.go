package webpanel

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/base64"
)

const (
	GetPlaza = `SELECT entry_id, artisan_id, initials, nickname, gender, country_id, wii_number, mii_id, likes, perm_likes, mii_data FROM miis ORDER BY entry_id`
	DeleteMii = `DELETE FROM miis WHERE entry_id = $1`
	GetMiiDetails = `SELECT entry_id, artisan_id, initials, nickname, gender, country_id, wii_number, mii_id, likes, perm_likes, mii_data FROM miis WHERE entry_id = $1`
)

type Plaza struct {
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

func (w *WebPanel) ViewPlaza(c *gin.Context) {
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

	c.HTML(http.StatusOK, "view_plaza.html", gin.H{
		"numberOfMiis": len(plaza),
		"Plaza":         plaza,

	})
}

func (w *WebPanel) ViewMiiDetails(c *gin.Context) {
	entryId := c.Param("entry_id")
	row := w.Pool.QueryRow(w.Ctx, GetMiiDetails, entryId)

	MiiDetails := Plaza{}
	err := row.Scan(&MiiDetails.EntryId, &MiiDetails.ArtisanId, &MiiDetails.Initials, &MiiDetails.Nickname, &MiiDetails.Gender, &MiiDetails.CountryId, &MiiDetails.WiiNumber, &MiiDetails.MiiId, &MiiDetails.Likes, &MiiDetails.PermLikes, &MiiDetails.MiiData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	MiiDetails.MiiDataEncoded = base64.StdEncoding.EncodeToString(MiiDetails.MiiData)

	c.HTML(http.StatusOK, "view_mii.html", gin.H{
		"MiiDetails": MiiDetails,
	})
}

func (w *WebPanel) DeleteMii(c *gin.Context) {
	entryId := c.Param("entry_id")
	_, err := w.Pool.Exec(w.Ctx, DeleteMii, entryId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/panel/plaza")
}