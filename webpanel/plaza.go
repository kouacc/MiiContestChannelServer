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
	GetArtisanInfo = `SELECT name FROM artisans where artisan_id = $1`
	SearchMiis = `SELECT entry_id, artisan_id, initials, nickname, gender, country_id, wii_number, mii_id, likes, perm_likes, mii_data FROM miis WHERE nickname ILIKE '%' || $1 || '%' ORDER BY entry_id`
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
	ArtisanName string
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

		row := w.Pool.QueryRow(w.Ctx, GetArtisanInfo, plazadata.ArtisanId)
		var artisanName string
		err = row.Scan(&artisanName)
		if err != nil {
    	c.HTML(http.StatusInternalServerError, "error.html", gin.H{
    	    "Error": err.Error(),
    	})
    		return
		}

		plazadata.ArtisanName = artisanName

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
	
	namerow := w.Pool.QueryRow(w.Ctx, GetArtisanInfo, MiiDetails.ArtisanId)
		var artisanName string
		err = namerow.Scan(&artisanName)
		if err != nil {
    	c.HTML(http.StatusInternalServerError, "error.html", gin.H{
    	    "Error": err.Error(),
    	})
    		return
		}
		
		MiiDetails.ArtisanName = artisanName

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

	c.Redirect(http.StatusFound, "/panel/plaza#delete_success")
}

func (w *WebPanel) SearchPlaza(c *gin.Context) {
	search := c.PostForm("search")

	rows, err := w.Pool.Query(w.Ctx, SearchMiis, search)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var SearchResults []Plaza
	for rows.Next() {
		search := Plaza{}
		err = rows.Scan(&search.EntryId, &search.ArtisanId, &search.Initials, &search.Nickname, &search.Gender, &search.CountryId, &search.WiiNumber, &search.MiiId, &search.Likes, &search.PermLikes, &search.MiiData)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		search.MiiDataEncoded = base64.StdEncoding.EncodeToString(search.MiiData)

		SearchResults = append(SearchResults, search)
	}

	c.HTML(http.StatusOK, "search_results.html", gin.H{
		"SearchResults": SearchResults,
		"SearchTerm":	 search,
	})
}