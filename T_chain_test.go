package biligo

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// TODO: use env
const testingIdentityPath = "./bilibili_identity.json"

func testingLoadIdentity(t interface{ Error(...any) }) {
	f, err := os.OpenFile(testingIdentityPath, os.O_RDONLY, 0o666)
	if err != nil {
		t.Error(err)
	}
	id := Identity{}
	err = json.NewDecoder(f).Decode(&id)
	if err != nil {
		t.Error(err)
	}
	err = StoreIdentity(id)
	if err != nil {
		t.Error(err)
	}
}

func TestChainVideoInfo(t *testing.T) {
	testingLoadIdentity(t)

	req := Chain{Req: ReqVideoInfo("1700234872")}
	err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	video, err := req.ToVideoInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(video.DoTemplate())
}

func TestChainVideoOnline(t *testing.T) {
	testingLoadIdentity(t)

	const aid = "1700234872"
	cid, err := FetchVideoCidTryCache(aid)
	if err != nil {
		t.Fatal(err)
	}
	req := Chain{Req: ReqVideoOnline(aid, cid)}
	err = req.Do()
	if err != nil {
		t.Fatal(err)
	}
	video, err := req.ToVideoOnline()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(video)
}

func TestChainVideoConclusion(t *testing.T) {
	testingLoadIdentity(t)

	const aid = "1700234872"
	cid, err := FetchVideoCidTryCache(aid)
	if err != nil {
		t.Fatal(err)
	}
	req := Chain{Req: ReqVideoConclusion(aid, cid)}
	err = req.Do()
	if err != nil {
		t.Fatal(err)
	}
	conclusion, err := req.ToVideoConclusion()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(conclusion.DoTemplate())
}

func TestLiveRoomInfo(t *testing.T) {
	testingLoadIdentity(t)

	const roomId = "25788785"
	req := Chain{Req: ReqLiveRoomInfo(roomId)}
	err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	room, err := req.ToLiveRoomInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(room.DoTemplate())
}

func TestLiveStatus(t *testing.T) {
	testingLoadIdentity(t)

	const uid = "1954091502"
	req := Chain{Req: ReqLiveStatus(uid)}
	err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	live, err := req.ToLiveStatusUid()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(live.Get(1954091502).DoTemplate())
}

func TestArticle(t *testing.T) {
	testingLoadIdentity(t)

	req := Chain{Req: ReqArticleInfo("cv41543233")}
	err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	article, err := req.ToArticle()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(article.WithCvid("cv41543233").DoTemplate())
}

func TestSpaceCard(t *testing.T) {
	testingLoadIdentity(t)

	req := Chain{Req: ReqSpaceCard("59442895")}
	err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	card, err := req.ToSpaceCard()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(card.DoTemplate())
}

func TestDynamicDetail(t *testing.T) {
	testingLoadIdentity(t)

	// https://t.bilibili.com/1061910907136770087
	const id = "1061910907136770087"
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(req.Body)
	dynamic, err := req.ToDynamicDetail()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dynamic.Modules.Dynamic.Major.Draw.Id)
	fmt.Println(dynamic.DoTemplate())
}

func TestReplyList(t *testing.T) {
	testingLoadIdentity(t)

	// https://t.bilibili.com/1061910907136770087
	// const id = "1061910907136770087"
	const id = "349150312"
	req := Chain{Req: ReqReplyList(COMMENT_TYPE_ALBUM, id)}
	err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(req.Body)
	reply, err := req.ToReplyList()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(reply)
}

func TestDynamicAll(t *testing.T) {
	testingLoadIdentity(t)

	req := Chain{Req: ReqDynamicAll()}
	err := req.Do()
	if err != nil {
		t.Fatal(err)
	}
	dynamic, err := req.ToDynamicAll()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(dynamic.Items))
	for i, item := range dynamic.Items {
		fmt.Println(i, item.IdStr)
	}
	fmt.Println(dynamic)
}
