package main

import (
	"encoding/hex"
	"github.com/WiiLink24/nwc24"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	IsContestOpen    = `SELECT EXISTS (SELECT 1 FROM contests WHERE contest_id = $1 AND status = 'open')`
	IsFirstPost      = `SELECT COUNT(*) = 0 FROM contest_miis WHERE mac_address = $1 AND artisan_id = $2 AND contest_id = $3`
	InsertContestMii = `INSERT INTO contest_miis (contest_id, artisan_id, country_id, wii_number, mac_address, mii_data) 
						VALUES ($1, $2, $3, $4, $5, $6)`
	UpdateContestMii = `UPDATE contest_miis SET mii_data = $1 WHERE artisan_id = $2 AND mac_address = $3 AND contest_id = $4`
)

func conPost(c *gin.Context) {
	contestId := c.PostForm("contestno")
	artisanId := c.PostForm("craftsno")
	countryId := c.PostForm("country")
	miiData := []byte(c.PostForm("miidata"))
	macAddress := c.PostForm("macadr")
	strwiiNumber := c.PostForm("wiino")

	// Validate Wii Number
	wiiNumber, err := strconv.ParseUint(strwiiNumber, 10, 64)
	if err != nil {
		writeResult(c, 400)
		return
	}

	number := nwc24.LoadWiiNumber(wiiNumber)
	if !number.CheckWiiNumber() {
		writeResult(c, 310)
		return
	}

	// Validate Mii data
	if len(miiData) != 76 {
		writeResult(c, 305)
		return
	}

	// Check if the contest is open.
	var isContestOpen bool
	err = pool.QueryRow(ctx, IsContestOpen, contestId).Scan(&isContestOpen)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	if !isContestOpen {
		data, _ := hex.DecodeString("435000000000000000000000000000000000000000000000FFFFFFFFFFFFFFFF")
		c.Data(http.StatusOK, "application/octet-stream", data)
		return
	}

	// Check if this is the artisans first time posting
	var isFirstPost bool
	err = pool.QueryRow(ctx, IsFirstPost, macAddress, artisanId, contestId).Scan(&isFirstPost)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	if isFirstPost {
		_, err = pool.Exec(ctx, InsertContestMii, contestId, artisanId, countryId, wiiNumber, macAddress, miiData)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}
	} else {
		_, err = pool.Exec(ctx, UpdateContestMii, miiData, artisanId, macAddress, contestId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}
	}

	data, _ := hex.DecodeString("435000000000000000000000000000000000000000000000FFFFFFFFFFFFFFFF4E4C001000000001FFFFFFFFFFFFFFFF")
	c.Data(http.StatusOK, "application/octet-stream", data)
}
