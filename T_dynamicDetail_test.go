package biligo

import (
	"fmt"
	"testing"
)

func TestVoteInfo(t *testing.T) {
	const id = `14349994`
	req := Chain{Req: ReqVoteInfo(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestVoteInfo_(t *testing.T) {
	const id = `14349994`
	req := Chain{Req: NewGet(URL_VOTE_INFO_).WithQuery("vote_id", id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailBlocked(t *testing.T) {
	// https://t.bilibili.com/1063386760933801988
	const id = `1063386760933801988`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

// 图文动态, 带图投票
func TestDynamicDetailDrawVote(t *testing.T) {
	// https://t.bilibili.com/1027118970485866514
	const id = `1027118970485866514`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

// 带图投票
func TestDynamicDetailVotePic(t *testing.T) {
	const id = `14349994`
	req := Chain{Req: ReqVoteInfo(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailDynamicVideo(t *testing.T) {
	// https://t.bilibili.com/1053778618491076616
	const id = `1053778618491076616`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailTextForwardVideo(t *testing.T) {
	// https://t.bilibili.com/1051224160886325257
	const id = `1051224160886325257`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailReserveVideo(t *testing.T) {
	// https://t.bilibili.com/1063755995030749203
	const id = `1063755995030749203`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailReserveLive(t *testing.T) {
	// https://t.bilibili.com/1063825951430803465
	const id = `1063825951430803465`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailRelatedGame(t *testing.T) {
	// https://t.bilibili.com/1063715738261389314
	const id = `1063715738261389314`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailArticle(t *testing.T) {
	// https://t.bilibili.com/1063713848460050435
	const id = `1063713848460050435`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

// 伪专栏, 获取不到标题
func TestDynamicDetailArticleSmall(t *testing.T) {
	// https://t.bilibili.com/1025375737754943528
	const id = `1025375737754943528`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailUgc(t *testing.T) {
	// https://t.bilibili.com/1063736375606509624
	const id = `1063736375606509624`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailForwardDelete(t *testing.T) {
	// https://t.bilibili.com/1063868488806825987
	const id = `1063868488806825987`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailTopic(t *testing.T) {
	// https://t.bilibili.com/1031434283062919170
	const id = `1031434283062919170`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailAudio(t *testing.T) {
	// https://t.bilibili.com/505035128347476154
	const id = `505035128347476154`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailLive(t *testing.T) {
	// https://t.bilibili.com/1063877443815735313
	const id = `1063877443815735313`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailLiveRcmd(t *testing.T) {
	// https://t.bilibili.com/1063837028173479941
	const id = `1063837028173479941`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}

func TestDynamicDetailBangumi(t *testing.T) {
	// https://t.bilibili.com/1062942566841843721
	const id = `1062942566841843721`
	req := Chain{Req: ReqDynamicDetail(id)}
	err := req.Do()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)
}
