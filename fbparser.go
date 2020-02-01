package fbparser

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	likeReg      *regexp.Regexp
	shareReg     *regexp.Regexp
	commentReg   *regexp.Regexp
	timestampReg *regexp.Regexp
)

func init() {
	likeReg = regexp.MustCompile("reaction_count:{count:(\\d+)")
	shareReg = regexp.MustCompile("share_count:{count:(\\d+)")
	commentReg = regexp.MustCompile("i18n_comment_count:\"(\\d+)")
	timestampReg = regexp.MustCompile("data-utime=\"(\\d+)\"")
}

func GetLikes(lnk string) (pi PostInfo, err error) {
	// Читаем страницу
	resp, err := http.Get(lnk)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Что-то пошло не так
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		log.Println("[error]", resp.Status, resp.StatusCode)
		return
	}

	// Читаем ответ
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	likes := likeReg.FindStringSubmatch(string(content))
	if len(likes) == 0 {
		err = errors.New("can't find like data")
		return
	}

	comments := commentReg.FindStringSubmatch(string(content))
	if len(comments) == 0 {
		err = errors.New("can't find comment data")
		return
	}

	reposts := shareReg.FindStringSubmatch(string(content))
	if len(reposts) == 0 {
		err = errors.New("can't find repost data")
		return
	}

	timestamp := timestampReg.FindStringSubmatch(string(content))
	if len(timestamp) == 0 {
		err = errors.New("can't find repost data")
		return
	}

	l, _ := strconv.ParseInt(likes[1], 10, 32)
	r, _ := strconv.ParseInt(reposts[1], 10, 32)
	c, _ := strconv.ParseInt(comments[1], 10, 32)
	t, _ := strconv.ParseInt(timestamp[1], 10, 32)

	pi = PostInfo{
		Likes:     int(l),
		Reposts:   int(r),
		Comments:  int(c),
		Published: int(t),
		Text:      string(content),
	}

	return
}
