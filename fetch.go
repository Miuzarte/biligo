package biligo

import (
	"net/http"
	"sync"

	"github.com/tidwall/gjson"
)

func fetchHeader(link string) (header http.Header, err error) {
	resp, err := httpClient.Do(NewHead(link))
	if err != nil {
		return
	}
	return resp.Header, nil
}

// fetchLocation 获取链接重定向地址
func fetchLocation(link string) (location string, err error) {
	header, err := fetchHeader(link)
	if err != nil {
		return "", err
	}
	loc := header.Get("Location")
	if loc == "" {
		return "", wrapErr(ErrFetchNoLocation, header)
	}
	return loc, nil
}

// fetchHelperWithResp 执行请求并返回结果 (data 字段) 和 http 响应
func fetchHelperWithResp(req *Request, path ...string) (j gjson.Result, resp *Response, err error) {
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	j, err = resp.ToGjson()
	if err != nil {
		return
	}
	if j.Get("code").Int() != 0 {
		err = wrapErr(ErrFetchRespCodeNotZero, j)
		return
	}

	k := j
	for i, p := range path {
		k = k.Get(p)
		if !k.Exists() {
			return j, resp, wrapErr(ErrFetchPathNotExists, path[:i+1])
		}
	}

	return k, resp, nil
}

