package biligo

const (
	URL_DOMAIN    = `bilibili.com`
	URL_MAIN_PAGE = `https://www.bilibili.com`
)

const (
	TIME_LAYOUT_L24  = `2006/01/02 15:04:05`
	TIME_LAYOUT_L24C = `2006年01月02日15时04分05秒`
	TIME_LAYOUT_M24  = `01/02 15:04:05`
	TIME_LAYOUT_M24C = `01月02日15时04分05秒`
	TIME_LAYOUT_S24  = `02 15:04:05`
	TIME_LAYOUT_S24C = `02日15时04分05秒`
	TIME_LAYOUT_T24  = `15:04:05`
	TIME_LAYOUT_T24C = `15时04分05秒`
)

const (
	REGEXP_SHORT   = `(?:b23\.tv|acg\.tv|bili2233\.cn)/([0-9A-Za-z]+)` // 0: url, 1: id(av/BV/shortId)
	REGEXP_VIDEO   = `video/(BV1[1-9A-HJ-NP-Za-km-z]{9}|av[0-9]+)`     // 1: id(av/BV)
	REGEXP_BANGUMI = `bangumi/(?:media|play)/((md|ss|ep)[0-9]+)`       // 1: id(md/ss/dp), 2: "md"|"ss"|"ep"
	REGEXP_LIVE    = `live\.bilibili\.com/([0-9]+)`                    // 1: room id
	REGEXP_ARTICLE = `(?:read/cv|read/mobile/)([0-9]+)`                // 1: cvid(numbers)
	REGEXP_AUDIO   = `audio/au([0-9]+)`                                // 1: audio id
	REGEXP_SPACE   = `(?:space|m)\.bilibili\.com/(?:space/)?([0-9]+)`  // 1: space id
	REGEXP_DYNAMIC = `(?:t.bilibili.com|dynamic|opus)/([0-9]+)`        // 1: dynamic id
)

const (
	REGEXP_INDEX_SHORT_URL, REGEXP_INDEX_SHORT_ID         = 0, 1
	REGEXP_INDEX_ARCHIVE                                  = 1
	REGEXP_INDEX_BANGUMI_ID, REGEXP_INDEX_BANGUMI_ID_TYPE = 1, 2
	REGEXP_INDEX_LIVE                                     = 1
	REGEXP_INDEX_ARTICLE                                  = 1
	REGEXP_INDEX_AUDIO                                    = 1
	REGEXP_INDEX_SPACE                                    = 1
	REGEXP_INDEX_DYNAMIC                                  = 1
)

type LinkType int

const (
	_ LinkType = iota
	LINK_TYPE_SHORT
	LINK_TYPE_ARCHIVE
	LINK_TYPE_MEDIA
	LINK_TYPE_LIVE
	LINK_TYPE_ARTICLE
	LINK_TYPE_SONG
	LINK_TYPE_SPACE
	LINK_TYPE_DYNAMIC
)

func (lt LinkType) String() string {
	switch lt {
	case LINK_TYPE_SHORT:
		return "short"
	case LINK_TYPE_ARCHIVE:
		return "video"
	case LINK_TYPE_MEDIA:
		return "media"
	case LINK_TYPE_LIVE:
		return "live"
	case LINK_TYPE_ARTICLE:
		return "article"
	case LINK_TYPE_SONG:
		return "audio"
	case LINK_TYPE_SPACE:
		return "space"
	case LINK_TYPE_DYNAMIC:
		return "dynamic"
	default:
		return "unknown"
	}
}

type SearchClass string

const (
	SEARCH_TYPE_VIDEO   SearchClass = "video"
	SEARCH_TYPE_VIDEO_C string      = "视频"

	SEARCH_TYPE_BANGUMI   SearchClass = "media_bangumi"
	SEARCH_TYPE_BANGUMI_C string      = "番剧"
	SEARCH_TYPE_FT        SearchClass = "media_ft"
	SEARCH_TYPE_FT_C      string      = "影视"

	SEARCH_TYPE_LIVE        SearchClass = "live" // 包括直播间和主播, 实际结果中不存在这个类型, 仅在请求中使用
	SEARCH_TYPE_LIVE_C      string      = "直播"
	SEARCH_TYPE_LIVE_USER   SearchClass = "live_user"
	SEARCH_TYPE_LIVE_USER_C string      = "主播"
	SEARCH_TYPE_LIVE_ROOM   SearchClass = "live_room"
	SEARCH_TYPE_LIVE_ROOM_C string      = "直播间"

	SEARCH_TYPE_ARTICLE   SearchClass = "article"
	SEARCH_TYPE_ARTICLE_C string      = "专栏"

	SEARCH_TYPE_BILI_USER   SearchClass = "bili_user"
	SEARCH_TYPE_BILI_USER_C string      = "用户"

	SEARCH_TYPE_TOPIC   SearchClass = "topic" // 老东西 搜不出啥
	SEARCH_TYPE_TOPIC_C string      = "话题"

	SEARCH_TYPE_PHOTO   SearchClass = "photo" // 搜不到东西
	SEARCH_TYPE_PHOTO_C string      = "相簿"
)

