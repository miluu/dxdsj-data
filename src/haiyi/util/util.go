package util

import (
	"bytes"
	"errors"
	"golanger.com/net/http/client"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func SaveImg(url, name, path string) error {
	out, err := os.Create(path + name)
	if err != nil {
		return err
	}
	defer out.Close()
	if resp, err := client.NewClient().Get(url); err != nil {
		return err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return errors.New("status code: " + strconv.Itoa(resp.StatusCode))
		}
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return err
		} else {
			if _, err := io.Copy(out, bytes.NewReader(body)); err != nil {
				return err
			}
		}
	}
	return nil
}

func ClearHtmlTags(str string) string {
	str = strings.TrimSpace(str)
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllStringFunc(str, strings.ToLower)
	reg := regexp.MustCompile(`<!--[^>]+>|<iframe[\S\s]+?</iframe>|<a[^>]+>|</a>|<script[\S\s]+?</script>`)
	str = reg.ReplaceAllString(str, "")
	return str
}

func GetPreMonth(y, m int) (year, month int) {
	month = m - 1
	year = y
	if month < 1 {
		month = 12
		year--
	}
	return
}

func GetNextMonth(y, m int) (year, month int) {
	month = m + 1
	year = y
	if month > 12 {
		month = 1
		year++
	}
	return
}

func FormatDate(y, m, d int) string {
	str := strconv.Itoa(y) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(d)
	return str
}
