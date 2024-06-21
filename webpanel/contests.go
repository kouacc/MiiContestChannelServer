package webpanel

import (
	"bytes"
	"fmt"
	"github.com/WiiLink24/MiiContestChannel/contest"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	GetContests = `SELECT contest_id, english_name, status, has_souvenir, has_thumbnail, has_special_award, open_time, close_time FROM contests ORDER BY contest_id`
	AddContest  = `INSERT INTO contests (has_special_award, has_thumbnail, has_souvenir, english_name, status, open_time, close_time) 
					VALUES ($1, $2, $3, $4, 'waiting', $5, $6) RETURNING contest_id`
	GetOneContest = `SELECT contest_id, english_name, status, has_souvenir, has_thumbnail, has_special_award, open_time, close_time FROM contests WHERE contest_id = $1`
)

type Contests struct {
	ContestId int
	Name      string
	Status    string
	HasSouvenir bool
	HasThumbnail bool
	HasSpecialAward bool
	OpenTime  time.Time
	CloseTime time.Time
}

func (w *WebPanel) ViewContests(c *gin.Context) {
	rows, err := w.Pool.Query(w.Ctx, GetContests)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var contests []Contests
	var activeContests []Contests
	var waitingContests []Contests
	for rows.Next() {
		contest := Contests{}
		err = rows.Scan(&contest.ContestId, &contest.Name, &contest.Status, &contest.HasSouvenir, &contest.HasThumbnail, &contest.HasSpecialAward, &contest.OpenTime, &contest.CloseTime)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		contests = append(contests, contest)
		if contest.Status == "waiting" {
			activeContests = append(activeContests, contest)
		} else if contest.Status == "results" {
			waitingContests = append(waitingContests, contest)
		}
	}

	c.HTML(http.StatusOK, "view_contests.html", gin.H{
		"numberOfContests": len(contests),
		"Contests":         contests,
		"numberOfActiveContests":  len(activeContests),
		"ActiveContests": 		activeContests,
		"numberOfWaitingContests": len(waitingContests),
		"WaitingContests": waitingContests,

	})
}

func (w *WebPanel) AddContest(c *gin.Context) {
	c.HTML(http.StatusOK, "add_contest.html", nil)
}

func (w *WebPanel) AddContestPOST(c *gin.Context) {
	name := c.PostForm("name")
	strSpecialAward := c.PostForm("special_award")
	strOpenTime := c.PostForm("open_time")
	thumbnail, _ := c.FormFile("thumbnail")
	souvenir, _ := c.FormFile("souvenir")

	specialAward := false
	hasThumbnail := false
	hasSouvenir := false
	if thumbnail != nil {
		hasThumbnail = true
	}

	if souvenir != nil {
		hasSouvenir = true
	}

	if strSpecialAward == "on" {
		specialAward = true
	}

	openTime, err := time.Parse("2006-01-02", strOpenTime)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var contestId uint32
	err = w.Pool.QueryRow(w.Ctx, AddContest, specialAward, hasThumbnail, hasSouvenir, name, openTime, openTime.AddDate(0, 0, 7)).Scan(&contestId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Write the images now that we have the contest ID.
	if hasThumbnail {
		f, err := thumbnail.Open()
		defer f.Close()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, f)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		resized, err := resize(buffer, 96, 96)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Create encrypted thumbnail
		err = contest.MakeThumbnail(resized, contestId)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Write unencrypted thumbnail
		err = os.WriteFile(fmt.Sprintf("%s/contest/%d/thumbnail.jpg", w.Config.AssetsPath, contestId), resized, 0666)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	c.Redirect(http.StatusPermanentRedirect, "/panel/contests")
}

func (w *WebPanel) EditContest(c *gin.Context) {
	contestId := c.Param("contest_id")
    //fetch the contest data
    row := w.Pool.QueryRow(w.Ctx, GetOneContest, contestId)

    var contestdata Contests
    err := row.Scan(&contestdata.ContestId, &contestdata.Name, &contestdata.Status, &contestdata.HasSouvenir, &contestdata.HasThumbnail, &contestdata.HasSpecialAward, &contestdata.OpenTime, &contestdata.CloseTime)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "Error": err.Error(),
        })
        return
    }


	c.HTML(http.StatusOK, "edit_contest.html", gin.H{
		"ContestInfo": contestdata,
	})
}

func (w *WebPanel) EditContestPOST(c *gin.Context) {
	c.HTML(http.StatusOK, "edit_contest.html", nil)
}