// FetchVideoInfo 获取视频信息
//
//	id: avid/bvid
func FetchVideoInfo(id string) (vi VideoInfo, err error) {
	aid, err := AnyToAid(id)
	if err != nil {
		return
	}

	req := Chain{Req: ReqVideoInfo(aid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToVideoInfo()
}

func FetchVideoCid(id string) (cid string, err error) {
	vi, err := FetchVideoInfo(id)
	if err != nil {
		return
	}
	if vi.Cid == 0 {
		err = wrapErr(ErrFetchCidNotExists, vi)
		return
	}
	return itoa(vi.Cid), nil
}

var (
	// 池中只会有一个对象,
	// mutex 被持有时表示池中对象正在被使用
	cidCachePool = sync.Pool{
		New: func() any {
			return make(map[string]string)
		},
	}
	cidCachePoolMu sync.Mutex
)

// FetchVideoCidTryCache 尝试从缓存中获取 cid,
// 如果没有则从网络获取并缓存
func FetchVideoCidTryCache(id string) (cid string, err error) {
	id, err = AnyToAid(id)
	if err != nil {
		return
	}
	if !cidCachePoolMu.TryLock() { // 池被占用
		return FetchVideoCid(id)
	}
	defer cidCachePoolMu.Unlock()

	m := cidCachePool.Get().(map[string]string)
	defer cidCachePool.Put(m)

	cid, ok := m[id]
	if ok {
		return
	}

	cid, err = FetchVideoCid(id)
	if err == nil {
		m[id] = cid
	}
	return
}

// FetchVideoOnline 获取视频在线人数
//
//	id: avid/bvid
//	cid: 分 p 对应的 cid
func FetchVideoOnline(id, cid string) (vo VideoOnline, err error) {
	aid, err := AnyToAid(id)
	if err != nil {
		return
	}
	if cid == "" || cid == "0" {
		cid, err = FetchVideoCidTryCache(id)
		if err != nil {
			return
		}
	}

	req := Chain{Req: ReqVideoOnline(aid, cid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToVideoOnline()
}

// FetchVideoConclusion 获取视频AI总结
//
//	id: avid/bvid
//	cid: 分 p 对应的 cid
func FetchVideoConclusion(id, cid string) (vc VideoConclusion, err error) {
	aid, err := AnyToAid(id)
	if err != nil {
		return
	}
	if cid == "" || cid == "0" {
		cid, err = FetchVideoCidTryCache(id)
		if err != nil {
			return
		}
	}

	req := Chain{Req: ReqVideoConclusion(aid, cid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToVideoConclusion()
}

// FetchVideoInfo2 获取视频信息、视频在线人数
//
//	id: avid/bvid
func FetchVideoInfo2(id string) (vi VideoInfo, vo VideoOnline, err error) {
	vi, err = FetchVideoInfo(id)
	if err != nil {
		return
	}
	vo, err = FetchVideoOnline(id, itoa(vi.Cid))
	if err != nil {
		return
	}
	vi.WithOnline(vo)
	return
}

// FetchVideoInfo3 获取视频信息、视频在线人数、总结
//
//	id: avid/bvid
func FetchVideoInfo3(id string) (vi VideoInfo, vo VideoOnline, vc VideoConclusion, err error) {
	vi, err = FetchVideoInfo(id)
	if err != nil {
		return
	}
	cid := itoa(vi.Cid)

	var oErr, cErr error
	await(
		func() { vo, oErr = FetchVideoOnline(id, cid) },
		func() { vc, cErr = FetchVideoConclusion(id, cid) },
	)
	if oErr != nil {
		err = oErr
		return
	}
	vi.WithOnline(vo)
	if cErr != nil {
		err = cErr
		return
	}
	return
}

// FetchMediaInfoBase 获取剧集基本信息
func FetchMediaInfoBase(mdid string) (mb MediaBase, err error) {
	req := Chain{Req: ReqMediaInfoBase(mdid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToMediaBase()
}

// FetchMediaInfoSsid 获取剧集明细
func FetchMediaInfoSsid(ssid string) (m Media, err error) {
	req := Chain{Req: ReqMediaInfoSsid(ssid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToMedia()
}

// FetchMediaInfoEpid 获取剧集明细
func FetchMediaInfoEpid(epid string) (m Media, err error) {
	req := Chain{Req: ReqMediaInfoEpid(epid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToMedia()
}

// FetchMediaSection 获取剧集分集信息
func FetchMediaSection(ssid string) (ms MediaSection, err error) {
	req := Chain{Req: ReqMediaSection(ssid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToMediaSection()
}

// FetchLiveStatus 批量获取直播间状态
//
// 有数据时返回 {..., "data": {"{{uid}}": {...}}}, json 正常解析
//
// 空数据时返回 {..., "data": []}, json 返回错误:
//
//	json: cannot unmarshal array into Go value of type biligo.LiveStatusUid
func FetchLiveStatus(uid ...string) (lsu LiveStatusUid, err error) {
	req := Chain{Req: ReqLiveStatus(uid...)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToLiveStatusUid()
}

// FetchLiveRoomInfo 获取直播间信息, 拿不到用户名
func FetchLiveRoomInfo(roomId string) (lri LiveRoomInfo, err error) {
	req := Chain{Req: ReqLiveRoomInfo(roomId)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToLiveRoomInfo()
}

// FetchArticleInfo 获取专栏文章基本信息
func FetchArticleInfo(id string) (ai ArticleInfo, err error) {
	req := Chain{Req: ReqArticleInfo(id)}
	err = req.Do()
	if err != nil {
		return
	}
	ai, err = req.ToArticle()
	if err != nil {
		return
	}
	ai.WithCvid(id) // 附加 cvid 用于格式化
	return
}

// FetchSong 获取音乐信息
func FetchSong(auid string) (song Song, err error) {
	i, t, m, l := ReqAudio(auid)
	iReq, tReq, mReq, lReq := Chain{Req: i}, Chain{Req: t}, Chain{Req: m}, Chain{Req: l}
	var iErr, tErr, mErr, lErr error

	await(
		func() { iErr = iReq.Do() },
		func() { tErr = tReq.Do() },
		func() { mErr = mReq.Do() },
		func() { lErr = lReq.Do() },
	)
	if iErr != nil {
		err = iErr
		return
	}
	if tErr != nil {
		err = tErr
		return
	}
	if mErr != nil {
		err = mErr
		return
	}
	if lErr != nil {
		err = lErr
		return
	}

	await(
		func() { song.Info, iErr = iReq.ToSongInfo() },
		func() { song.Tag, tErr = tReq.ToSongTag() },
		func() { song.Member, mErr = mReq.ToSongMember() },
		func() { song.Lyric = lReq.ToSongLyric() },
	)
	if iErr != nil {
		err = iErr
		return
	}
	if tErr != nil {
		err = tErr
		return
	}
	if mErr != nil {
		err = mErr
		return
	}

	return
}

// FetchSpaceCard 获取用户名片信息
func FetchSpaceCard(uid string) (sc SpaceCard, err error) {
	req := Chain{Req: ReqSpaceCard(uid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToSpaceCard()
}

// FetchRelationStat 获取关注/粉丝数
func FetchRelationStat(uid string) (rs RelationStat, err error) {
	req := Chain{Req: ReqRelationStat(uid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToRelationStat()
}

// FetchDynamicAll 获取全部动态列表
func FetchDynamicAll() (da DynamicAll, err error) {
	req := Chain{Req: ReqDynamicAll()}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToDynamicAll()
}

// FetchDynamicAllUpdate 获取新动态数量
func FetchDynamicAllUpdate(updateBaseline string) (dau DynamicAllUpdate, err error) {
	req := Chain{Req: ReqDynamicAllUpdate(updateBaseline)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToDynamicAllUpdate()
}

// FetchVideoPlayurl 获取视频流地址,
// 不提供 cid 时默认获取 P1,
// 取 dash 流时不需要 qn,
// 取 mp4 流时不传 qn 获取的清晰度应该是 720p, 未经大量测试,
// 同时, B站引入 dash 后较新的视频 mp4 只有 720p 与 360p
func FetchVideoPlayurl(id, cid string, fnval VideoFormat, qn ...VideoQuality) (vp VideoPlayurl, err error) {
	aid, err := AnyToAid(id)
	if err != nil {
		return
	}
	if cid == "" || cid == "0" {
		cid, err = FetchVideoCidTryCache(id)
		if err != nil {
			return
		}
	}

	req := Chain{Req: ReqVideoPlayurl().WithQuerys(
		"avid", aid, "cid", cid, "fnval", itoa(fnval),
		"fourk", "1", // 请求 4K
		"try_look", "1", // 游客高清晰度
	)}
	if len(qn) != 0 { // dash 不需要
		req.Req.WithQueryList("qn", itoa(qn[0]))
	}

	err = req.Do()
	if err != nil {
		return
	}
	return req.ToVideoPlayurl()
}

// FetchDynamicDetail 获取动态信息
func FetchDynamicDetail(id string) (dd DynamicDetail, err error) {
	req := Chain{Req: ReqDynamicDetail(id)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToDynamicDetail()
}

// FetchVoteInfo 获取投票信息
func FetchVoteInfo(voteId string) (vi VoteInfo, err error) {
	req := Chain{Req: ReqVoteInfo(voteId)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToVoteInfo()
}

// FetchSearchAll 获取综合搜索
func FetchSearchAll(keyword string) (j gjson.Result, err error) {
	j, _, err = fetchHelperWithResp(ReqSearchAll(keyword), "data.result")
	return
}

// FetchSearchType 获取分类搜索
func FetchSearchType(searchType SearchClass, keyword string) (j gjson.Result, err error) {
	j, _, err = fetchHelperWithResp(ReqSearchType(searchType, keyword), "data.result")
	return
}

// FetchLiveDanmuInfo 获取信息流认证秘钥
func FetchLiveDanmuInfo(roomId string) (ldi LiveDanmuInfo, err error) {
	req := Chain{Req: ReqLiveDanmuInfo(roomId)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToLiveDanmuInfo()
}

// FetchReplyList 获取评论列表
func FetchReplyList(typ CommentClass, oid string) (rl ReplyList, err error) {
	req := Chain{Req: ReqReplyList(typ, oid)}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToReplyList()
}

// FetchLoginQrcodeGenerate 申请登录二维码
func FetchLoginQrcodeGenerate() (qg QrcodeGenerate, err error) {
	req := Chain{Req: ReqLoginQrcodeGenerate()}
	err = req.Do()
	if err != nil {
		return
	}
	return req.ToQrcodeGenerate()
}

// FetchLoginQrcodePoll 扫码登录轮询
func FetchLoginQrcodePoll(qrcodeKey string) (qp QrcodePoll, header http.Header, err error) {
	req := Chain{Req: ReqLoginQrcodePoll(qrcodeKey)}
	err = req.Do()
	if req.Resp != nil {
		header = req.Resp.Header
	}
	if err != nil {
		return
	}
	qp, err = req.ToQrcodePoll()
	return
}

// FetchNav 获取导航栏用户信息
func FetchNav() (nav Nav, err error) {
	req := Chain{Req: ReqNav()}
	err = req.Do()
	if err != nil {
		// 未登录时响应状态码不为 0, 会返回 err
		if !err.(*Error).Is(ErrChainRespCodeNotZero) {
			return
		}
	}
	return req.ToNav()
}
