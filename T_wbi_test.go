package biligo

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"testing"
)

func TestBytesReplaceAll(t *testing.T) {
	unwantedChars := []byte{'!', '\'', '(', ')', '*'}
	myBytes1 := []byte("1! 2' 3( 4) 5* 6")
	myBytes2 := []byte("!'()* 1 2 3 4 5 6")
	t.Log(string(bytes.ReplaceAll(myBytes1, unwantedChars, nil)))
	t.Log(string(bytes.ReplaceAll(myBytes2, unwantedChars, nil)))
}

func TestRegReplaceAll(t *testing.T) {
	r := regexp.MustCompile("[!()'*]")
	myBytes1 := []byte("1! 2' 3( 4) 5* 6")
	t.Log(string(r.ReplaceAll(myBytes1, nil)))
}

const TestWbiSignUrl = "https://api.bilibili.com/x/space/wbi/acc/info?mid=59442895"

func TestWbiSign(t *testing.T) {
	u, err := url.Parse(TestWbiSignUrl)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("before: %s", u)
	err = WbiSign(u)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signed: %s", u)
}

func TestWbiSignAndGet(t *testing.T) {
	testingLoadIdentity(t)

	t.Logf("before: %s", TestWbiSignUrl)
	u, err := url.Parse(TestWbiSignUrl)
	if err != nil {
		t.Fatal(err)
	}
	err = WbiSign(u)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signed: %s", u)

	req := Chain{Req: NewGet(TestWbiSignUrl).WbiSign()}
	err = req.Do()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(req.Body)
}
