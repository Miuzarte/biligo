package biligo

const (
	// 视频详细信息 [VideoInfo]
	//	.WithQuery("aid", aid)
	//	.WithQuery("bvid", bvid)
	URL_VIDEO_INFO_WBI = `https://api.bilibili.com/x/web-interface/wbi/view`
	URL_VIDEO_INFO     = `https://api.bilibili.com/x/web-interface/view`
	// 视频在线人数 [VideoOnline]
	//	.WithQuerys("aid", aid, "cid", cid)
	URL_VIDEO_ONLINE = `https://api.bilibili.com/x/player/online/total`
	// 视频AI总结内容 [VideoConclusion]
	//	.WithQuerys("aid", aid, "cid", cid)
	//	?.WithQuery("up_mid", uid) // 似乎非必须
	URL_VIDEO_CONCLUSION_WBI = `https://api.bilibili.com/x/web-interface/view/conclusion/get`

	// 剧集基本信息, 返回的结果字段为 `result`, 信息很少 [MediaBase]
	//	.WithQuery("media_id", mdid)
	URL_MEDIA_INFO_BASE = `https://api.bilibili.com/pgc/review/user`
	// 剧集明细, 返回的结果字段为 `result`, 两种 query 返回的结果一致 (仅测试了番剧) [Media]
	//	.WithQuery("season_id", ssid)
	//	.WithQuery("ep_id", epid)
	URL_MEDIA_INFO_DETAIL = `https://api.bilibili.com/pgc/view/web/season`
	// 剧集分集信息, 返回的结果字段为 `result` [MediaSection]
	//	.WithQuery("season_id", ssid)
	URL_MEDIA_SECTION = `https://api.bilibili.com/pgc/web/season/section`

	// 直播间状态, 可批量 [LiveStatus]
	//	.WithQueryList("uids[]", uids...)
	URL_LIVE_STATUS_UIDS = `https://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids`
	// 直播间信息, 拿不到用户名 [LiveRoomInfo]
	//	.WithQuery("room_id", roomId)
	URL_LIVE_INFO_ROOMID = `https://api.live.bilibili.com/room/v1/Room/get_info`

	// 专栏文章基本信息 [ArticleInfo]
	//	.WithQuery("id", id) // TrimPrefix(id, "cv")
	URL_ARTICLE_INFO = `https://api.bilibili.com/x/article/viewinfo`

	// 整合: [Song]
	// 歌曲基本信息 [SongInfo]
	//	.WithQuery("sid", auid) // TrimPrefix(auid, "au")
	URL_SONG_INFO = `https://www.bilibili.com/audio/music-service-c/web/song/info`
	// 歌曲TAG [SongTag]
	//	.WithQuery("sid", auid) // TrimPrefix(auid, "au")
	URL_SONG_TAG = `https://www.bilibili.com/audio/music-service-c/web/tag/song`
	// 歌曲创作成员列表 [SongMember]
	//	.WithQuery("sid", auid) // TrimPrefix(auid, "au")
	URL_SONG_MEMBER = `https://www.bilibili.com/audio/music-service-c/web/member/song`
	// 歌曲歌词 [SongLyric]
	//	.WithQuery("sid", auid) // TrimPrefix(auid, "au")
	URL_SONG_LIRIC = `https://www.bilibili.com/audio/music-service-c/web/song/lyric`

	// 用户名片信息 [SpaceCard]
	//	.WithQuery("mid", uid)
	//	.WithQuery("photo", "true") // 请求用户主页头图 space.s_img, space.l_img
	URL_SPACE_CARD = `https://api.bilibili.com/x/web-interface/card`

	// 用户关系状态数 [RelationStat]
	//	.WithQuery("vmid", uid)
	URL_RELATION_STAT = `https://api.bilibili.com/x/relation/stat`

	// 全部动态列表 [DynamicAll]
	URL_DYNAMIC_ALL = `https://api.bilibili.com/x/polymer/web-dynamic/v1/feed/all`
	// 新动态数量 [DynamicAllUpdate]
	//	.WithQuery("update_baseline", updateBaseline)
	URL_DYNAMIC_ALL_UPDATE = `https://api.bilibili.com/x/polymer/web-dynamic/v1/feed/all/update`

	// 动态详情 [DynamicDetail]
	//	.WithQuery("id", id)
	// 图文动态无文字内容
	URL_DYNAMIC_DETAIL = `https://api.bilibili.com/x/polymer/web-dynamic/v1/detail`
	// 动态详情 [DynamicDetailDesktop]
	//	.WithQuery("id", id)
	// 二次请求补充文字内容
	URL_DYNAMIC_DETAIL_DESKTOP = `https://api.bilibili.com/x/polymer/web-dynamic/desktop/v1/detail`

	// 投票信息 [VoteInfo]
	//	.WithQuery("vote_id", voteId)
	URL_VOTE_INFO = `https://api.bilibili.com/x/vote/vote_info`
	// 投票信息(无图片)
	//	.WithQuery("vote_id", voteId)
	URL_VOTE_INFO_ = `https://api.vc.bilibili.com/vote_svr/v1/vote_svr/vote_info`

	// 综合搜索 [SearchAll]
	//	.WithQuery("keyword", keyword)
	URL_SEARCH_ALL_WBI = `https://api.bilibili.com/x/web-interface/wbi/search/all/v2`
	// 综合搜索(旧链接)
	URL_SEARCH_ALL_OLD = `https://api.bilibili.com/x/web-interface/search/all/v2`
	// 分类搜索 [SearchType] [SearchTypeLive]
	//	.WithQuerys("search_type", string([SearchClass]), "keyword", keyword)
	URL_SEARCH_TYPE_WBI = `https://api.bilibili.com/x/web-interface/wbi/search/type`
	// 分类搜索(旧链接)
	URL_SEARCH_TYPE_OLD = `https://api.bilibili.com/x/web-interface/search/type`

	// 直播间信息流认证秘钥 [LiveDanmuInfo]
	//	.WithQuery("id", roomId) // 直播间真实id
	URL_LIVE_DANMU_INFO = `https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo`

	// 视频流地址 [VideoPlayurl]
	//	.WithQuerys(
	//		"avid", aid, "cid", cid, "fnval", fnval, // [VIDED_FNVAL_DASHALL] / [VIDEO_FNVAL_MP4]
	//		"fourk", "1", // 请求 4K
	//		"try_look", "1", // 游客高清晰度
	//	)
	URL_VIDEO_PALYURL_WBI = `https://api.bilibili.com/x/player/wbi/playurl`
	// 视频流地址(旧链接)
	URL_VIDEO_PALYURL_OLD = `https://api.bilibili.com/x/player/playurl`

	// 评论区评论列表 [ReplyList]
	//	.WithQuerys("type", itoa([CommentClass]), "oid", id)
	URL_REPLY_LIST = `https://api.bilibili.com/x/v2/reply`

	// 申请登录二维码 [QrcodeGenerate]
	URL_LOGIN_QRCODE_GENERATE = `https://passport.bilibili.com/x/passport-login/web/qrcode/generate`
	// 扫码登录 [QrcodePoll]
	//	.WithQuery("qrcode_key", qrcodeKey)
	URL_LOGIN_QRCODE_POLL = `https://passport.bilibili.com/x/passport-login/web/qrcode/poll`

	// 导航栏用户信息 主要用于获取 wbi keys [Nav]
	// NOTE: 未登录时也能正常获取到 wbi keys,
	// 但是此时响应状态码不为 0
	URL_NAV = `https://api.bilibili.com/x/web-interface/nav`

	// URL_BUVID3  = `https://api.bilibili.com/x/web-frontend/getbuvid`
	URL_BUVID34 = `https://api.bilibili.com/x/frontend/finger/spi`
)
