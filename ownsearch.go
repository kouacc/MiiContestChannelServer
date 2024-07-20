package main

import (
	"github.com/WiiLink24/MiiContestChannel/common"
	"github.com/WiiLink24/MiiContestChannel/plaza"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const GetMiisForArtisan = `SELECT entry_id, initials, perm_likes, skill, country_id, mii_data
				FROM miis WHERE artisan_id = $1
				ORDER BY likes LIMIT 50`

func ownSearch(c *gin.Context) {
	artisanId := c.Query("craftsno")

	var miis []common.MiiWithArtisan
	rows, err := pool.Query(ctx, GetMiisForArtisan, artisanId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	for rows.Next() {
		var likes int
		mii := common.MiiWithArtisan{}
		err = rows.Scan(&mii.EntryNumber, &mii.Initials, &likes, &mii.Skill, &mii.CountryCode, &mii.MiiData)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}

		// Downcast to u8 as database contains numbers larger.
		mii.Likes = uint8(likes)

		miis = append(miis, mii)
	}

	intArtisanId, err := strconv.Atoi(artisanId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		writeResult(c, 400)
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", plaza.MakeOwnSearch(miis, uint32(intArtisanId)))
}
