package biligo

import (
	"fmt"
	"testing"
)

const (
	SEARCH_KEYWORD = "岁己"
	SEARCH_TYPE    = SEARCH_TYPE_LIVE
)

func TestSearchAllWbi(t *testing.T) {
	req := Chain{
		Req: NewGet(URL_SEARCH_ALL_WBI).WbiSign().
			WithQuery("keyword", SEARCH_KEYWORD),
	}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestSearchAllOld(t *testing.T) {
	req := Chain{
		Req: NewGet(URL_SEARCH_ALL_OLD).
			WithQuery("keyword", SEARCH_KEYWORD),
	}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestSearchTypeWbi(t *testing.T) {
	req := Chain{
		Req: NewGet(URL_SEARCH_TYPE_WBI).WbiSign().
			WithQuerys("search_type", string(SEARCH_TYPE), "keyword", SEARCH_KEYWORD),
	}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestSearchTypeOld(t *testing.T) {
	req := Chain{
		Req: NewGet(URL_SEARCH_TYPE_OLD).
			WithQuerys("search_type", string(SEARCH_TYPE), "keyword", SEARCH_KEYWORD),
	}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}
