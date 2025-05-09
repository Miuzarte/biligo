package biligo

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// WbiSign 为链接签名
func WbiSign(u *url.URL) error {
	return wbiKeys.Sign(u)
}

// WbiUpdate 无视过期时间更新
func WbiUpdate() error {
	return wbiKeys.Update()
}

func WbiGet() (wk WbiKeys, err error) {
	if err = wk.update(false); err != nil {
		return
	}
	return wbiKeys, nil
}

var wbiKeys WbiKeys

type WbiKeys struct {
	Img            string
	Sub            string
	Mixin          string
	lastUpdateTime time.Time
}

// Sign 为链接签名
func (wk *WbiKeys) Sign(u *url.URL) (err error) {
	if err = wk.update(false); err != nil {
		return err
	}
	values := u.Query()
	values = removeUnwantedChars(values, '!', '\'', '(', ')', '*') // 必要性存疑?
	values.Set("wts", strconv.FormatInt(time.Now().Unix(), 10))
	hash := md5.Sum([]byte(values.Encode() + wk.Mixin))
	values.Set("w_rid", hex.EncodeToString(hash[:]))
	u.RawQuery = values.Encode()
	return nil
}

// Update 无视过期时间更新
func (wk *WbiKeys) Update() (err error) {
	return wk.update(true)
}

// update 按需更新
func (wk *WbiKeys) update(purge bool) error {
	if !purge && time.Since(wk.lastUpdateTime) < time.Hour {
		return nil
	}

	nav, err := FetchNav()
	if err != nil {
		return err
	}

	img := nav.WbiImg.ImgUrl
	sub := nav.WbiImg.SubUrl
	if img == "" || sub == "" {
		return wrapErr(ErrWbiEmptyUrls, nav)
	}

	imgParts := strings.Split(img, "/")
	subParts := strings.Split(sub, "/")
	imgPng := imgParts[len(imgParts)-1]
	subPng := subParts[len(subParts)-1]
	wbiKeys.Img = strings.TrimSuffix(imgPng, ".png")
	wbiKeys.Sub = strings.TrimSuffix(subPng, ".png")
	wbiKeys.mixin()
	wbiKeys.lastUpdateTime = time.Now()
	return nil
}

func (wk *WbiKeys) mixin() {
	var mixin [32]byte
	wbi := wk.Img + wk.Sub
	for i := range mixin {
		mixin[i] = wbi[mixinKeyEncTab[i]]
	}
	wk.Mixin = string(mixin[:])
}

var mixinKeyEncTab = [...]int{
	46, 47, 18, 2, 53, 8, 23, 32,
	15, 50, 10, 31, 58, 3, 45, 35,
	27, 43, 5, 49, 33, 9, 42, 19,
	29, 28, 14, 39, 12, 38, 41, 13,
	37, 48, 7, 16, 24, 55, 40, 61,
	26, 17, 0, 1, 60, 51, 30, 4,
	22, 25, 54, 21, 56, 59, 6, 63,
	57, 62, 11, 36, 20, 34, 44, 52,
}

func removeUnwantedChars(v url.Values, chars ...byte) url.Values {
	b := []byte(v.Encode())
	for _, c := range chars {
		b = bytes.ReplaceAll(b, []byte{c}, nil)
	}
	s, err := url.ParseQuery(string(b))
	if err != nil {
		panic(err)
	}
	return s
}
