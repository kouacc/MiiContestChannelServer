package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"net/http"
)

const (
	CheckMii = `SELECT entry_id FROM miis WHERE artisan_id = $1 AND mii_id = $2`
)

func check(c *gin.Context) {
	artisanId := c.PostForm("craftsno")
	strMiiId := c.PostForm("miiid")

	data, _ := hex.DecodeString("434800000000000000000000000000000000000000000000FFFFFFFFFFFFFFFF454E000C00000001")

	// Validity check for the Mii ID.
	miiId, err := hex.DecodeString(strMiiId)
	if err != nil {
		// TODO: Find invalid Mii ID error code
		c.Status(http.StatusBadRequest)
		writeResult(c, 400)
		return
	}

	// Check if the Mii is in the database.
	var entryId uint32
	err = pool.QueryRow(ctx, CheckMii, artisanId, miiId).Scan(&entryId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Mii does not exist, tell the channel to update, not add.
			data = binary.BigEndian.AppendUint32(data, 0)
		} else {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}
	} else {
		// Mii exists, set the entry id.
		data = binary.BigEndian.AppendUint32(data, entryId)
	}

	c.Data(http.StatusOK, "application/octet-stream", data)
}
