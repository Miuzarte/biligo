package biligo

import (
	"encoding/json"
	"io"

	"github.com/tidwall/gjson"
)

type Chain struct {
	Req  *Request
	Resp *Response
	Body string
}

func (c *Chain) Do() (err error) {
	c.Resp, err = httpClient.Do(c.Req)
	if err != nil {
		return err
	}
	defer c.Resp.Response.Body.Close()

	body, err := io.ReadAll(c.Resp.Response.Body)
	if err != nil {
		return err
	}
	c.Body = toString(body)

	if gjson.Get(c.Body, "code").Int() != 0 {
		return wrapErr(ErrChainRespCodeNotZero, c.Body)
	}
	return nil
}

// ParseTo 解析 json 数据到结构体中,
// path 为 json 中的路径, 例如 "data"
func (c *Chain) ParseTo(v any, path ...string) (err error) {
	if c.Resp == nil {
		return wrapErr(ErrChainNilResp, c)
	}
	if c.Body == "" {
		return wrapErr(ErrChainNilBody, c)
	}

	body := c.Body
	for i, p := range path {
		body = gjson.Get(body, p).Raw
		if body == "" {
			return wrapErr(ErrChainPathNotExists, path[:i])
		}
	}

	err = json.Unmarshal([]byte(body), v)
	if err != nil {
		return err
	}
	return nil
}

func (c *Chain) ToVideoInfo() (vi VideoInfo, err error) {
	err = c.ParseTo(&vi, "data")
	return
}

func (c *Chain) ToVideoOnline() (vo VideoOnline, err error) {
	err = c.ParseTo(&vo, "data")
	return
}

func (c *Chain) ToVideoConclusion() (vc VideoConclusion, err error) {
	err = c.ParseTo(&vc, "data")
	return
}

func (c *Chain) ToMediaBase() (mb MediaBase, err error) {
	err = c.ParseTo(&mb, "result")
	return
}

func (c *Chain) ToMedia() (m Media, err error) {
	err = c.ParseTo(&m, "result")
	return
}

func (c *Chain) ToMediaSection() (ms MediaSection, err error) {
	err = c.ParseTo(&ms, "result")
	return
}

func (c *Chain) ToLiveRoomInfo() (lri LiveRoomInfo, err error) {
	err = c.ParseTo(&lri, "data")
	return
}

// 有数据时返回 {..., "data": {"{{uid}}": {...}}}, json 正常解析
//
// 空数据时返回 {..., "data": []}, json 返回错误:
//
//	json: cannot unmarshal array into Go value of type biligo.LiveStatusUid
func (c *Chain) ToLiveStatusUid() (lsu LiveStatusUid, err error) {
	err = c.ParseTo(&lsu, "data")
	return
}

func (c *Chain) ToArticle() (ai ArticleInfo, err error) {
	err = c.ParseTo(&ai, "data")
	return
}

func (c *Chain) ToSongInfo() (si SongInfo, err error) {
	err = c.ParseTo(&si, "data")
	return
}

func (c *Chain) ToSongTag() (st SongTag, err error) {
	err = c.ParseTo(&st, "data")
	return
}

func (c *Chain) ToSongMember() (sm SongMember, err error) {
	err = c.ParseTo(&sm, "data")
	return
}

func (c *Chain) ToSongLyric() (sl SongLyric) {
	return gjson.Get(c.Body, "data").String()
}

func (c *Chain) ToSpaceCard() (sc SpaceCard, err error) {
	err = c.ParseTo(&sc, "data")
	return
}

func (c *Chain) ToVideoPlayurl() (vp VideoPlayurl, err error) {
	err = c.ParseTo(&vp, "data")
	return
}

func (c *Chain) ToReplyList() (rl ReplyList, err error) {
	err = c.ParseTo(&rl, "data")
	return
}

func (c *Chain) ToDynamicSpace() (ds DynamicSpace, err error) {
	err = c.ParseTo(&ds, "data")
	return
}

func (c *Chain) ToDynamicSpaceDesktop() (dsd DynamicSpaceDesktop, err error) {
	err = c.ParseTo(&dsd, "data")
	return
}

func (c *Chain) ToDynamicAll() (da DynamicAll, err error) {
	err = c.ParseTo(&da, "data")
	return
}

func (c *Chain) ToDynamicAllUpdate() (dau DynamicAllUpdate, err error) {
	err = c.ParseTo(&dau, "data")
	return
}

func (c *Chain) ToDynamicDetail() (dd DynamicDetail, err error) {
	err = c.ParseTo(&dd, "data.item")
	return
}

func (c *Chain) ToDynamicDetailDesktop() (dd DynamicDetailDesktop, err error) {
	err = c.ParseTo(&dd, "data")
	return
}

func (c *Chain) ToVoteInfo() (vi VoteInfo, err error) {
	err = c.ParseTo(&vi, "data.vote_info")
	return
}

func (c *Chain) ToQrcodeGenerate() (qg QrcodeGenerate, err error) {
	err = c.ParseTo(&qg, "data")
	return
}

func (c *Chain) ToQrcodePoll() (qp QrcodePoll, err error) {
	err = c.ParseTo(&qp, "data")
	return
}

func (c *Chain) ToRelationStat() (rs RelationStat, err error) {
	err = c.ParseTo(&rs, "data")
	return
}

func (c *Chain) ToLiveDanmuInfo() (ldi LiveDanmuInfo, err error) {
	err = c.ParseTo(&ldi, "data")
	return
}

func (c *Chain) ToNav() (n Nav, err error) {
	err = c.ParseTo(&n, "data")
	return
}

func (c *Chain) ToBuvid34() (b Buvid34, err error) {
	err = c.ParseTo(&b, "data")
	return
}