func (st SearchClass) String() string {
	return SearchTypeMap[st]
}

var (
	SearchTypeMap = map[SearchClass]string{
		SEARCH_TYPE_VIDEO:     SEARCH_TYPE_VIDEO_C,
		SEARCH_TYPE_BANGUMI:   SEARCH_TYPE_BANGUMI_C,
		SEARCH_TYPE_FT:        SEARCH_TYPE_FT_C,
		SEARCH_TYPE_LIVE:      SEARCH_TYPE_LIVE_C,
		SEARCH_TYPE_LIVE_ROOM: SEARCH_TYPE_LIVE_ROOM_C,
		SEARCH_TYPE_LIVE_USER: SEARCH_TYPE_LIVE_USER_C,
		SEARCH_TYPE_ARTICLE:   SEARCH_TYPE_ARTICLE_C,
		SEARCH_TYPE_BILI_USER: SEARCH_TYPE_BILI_USER_C,

		SEARCH_TYPE_TOPIC: SEARCH_TYPE_TOPIC_C,
		SEARCH_TYPE_PHOTO: SEARCH_TYPE_PHOTO_C,
	}
	SearchTypePam = map[string]SearchClass{
		SEARCH_TYPE_VIDEO_C:     SEARCH_TYPE_VIDEO,
		SEARCH_TYPE_BANGUMI_C:   SEARCH_TYPE_BANGUMI,
		SEARCH_TYPE_FT_C:        SEARCH_TYPE_FT,
		SEARCH_TYPE_LIVE_C:      SEARCH_TYPE_LIVE,
		SEARCH_TYPE_LIVE_ROOM_C: SEARCH_TYPE_LIVE_ROOM,
		SEARCH_TYPE_LIVE_USER_C: SEARCH_TYPE_LIVE_USER,
		SEARCH_TYPE_ARTICLE_C:   SEARCH_TYPE_ARTICLE,
		SEARCH_TYPE_BILI_USER_C: SEARCH_TYPE_BILI_USER,

		SEARCH_TYPE_TOPIC_C: SEARCH_TYPE_TOPIC,
		SEARCH_TYPE_PHOTO_C: SEARCH_TYPE_PHOTO,
	}
)

type LoginCodeState int

const (
	LOGIN_CODE_STATE_SUCCESS   LoginCodeState = 0     // 扫码登录成功
	LOGIN_CODE_STATE_EXPIRED   LoginCodeState = 86038 // 二维码已失效
	LOGIN_CODE_STATE_SCANNED   LoginCodeState = 86090 // 二维码已扫码未确认
	LOGIN_CODE_STATE_UNSCANNED LoginCodeState = 86101 // 未扫码
)

func (lcs LoginCodeState) String() string {
	switch lcs {
	case LOGIN_CODE_STATE_SUCCESS:
		return "success"
	case LOGIN_CODE_STATE_EXPIRED:
		return "expired"
	case LOGIN_CODE_STATE_SCANNED:
		return "scaned"
	case LOGIN_CODE_STATE_UNSCANNED:
		return "unscaned"
	default:
		return "unknown"
	}
}

type LiveMsgCmd = string

const (
	LIVE_MSG_STREAM_LIVE      LiveMsgCmd = "LIVE"        // 直播开始
	LIVE_MSG_STREAM_PREPARING LiveMsgCmd = "PREPARING"   // 直播准备中 (结束)
	LIVE_MSG_STREAM_CHANGE    LiveMsgCmd = "ROOM_CHANGE" // 房间信息变更
	LIVE_MSG_STREAM_WARNING   LiveMsgCmd = "WARNING"     // 警告
	LIVE_MSG_STREAM_CUT_OFF   LiveMsgCmd = "CUT_OFF"     // 切断
)

type CommentClass = int

