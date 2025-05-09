package biligo

import (
	"testing"
)

func TestLocation(t *testing.T) {
	const shortUrl = `https://b23.tv/7qK0pRN`
	resp, err := httpClient.Do(NewHead(shortUrl))
	if err != nil {
		return
	}
	loc, ok := resp.Header["Location"]
	for k, v := range resp.Header {
		t.Log(k, v)
	}
	t.Log(!ok)
	t.Log(len(loc))
	t.Log(loc[0])
}

func TestGuestCookie(t *testing.T) {
	resp, err := httpClient.Do(NewGet("www.bilibili.com"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	t.Log(resp.Header.Get("Set-Cookie"))
}
