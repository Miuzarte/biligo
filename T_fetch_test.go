package biligo

import (
	"fmt"
	"testing"
)

func TestGetUid(t *testing.T) {
	testingLoadIdentity(t)

	u, err := getSpaceUid()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}

func TestFetchShortLink(t *testing.T) {
	const u = `b23.tv/7qK0pRN`
	_, err := fetchLocation(u)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: complete result judgement

func TestFetchVideoInfo3(t *testing.T) {
	const id = `BV1UK42117py`
	_, _, _, err := FetchVideoInfo3(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchLiveStatus(t *testing.T) {
	const id = `1954091502`
	lsu, err := FetchLiveStatus(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lsu.Get(id).DoTemplate())
}

func TestFetchLiveRoomInfo(t *testing.T) {
	const id = `25788785`
	lri, err := FetchLiveRoomInfo(id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lri.DoTemplate())
}

func TestFetchArticle(t *testing.T) {
	const id = `cv40535917`
	_, err := FetchArticleInfo(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchSong(t *testing.T) {
	const id = `au3905874`
	_, err := FetchSong(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchSpaceCard(t *testing.T) {
	const id = `1954091502`
	_, err := FetchSpaceCard(id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFetchVideoPlayurlMp4(t *testing.T) {
	testingLoadIdentity(t)

	const bvid = "BV1UK42117py"
	cid, err := FetchVideoCidTryCache(bvid)
	if err != nil {
		t.Fatal(err)
	}

	req := Chain{Req: ReqVideoPlayurl().WithQuerys(
		"bvid", bvid, "cid", cid, "fnval", itoa(VIDEO_FNVAL_MP4),
		"fourk", "1", // 请求 4K
		"try_look", "1", // 游客高清晰度
	)}

	err = req.Do()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(req.Body)
	vp, err := req.ToVideoPlayurl()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(vp)
}

func TestFetchVideoPlayurlDash(t *testing.T) {
	testingLoadIdentity(t)

	const bvid = "BV1UK42117py"
	cid, err := FetchVideoCid(bvid)
	if err != nil {
		t.Fatal(err)
	}

	req := Chain{Req: ReqVideoPlayurl().WithQuerys(
		"bvid", bvid, "cid", cid, "fnval", itoa(VIDED_FNVAL_DASHALL),
		"fourk", "1", // 请求 4K
		"try_look", "1", // 游客高清晰度
	)}

	err = req.Do()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(req.Body)
	vp, err := req.ToVideoPlayurl()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(vp)
}