const (
	_ CommentClass = iota

	COMMENT_TYPE_VIDEO // 视频稿件 (oid: 稿件 avid)
	COMMENT_TYPE_TOPIC // 话题 (oid: 话题 id)
	_
	COMMENT_TYPE_ACTIVITY          // 活动 (oid: 活动 id)
	COMMENT_TYPE_SHORT_VIDEO       // 小视频 (oid: 小视频 id)
	COMMENT_TYPE_BLACK_ROOM        // 小黑屋封禁信息 (oid: 封禁公示 id)
	COMMENT_TYPE_ANNOUNCEMENT      // 公告信息 (oid: 公告 id)
	COMMENT_TYPE_LIVE_ACTIVITY     // 直播活动 (oid: 直播间 id)
	COMMENT_TYPE_ACTIVITY_DRAFT    // 活动稿件 (oid: ?)
	COMMENT_TYPE_LIVE_ANNOUNCEMENT // 直播公告 (oid: ?)
	COMMENT_TYPE_ALBUM             // 相簿（图片动态） (oid: 相簿 id)
	COMMENT_TYPE_COLUMN            // 专栏 (oid: 专栏 cvid)
	COMMENT_TYPE_TICKET            // 票务 (oid: ?)
	COMMENT_TYPE_AUDIO             // 音频 (oid: 音频 auid)
	COMMENT_TYPE_JUDGMENT          // 风纪委员会 (oid: 众裁项目 id)
	COMMENT_TYPE_REVIEW            // 点评 (oid: ?)
	COMMENT_TYPE_DYNAMIC           // 动态（纯文字动态&分享） (oid: 动态 id)
	COMMENT_TYPE_PLAYLIST          // 播单 (oid: ?)
	COMMENT_TYPE_MUSIC_PLAYLIST    // 音乐播单 (oid: ?)
	COMMENT_TYPE_COMIC1            // 漫画 (oid: ?)
	COMMENT_TYPE_COMIC2            // 漫画 (oid: ?)
	COMMENT_TYPE_COMIC3            // 漫画 (oid: 漫画 mcid)

	COMMENT_TYPE_COURSE CommentClass = 33 // 课程 (oid: 课程 epid)
)

type VideoQuality = int

const (
	VIDEO_QN_240      VideoQuality = 6
	VIDEO_QN_360      VideoQuality = 16
	VIDEO_QN_480      VideoQuality = 32
	VIDEO_QN_720      VideoQuality = 64
	VIDEO_QN_720P60   VideoQuality = 74
	VIDEO_QN_1080     VideoQuality = 80
	VIDEO_QN_1080PLUS VideoQuality = 112
	VIDEO_QN_1080P60  VideoQuality = 116
	VIDEO_QN_4K       VideoQuality = 120
	VIDEO_QN_HDR      VideoQuality = 125
	VIDEO_QN_DOLBY    VideoQuality = 126
	VIDEO_QN_8K       VideoQuality = 127
)

type VideoFormat = int

const (
	VIDEO_FNVAL_FLV VideoFormat = 0
	VIDEO_FNVAL_MP4 VideoFormat = 1

	VIDEO_FNVAL_DASH        VideoFormat = 16
	VIDEO_FNVAL_HDR         VideoFormat = 64
	VIDEO_FNVAL_4K          VideoFormat = 128
	VIDEO_FNVAL_DOLBYAUDIO  VideoFormat = 256
	VIDEO_FNVAL_DOLBYVISION VideoFormat = 512
	VIDEO_FNVAL_8K          VideoFormat = 1024
	VIDEO_FNVAL_AV1         VideoFormat = 2048

	VIDED_FNVAL_DASHALL VideoFormat = // 4048
	VIDEO_FNVAL_DASH |
		VIDEO_FNVAL_HDR |
		VIDEO_FNVAL_4K |
		VIDEO_FNVAL_DOLBYAUDIO |
		VIDEO_FNVAL_DOLBYVISION |
		VIDEO_FNVAL_8K |
		VIDEO_FNVAL_AV1
)

type VideoCodec = string

const (
	// [VIDEO_QN_4K] "avc1.640034"
	// [VIDEO_QN_1080] [VIDEO_QN_1080PLUS] [VIDEO_QN_1080P60] "avc1.640032"
	// [VIDEO_QN_720] "avc1.640028"
	// [VIDEO_QN_480] "avc1.64001F"
	// [VIDEO_QN_360] "avc1.64001E"
	VIDEO_CODEC_AVC  VideoCodec = "avc1.6400**"
	VIDEO_CODEC_HEVC VideoCodec = "hev1.1.6.L153.90"
	VIDEO_CODEC_AV1  VideoCodec = "av01.0.00M.10.0.110.01.01.01.0"
)

type VideoCodecId = int

const (
	VIDEO_CODEC_ID_AVC  VideoCodecId = 7  // H.264
	VIDEO_CODEC_ID_HEVC VideoCodecId = 12 // H.265
	VIDEO_CODEC_ID_AV1  VideoCodecId = 13 // AV1
)

const (
	// 30232 "mp4a.40.2"
	// 39489 "mp4a.40.5"
	// 30280 "mp4a.40.2"
	AUDIO_CODEC_M4A VideoCodec = "mp4a.40.*"
)
