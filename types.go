package biligo

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	videoDimension struct {
		Width  int `json:"width" mapstructure:"width"`
		Height int `json:"height" mapstructure:"height"`
		Rotate int `json:"rotate" mapstructure:"rotate"`
	}
	videoPage struct {
		Cid        int            `json:"cid" mapstructure:"cid"`
		Page       int            `json:"page" mapstructure:"page"`
		From       string         `json:"from" mapstructure:"from"` // "vupload"
		Part       string         `json:"part" mapstructure:"part"`
		Duration   int            `json:"duration" mapstructure:"duration"` // (s)
		Vid        string         `json:"vid" mapstructure:"vid"`           // ""
		Weblink    string         `json:"weblink" mapstructure:"weblink"`   // ""
		Dimension  videoDimension `json:"dimension" mapstructure:"dimension"`
		FirstFrame string         `json:"first_frame" mapstructure:"first_frame"`
		CTime      int            `json:"ctime" mapstructure:"ctime"` // 投稿时间
	}
)

// [URL_VIDEO_INFO_WBI], [URL_VIDEO_INFO]
// "data"
type VideoInfo struct {
	Bvid      string `json:"bvid" mapstructure:"bvid"`
	Aid       int    `json:"aid" mapstructure:"aid"`
	Videos    int    `json:"videos" mapstructure:"videos"`       // 分P数
	Copyright int    `json:"copyright" mapstructure:"copyright"` // 1: 原创, 2: 转载
	Pic       string `json:"pic" mapstructure:"pic"`             // 封面 url
	Title     string `json:"title" mapstructure:"title"`         // 标题
	PubDate   int    `json:"pubdate" mapstructure:"pubdate"`     // 发布时间
	CTime     int    `json:"ctime" mapstructure:"ctime"`         // 投稿时间
	Desc      string `json:"desc" mapstructure:"desc"`           // 简介
	State     int    `json:"state" mapstructure:"state"`         // 状态 // TODO: complete
	Duration  int    `json:"duration" mapstructure:"duration"`   // 所有分P总时长(s)

	Owner struct {
		Mid  int    `json:"mid" mapstructure:"mid"`
		Name string `json:"name" mapstructure:"name"`
		Face string `json:"face" mapstructure:"face"`
	} `json:"owner" mapstructure:"owner"`

	Stat struct {
		Aid        int    `json:"aid" mapstructure:"aid"`
		View       int    `json:"view" mapstructure:"view"`         // 播放
		Danmaku    int    `json:"danmaku" mapstructure:"danmaku"`   // 弹幕
		Reply      int    `json:"reply" mapstructure:"reply"`       // 评论
		Favorite   int    `json:"favorite" mapstructure:"favorite"` // 收藏
		Coin       int    `json:"coin" mapstructure:"coin"`         // 投币
		Share      int    `json:"share" mapstructure:"share"`
		NowRank    int    `json:"now_rank" mapstructure:"now_rank"`
		HisRank    int    `json:"his_rank" mapstructure:"his_rank"`
		Like       int    `json:"like" mapstructure:"like"`       // 点赞
		Dislike    int    `json:"dislike" mapstructure:"dislike"` // 0, Deprecated?
		Evaluation string `json:"evaluation" mapstructure:"evaluation"`
		Vt         int    `json:"vt" mapstructure:"vt"`
	} `json:"stat" mapstructure:"stat"`

	ArgueInfo struct { // 争议/警告信息
		ArgueMsg  string `json:"argue_msg" mapstructure:"argue_msg"`
		ArgueType int    `json:"argue_type" mapstructure:"argue_type"`
		ArgueLink string `json:"argue_link" mapstructure:"argue_link"`
	} `json:"argue_info" mapstructure:"argue_info"`

	Dynamic   string         `json:"dynamic" mapstructure:"dynamic"`     // 视频同步发布的的动态的文字内容
	Cid       int            `json:"cid" mapstructure:"cid"`             // P1的 cid
	Dimension videoDimension `json:"dimension" mapstructure:"dimension"` // P1的分辨率

	IsChargeableSeason      bool `json:"is_chargeable_season" mapstructure:"is_chargeable_season"`
	IsStory                 bool `json:"is_story" mapstructure:"is_story"`
	IsUpowerExclusive       bool `json:"is_upower_exclusive" mapstructure:"is_upower_exclusive"` // 充电专属视频
	IsUpowerPay             bool `json:"is_upower_pay" mapstructure:"is_upower_pay"`
	IsUpowerPreview         bool `json:"is_upower_preview" mapstructure:"is_upower_preview"` // 是否支持试看
	IsUpowerExclusiveWithQa bool `json:"is_upower_exclusive_with_qa" mapstructure:"is_upower_exclusive_with_qa"`
	NoCache                 bool `json:"no_cache" mapstructure:"no_cache"` // 是否禁止缓存

	Pages []videoPage `json:"pages" mapstructure:"pages"` // 分P信息

	Online VideoOnline `json:"-" mapstructure:"-"` // 在线人数, 仅用于格式化
}

// WithOnline 为 [VideoInfo] 附加在线人数信息用于格式化
func (vi *VideoInfo) WithOnline(vo VideoOnline) *VideoInfo {
	vi.Online = vo
	return vi
}

func (vi *VideoInfo) DoTemplate() string {
	return DoTemplate(vi)
}

func (vi *VideoInfo) GetReplyList() (rl ReplyList, err error) {
	return FetchReplyList(COMMENT_TYPE_VIDEO, itoa(vi.Aid))
}

// [URL_VIDEO_ONLINE]
// "data"
type VideoOnline struct {
	// 所有平台
	Total string `json:"total" mapstructure:"total"`
	// web端
	Count string `json:"count" mapstructure:"count"`

	ShowSwitch struct {
		Total bool `json:"total" mapstructure:"total"`
		Count bool `json:"count" mapstructure:"count"`
	} `json:"show_switch" mapstructure:"show_switch"`
	AbTest struct {
		Group string `json:"group" mapstructure:"group"`
	} `json:"abtest" mapstructure:"abtest"`
}

func (vo VideoOnline) Ok() bool {
	return vo.Total != "" && vo.Total != "1"
}

type (
	conclusionPartOutline struct {
		Timestamp int    `json:"timestamp" mapstructure:"timestamp"`
		Content   string `json:"content" mapstructure:"content"`
	}
	conclusionOutline struct {
		Title       string                  `json:"title" mapstructure:"title"`
		PartOutline []conclusionPartOutline `json:"part_outline" mapstructure:"part_outline"`
		Timestamp   int                     `json:"timestamp" mapstructure:"timestamp"`
	}
	conclusionPartSubtitle struct {
		StartTimestamp int    `json:"start_timestamp" mapstructure:"start_timestamp"`
		EndTimestamp   int    `json:"end_timestamp" mapstructure:"end_timestamp"`
		Content        string `json:"content" mapstructure:"content"`
	}
	conclusionSubtitle struct {
		Title        string                   `json:"title" mapstructure:"title"`
		PartSubtitle []conclusionPartSubtitle `json:"part_subtitle" mapstructure:"part_subtitle"`
		Timestamp    int                      `json:"timestamp" mapstructure:"timestamp"`
	}
	conclusionResult struct {
		ResultType int                  `json:"result_type" mapstructure:"result_type"`
		Summary    string               `json:"summary" mapstructure:"summary"`   // 大总结
		Outline    []conclusionOutline  `json:"outline" mapstructure:"outline"`   // 大纲
		Subtitle   []conclusionSubtitle `json:"subtitle" mapstructure:"subtitle"` // 字幕
	}
)

// [URL_VIDEO_CONCLUSION_WBI]
// "data"
type VideoConclusion struct {
	Code        int              `json:"code" mapstructure:"code"`
	ModelResult conclusionResult `json:"model_result" mapstructure:"model_result"`
	Stid        string           `json:"stid" mapstructure:"stid"`
	Status      int              `json:"status" mapstructure:"status"`
	LikeNum     int              `json:"like_num" mapstructure:"like_num"`
	DislikeNum  int              `json:"dislike_num" mapstructure:"dislike_num"`
}

func (vc *VideoConclusion) DoTemplate() string {
	return DoTemplate(vc)
}

func (vcf *VideoConclusion) Ok() bool {
	return vcf != nil && vcf.Code == 0 && vcf.ModelResult.Summary != ""
}

type (
	mediaArea struct {
		Id   int    `json:"id" mapstructure:"id"`     // 2
		Name string `json:"name" mapstructure:"name"` // "日本"
	}
	mediaRating struct {
		Count int     `json:"count" mapstructure:"count"`
		Score float64 `json:"score" mapstructure:"score"`
	}
)

// [URL_MEDIA_INFO_BASE]
// "result"
type MediaBase struct {
	Media struct {
		Areas []mediaArea `json:"areas" mapstructure:"areas"`

		Cover             string `json:"cover" mapstructure:"cover"`
		HorizontalPicture string `json:"horizontal_picture" mapstructure:"horizontal_picture"`
		MediaId           int    `json:"media_id" mapstructure:"media_id"` // 1

		NewEp struct {
			Id        int    `json:"id" mapstructure:"id"`
			Index     string `json:"index" mapstructure:"index"`
			IndexShow string `json:"index_show" mapstructure:"index_show"` // "第1话"
		} `json:"new_ep" mapstructure:"new_ep"`

		Rating mediaRating `json:"rating" mapstructure:"rating"`

		SeasonId int    `json:"season_id" mapstructure:"season_id"`
		ShareUrl string `json:"share_url" mapstructure:"share_url"`
		Title    string `json:"title" mapstructure:"title"`
		Type     int    `json:"type" mapstructure:"type"`           // 1
		TypeName string `json:"type_name" mapstructure:"type_name"` // "番剧"
	} `json:"media" mapstructure:"media"`
}

func (mb *MediaBase) DoTemplate() string {
	return DoTemplate(mb)
}

// incomplete
type mediaEpisode struct {
	Aid       int    `json:"aid" mapstructure:"aid"`
	Cid       int    `json:"cid" mapstructure:"cid"`
	Cover     string `json:"cover" mapstructure:"cover"`
	From      string `json:"from" mapstructure:"from"` // "bangumi"
	Id        int    `json:"id" mapstructure:"id"`
	LongTitle string `json:"long_title" mapstructure:"long_title"`
	ShareUrl  string `json:"share_url" mapstructure:"share_url"`
	Title     string `json:"title" mapstructure:"title"`
}

// [URL_MEDIA_INFO_DETAIL]
// "result"
//
// incomplete
type Media struct {
	Actors string `json:"actors" mapstructure:"actors"`
	Alias  string `json:"alias" mapstructure:"alias"`

	Areas []mediaArea `json:"areas" mapstructure:"areas"`

	BkgCover string `json:"bkg_cover" mapstructure:"bkg_cover"`
	Cover    string `json:"cover" mapstructure:"cover"`

	Episodes []mediaEpisode `json:"episodes" mapstructure:"episodes"`

	Evaluate string `json:"evaluate" mapstructure:"evaluate"` // 简介
	JpTitle  string `json:"jp_title" mapstructure:"jp_title"`
	Link     string `json:"link" mapstructure:"link"` // mdid
	MediaId  int    `json:"media_id" mapstructure:"media_id"`
	Mode     int    `json:"mode" mapstructure:"mode"`

	NewEp struct {
		Desc  string `json:"desc" mapstructure:"desc"` // "已完结, 全12话"
		Id    int    `json:"id" mapstructure:"id"`
		IsNew int    `json:"is_new" mapstructure:"is_new"`
		Title string `json:"title" mapstructure:"title"`
	} `json:"new_ep" mapstructure:"new_ep"`

	Publish struct {
		IsFinish      int    `json:"is_finish" mapstructure:"is_finish"`
		IsStarted     int    `json:"is_started" mapstructure:"is_started"`
		PubTime       string `json:"pub_time" mapstructure:"pub_time"`           // "2006-01-02 15:04:05"
		PubTimeShow   string `json:"pub_time_show" mapstructure:"pub_time_show"` // "2006年01月02日15:04"
		UnknowPubDate int    `json:"unknow_pub_date" mapstructure:"unknow_pub_date"`
		Weekday       int    `json:"weekday" mapstructure:"weekday"`
	} `json:"publish" mapstructure:"publish"`

	Rating mediaRating `json:"rating" mapstructure:"rating"`
	Record string      `json:"record" mapstructure:"record"` // 备案号

	SeasonId    int    `json:"season_id" mapstructure:"season_id"`
	SeasonTitle string `json:"season_title" mapstructure:"season_title"`
	// Seasons []struct{} `json:"seasons" mapstructure:"seasons"`
	// Section []struct{} `json:"section" mapstructure:"section"`

	Series struct {
		DisplayType int    `json:"display_type" mapstructure:"display_type"`
		SeriesId    int    `json:"series_id" mapstructure:"series_id"`
		SeriesTitle string `json:"series_title" mapstructure:"series_title"`
	} `json:"series" mapstructure:"series"`

	ShareCopy     string `json:"share_copy" mapstructure:"share_copy"`
	ShareSubTitle string `json:"share_sub_title" mapstructure:"share_sub_title"`
	ShareUrl      string `json:"share_url" mapstructure:"share_url"` // ssid

	SquareCover string `json:"square_cover" mapstructure:"square_cover"`
	Staff       string `json:"staff" mapstructure:"staff"`

	Stat struct {
		Coins      int    `json:"coins" mapstructure:"coins"`
		Danmakus   int    `json:"danmakus" mapstructure:"danmakus"`
		Favorite   int    `json:"favorite" mapstructure:"favorite"`
		Favorites  int    `json:"favorites" mapstructure:"favorites"`
		FollowText string `json:"follow_text" mapstructure:"follow_text"`
		Likes      int    `json:"likes" mapstructure:"likes"`
		Reply      int    `json:"reply" mapstructure:"reply"`
		Share      int    `json:"share" mapstructure:"share"`
		Views      int    `json:"views" mapstructure:"views"`
		Vt         int    `json:"vt" mapstructure:"vt"`
	} `json:"stat" mapstructure:"stat"`

	Status   int      `json:"status" mapstructure:"status"`
	Styles   []string `json:"styles" mapstructure:"styles"`
	SubTitle string   `json:"subtitle" mapstructure:"subtitle"` // "已观看1.7亿次"
	Title    string   `json:"title" mapstructure:"title"`
	Total    int      `json:"total" mapstructure:"total"`
	Type     int      `json:"type" mapstructure:"type"`
}

func (m *Media) DoTemplate() string {
	return DoTemplate(m)
}

type mediaSection struct {
	Episodes []mediaEpisode `json:"episodes" mapstructure:"episodes"`
	Id       int            `json:"id" mapstructure:"id"`
	Title    string         `json:"title" mapstructure:"title"` // "正片" / "预告" / "PV"
	Type     int            `json:"type" mapstructure:"type"`   // 对应 Title
}

// [URL_MEDIA_SECTION]
// "result"
type MediaSection struct {
	MainSection mediaSection   `json:"main_section" mapstructure:"main_section"`
	Section     []mediaSection `json:"section" mapstructure:"section"`
}

// "data"
type LiveStatusUid map[string]*LiveStatus // Uid -> LiveStatus

func (lsu LiveStatusUid) Get(uid any) *LiveStatus {
	switch uid := uid.(type) {
	case string:
		return lsu[uid]
	case int:
		return lsu[itoa(uid)]
	case int64:
		return lsu[itoa(uid)]
	case int32:
		return lsu[itoa(uid)]
	case fmt.Stringer:
		return lsu[uid.String()]
	case nil:
		return nil
	default:
		return lsu[fmt.Sprint(uid)]
	}
}

func (lsu LiveStatusUid) DoTemplate() map[string]string {
	if lsu == nil {
		return nil
	}
	result := make(map[string]string, len(lsu))
	for uid, ls := range lsu {
		result[uid] = DoTemplate(ls)
	}
	return result
}

// [URL_LIVE_STATUS_UIDS]
// "data.{{uid}}"
type LiveStatus struct {
	Title      string `json:"title" mapstructure:"title"`
	RoomId     int    `json:"room_id" mapstructure:"room_id"`
	Uid        int    `json:"uid" mapstructure:"uid"`
	Online     int    `json:"online" mapstructure:"online"`           // 在线人数
	LiveTime   int    `json:"live_time" mapstructure:"live_time"`     // 开播时间戳 (s)
	LiveStatus int    `json:"live_status" mapstructure:"live_status"` // 0 未开播, 1 直播中, 2 轮播中
	ShortId    int    `json:"short_id" mapstructure:"short_id"`       // 短号

	Area             int    `json:"area_id" mapstructure:"area_id"`                         // 分区id
	AreaName         string `json:"area_name" mapstructure:"area_name"`                     // 分区名
	AreaV2Id         int    `json:"area_v2_id" mapstructure:"area_v2_id"`                   // 新版分区id
	AreaV2Name       string `json:"area_v2_name" mapstructure:"area_v2_name"`               // 新版分区名
	AreaV2ParentName string `json:"area_v2_parent_name" mapstructure:"area_v2_parent_name"` // 父分区名
	AreaV2ParentId   int    `json:"area_v2_parent_id" mapstructure:"area_v2_parent_id"`     // 父分区id

	Uname   string `json:"uname" mapstructure:"uname"`       // 用户名
	Face    string `json:"face" mapstructure:"face"`         // 头像url
	TagName string `json:"tag_name" mapstructure:"tag_name"` // 标签, ','分隔
	Tags    string `json:"tags" mapstructure:"tags"`

	CoverFromUser string `json:"cover_from_user" mapstructure:"cover_from_user"` // 直播间封面url
	Keyframe      string `json:"keyframe" mapstructure:"keyframe"`               // 直播间关键帧url, 开播时才有

	LockTill      string `json:"lock_till" mapstructure:"lock_till"`
	HiddenTill    string `json:"hidden_till" mapstructure:"hidden_till"`
	BroadCastType int    `json:"broadcast_type" mapstructure:"broadcast_type"` // 0 普通, 1 手机直播
}

func (ls *LiveStatus) DoTemplate() string {
	return DoTemplate(ls)
}

// [URL_LIVE_INFO_ROOMID]
// "data"
type LiveRoomInfo struct {
	Uid     int `json:"uid" mapstructure:"uid"`           // 主播mid
	RoomId  int `json:"room_id" mapstructure:"room_id"`   // 直播间长号
	ShortId int `json:"short_id" mapstructure:"short_id"` // 直播间短号, 为0是无短号

	Attention   int    `json:"attention" mapstructure:"attention"`     // 关注数量
	Online      int    `json:"online" mapstructure:"online"`           // 观看人数
	IsPortrait  bool   `json:"is_portrait" mapstructure:"is_portrait"` // 是否竖屏
	Description string `json:"description" mapstructure:"description"` // 描述
	LiveStatus  int    `json:"live_status" mapstructure:"live_status"` // 直播状态, 0: 未开播, 1: 直播中, 2: 轮播中

	AreaId         int    `json:"area_id" mapstructure:"area_id"`                   // 分区id
	ParentAreaId   int    `json:"parent_area_id" mapstructure:"parent_area_id"`     // 父分区id
	ParentAreaName string `json:"parent_area_name" mapstructure:"parent_area_name"` // 父分区名称
	OldAreaId      int    `json:"old_area_id" mapstructure:"old_area_id"`           // 旧版分区id

	Background string `json:"background" mapstructure:"background"` // 背景图片链接
	Title      string `json:"title" mapstructure:"title"`           // 标题
	UserCover  string `json:"user_cover" mapstructure:"user_cover"` // 封面
	Keyframe   string `json:"keyframe" mapstructure:"keyframe"`     // 关键帧, 开播后一段时间才有

	IsStrictRoom bool `json:"is_strict_room" mapstructure:"is_strict_room"`

	LiveTime string `json:"live_time" mapstructure:"live_time"` // 直播开始时间, 格式: `2006-01-02 15:04:05`
	Tags     string `json:"tags" mapstructure:"tags"`           // 标签, ','分隔
	IsAnchor int    `json:"is_anchor" mapstructure:"is_anchor"`

	RoomSilentType   string `json:"room_silent_type" mapstructure:"room_silent_type"`     // 禁言状态
	RoomSilentLevel  int    `json:"room_silent_level" mapstructure:"room_silent_level"`   // 禁言等级
	RoomSilentSecond int    `json:"room_silent_second" mapstructure:"room_silent_second"` // 禁言时间, 单位是秒

	AreaName     string `json:"area_name" mapstructure:"area_name"` // 分区名称
	Pendants     string `json:"pendants" mapstructure:"pendants"`
	AreaPendants string `json:"area_pendants" mapstructure:"area_pendants"`

	HotWords       []string `json:"hot_words" mapstructure:"hot_words"`               // 热词
	HotWordsStatus int      `json:"hot_words_status" mapstructure:"hot_words_status"` // 热词状态

	Verify string `json:"verify" mapstructure:"verify"`

	// NewPendants map[string]any `json:"new_pendants" mapstructure:"new_pendants"`

	UpSession            string `json:"up_session" mapstructure:"up_session"`
	PkStatus             int    `json:"pk_status" mapstructure:"pk_status"` // pk状态
	PkId                 int    `json:"pk_id" mapstructure:"pk_id"`         // pk id
	BattleId             int    `json:"battle_id" mapstructure:"battle_id"`
	AllowChangeAreaTime  int    `json:"allow_change_area_time" mapstructure:"allow_change_area_time"`
	AllowUploadCoverTime int    `json:"allow_upload_cover_time" mapstructure:"allow_upload_cover_time"`

	// StudioInfo struct {
	// 	Status     int   `json:"status" mapstructure:"status"`
	// 	MasterList []any `json:"master_list" mapstructure:"master_list"`
	// }
}

// [LiveRoomInfo] 中没有用户名信息, 更建议使用 [LiveStatus.DoTemplate]
func (lri *LiveRoomInfo) DoTemplate() string {
	return DoTemplate(lri)
}

// [URL_ARTICLE_INFO]
// "data"
type ArticleInfo struct {
	Like      int  `json:"like" mapstructure:"like"`           // 0: 未点赞, 1: 已点赞
	Attention bool `json:"attention" mapstructure:"attention"` // 关注
	Favorite  bool `json:"favorite" mapstructure:"favorite"`   // 收藏
	Coin      int  `json:"coin" mapstructure:"coin"`

	Stats struct {
		View     int `json:"view" mapstructure:"view"`         // 阅读数
		Favorite int `json:"favorite" mapstructure:"favorite"` // 收藏数
		Like     int `json:"like" mapstructure:"like"`         // 点赞数
		Dislike  int `json:"dislike" mapstructure:"dislike"`   // 点踩数
		Reply    int `json:"reply" mapstructure:"reply"`       // 评论数
		Share    int `json:"share" mapstructure:"share"`       // 分享数
		Coin     int `json:"coin" mapstructure:"coin"`         // 投币数
		Dynamic  int `json:"dynamic" mapstructure:"dynamic"`   // 动态转发数
	} `json:"stats" mapstructure:"stats"`

	Title      string `json:"title" mapstructure:"title"`             // 标题
	BannerUrl  string `json:"banner_url" mapstructure:"banner_url"`   // 头图url
	Mid        int    `json:"mid" mapstructure:"mid"`                 // 作者uid
	AuthorName string `json:"author_name" mapstructure:"author_name"` // 作者昵称

	IsAuthor bool `json:"is_author" mapstructure:"is_author"`

	ImageUrls       []string `json:"image_urls" mapstructure:"image_urls"`               // 动态封面
	OriginImageUrls []string `json:"origin_image_urls" mapstructure:"origin_image_urls"` // 封面图片

	Shareable       bool `json:"sharable" mapstructure:"sharable"`
	ShowLaterWatch  bool `json:"show_later_watch" mapstructure:"show_later_watch"`
	ShowSmallWindow bool `json:"show_small_window" mapstructure:"show_small_window"`

	InList bool `json:"in_list" mapstructure:"in_list"` // 是否收于文集
	Pre    int  `json:"pre" mapstructure:"pre"`         // 上一篇文章cvid
	Next   int  `json:"next" mapstructure:"next"`       // 下一篇文章cvid

	// ShareChannels []map[string]any `json:"share_channels" mapstructure:"share_channels"`

	Type     int    `json:"type" mapstructure:"type"`           // 0: 文章, 2: 笔记
	VideoUrl string `json:"video_url" mapstructure:"video_url"` // 专栏转视频的活动url
	Location string `json:"location" mapstructure:"location"`   // IP属地

	DisableShare bool `json:"disable_share" mapstructure:"disable_share"`

	Cvid string `json:"-" mapstructure:"-"` // 专栏id, 仅用于格式化
}

// WithCvid 为 [ArticleInfo] 附加 cvid 用于格式化
func (ai *ArticleInfo) WithCvid(cvid string) *ArticleInfo {
	ai.Cvid = strings.TrimPrefix(cvid, "cv")
	return ai
}

func (ai *ArticleInfo) DoTemplate() string {
	return DoTemplate(ai)
}

// [URL_SONG_INFO]
// "data"
type SongInfo struct {
	Id     int    `json:"id" mapstructure:"id"`
	Uid    int    `json:"uid" mapstructure:"uid"`
	Uname  string `json:"uname" mapstructure:"uname"`
	Author string `json:"author" mapstructure:"author"` // 作者名
	Title  string `json:"title" mapstructure:"title"`
	Cover  string `json:"cover" mapstructure:"cover"` // 封面图片url
	Intro  string `json:"intro" mapstructure:"intro"` // 简介
	Lyric  string `json:"lyric" mapstructure:"lyric"` // lrc歌词url

	Crtype int `json:"crtype" mapstructure:"crtype"` // 1

	Duration int `json:"duration" mapstructure:"duration"` // 时长 (s)
	Passtime int `json:"passtime" mapstructure:"passtime"` // 发布时间 (s)
	Curtime  int `json:"curtime" mapstructure:"curtime"`   // 当前请求时间 (s)

	Aid  int    `json:"aid" mapstructure:"aid"`   // 关联稿件
	Bvid string `json:"bvid" mapstructure:"bvid"` // 关联稿件
	Cid  int    `json:"cid" mapstructure:"cid"`   // 关联稿件

	Msid       int    `json:"msid" mapstructure:"msid"`
	Attr       int    `json:"attr" mapstructure:"attr"`
	Limit      int    `json:"limit" mapstructure:"limit"`
	ActivityId int    `json:"activity_id" mapstructure:"activity_id"`
	Limitdesc  string `json:"limitdesc" mapstructure:"limitdesc"`

	CoinNum int `json:"coin_num" mapstructure:"coin_num"` // 投币
	CTime   int `json:"ctime" mapstructure:"ctime"`       // 投稿时间 (ms)

	Statistic struct {
		Sid     int `json:"sid" mapstructure:"sid"`         // auid
		Play    int `json:"play" mapstructure:"play"`       // 播放
		Collect int `json:"collect" mapstructure:"collect"` // 收藏
		Comment int `json:"comment" mapstructure:"comment"` // 评论
		Share   int `json:"share" mapstructure:"share"`     // 分享
	} `json:"statistic" mapstructure:"statistic"`

	// VipInfo struct {
	// 	Type       int `json:"type" mapstructure:"type"` // 1: 月, 2: 年
	// 	Status     int `json:"status" mapstructure:"status"`
	// 	DueDate    int `json:"due_date" mapstructure:"due_date"` // (ms)
	// 	VipPayType int `json:"vip_pay_type" mapstructure:"vip_pay_type"`
	// } `json:"vipInfo" mapstructure:"vipInfo"` // UP主会员状态

	CollectIds []int `json:"collectIds" mapstructure:"collectIds"` // 歌曲所在的收藏夹mlid
	IsCooper   bool  `json:"isCooper" mapstructure:"isCooper"`     // ?
}

// [URL_SONG_TAG]
// "data"
//
// 只有 Info 字段是有用的
type SongTag []struct {
	Type    string `json:"type" mapstructure:"type"` // "song"
	Subtype int    `json:"subtype" mapstructure:"subtype"`
	Key     int    `json:"key" mapstructure:"key"`
	Info    string `json:"info" mapstructure:"info"` // "音乐", "人声", "翻唱", "粤语", "偶像", "流行"...
}

func (st SongTag) String() string {
	return st.Join(", ")
}

// Join 格式化 Info 字段, 用 sep 分隔
func (st SongTag) Join(sep string) string {
	if len(st) == 0 {
		return ""
	}
	tags := make([]string, len(st))
	for i := range st {
		tags[i] = st[i].Info
	}
	return strings.Join(tags, sep)
}

func (st SongTag) Ok() bool {
	return len(st) > 0
}

// [URL_SONG_MEMBER]
// "data"
//
// 只有 Type 与 List.Name 字段是有用的
type SongMember []struct {
	Type int `json:"type" mapstructure:"type"` // [SongMemberMap]
	List []struct {
		Mid      int    `json:"mid" mapstructure:"mid"`             // 0
		MemberId int    `json:"member_id" mapstructure:"member_id"` // ?
		Name     string `json:"name" mapstructure:"name"`           // 昵称
	} `json:"list" mapstructure:"list"`
}

func (sm SongMember) String() string {
	return sm.Format("%s：%s", "\n", "/")
}

// Format 格式化成员信息
//
//	"%s：%s", "\n", "/"
//
//	歌手：师欣/琉盈君/花音/小小六
//	作词：小小六
//	...
func (sm SongMember) Format(format string, sep string, subSep string) string {
	if len(sm) == 0 {
		return ""
	}

	if format == "" {
		format = "%s：%s"
	}

	members := make([]string, len(sm))
	for i := range sm {
		names := make([]string, len(sm[i].List))
		for j := range sm[i].List {
			names[j] = sm[i].List[j].Name
		}

		members[i] = fmt.Sprintf(
			format,
			SongMemberTable(sm[i].Type),
			strings.Join(names, subSep),
		)
	}
	return strings.Join(members, sep)
}

func (sm SongMember) Ok() bool {
	return len(sm) > 0
}

// SongMemberTable 返回歌曲成员类型的名称,
// 未知类型返回数字字符串
func SongMemberTable[T integer](i T) (s string) {
	if i < 1 || int(i) > len(songMemberTable) || songMemberTable[i] == "" {
		return itoa(i)
	}
	return songMemberTable[i]
}

var songMemberTable = [...]string{
	1:  "歌手",
	2:  "作词",
	3:  "作曲",
	4:  "编曲",
	5:  "后期/混音",
	6:  "封面制作",
	7:  "封面制作",
	8:  "音源",
	9:  "调音",
	10: "演奏",
	11: "乐器",

	127: "UP主",
}

// [URL_SONG_LIRIC]
// "data"
//
// lrc 格式
type SongLyric = string

type Song struct {
	Info   SongInfo
	Tag    SongTag
	Member SongMember
	Lyric  SongLyric
}

func (s *Song) DoTemplate() string {
	return DoTemplate(s)
}

// [URL_SPACE_CARD]
// "data"
type SpaceCard struct {
	Card struct {
		Mid  string `json:"mid" mapstructure:"mid"`   // uid
		Name string `json:"name" mapstructure:"name"` // 昵称
		Sex  string `json:"sex" mapstructure:"sex"`   // 性别 // "男"
		Face string `json:"face" mapstructure:"face"` // 头像 url

		Fans      int `json:"fans" mapstructure:"fans"`           // 粉丝数
		Friend    int `json:"friend" mapstructure:"friend"`       // 关注数
		Attention int `json:"attention" mapstructure:"attention"` // 关注数

		Sign string `json:"sign" mapstructure:"sign"` // 签名

		LevelInfo struct {
			CurrentLevel int `json:"current_level" mapstructure:"current_level"`
			CurrentMin   int `json:"current_min" mapstructure:"current_min"`
			CurrentExp   int `json:"current_exp" mapstructure:"current_exp"`
			NextExp      int `json:"next_exp" mapstructure:"next_exp"`
		} `json:"level_info" mapstructure:"level_info"`

		Pendant struct { // 头像框
			Pid               int    `json:"pid" mapstructure:"pid"`
			Name              string `json:"name" mapstructure:"name"`
			Image             string `json:"image" mapstructure:"image"`
			Expire            int    `json:"expire" mapstructure:"expire"`
			ImageEnhance      string `json:"image_enhance" mapstructure:"image_enhance"`
			ImageEnhanceFrame string `json:"image_enhance_frame" mapstructure:"image_enhance_frame"`
			NPid              int    `json:"n_pid" mapstructure:"n_pid"`
		} `json:"pendant" mapstructure:"pendant"`
	} `json:"card" mapstructure:"card"`

	Space struct { // 主页头图
		SImg string `json:"s_img" mapstructure:"s_img"`
		LImg string `json:"l_img" mapstructure:"l_img"`
	}

	ArchiveCount int `json:"archive_count" mapstructure:"archive_count"` // 投稿数
	ArticleCount int `json:"article_count" mapstructure:"article_count"` // 专栏数 // 0, Deprecated?
	Follower     int `json:"follower" mapstructure:"follower"`           // 粉丝数
	LikeNum      int `json:"like_num" mapstructure:"like_num"`           // 获赞数
}

func (sc *SpaceCard) DoTemplate() string {
	return DoTemplate(sc)
}

// "data"
type RelationStat struct {
	Mid       int `json:"mid" mapstructure:"mid"`
	Following int `json:"following" mapstructure:"following"` // 关注数
	Whisper   int `json:"whisper" mapstructure:"whisper"`
	Black     int `json:"black" mapstructure:"black"`
	Follower  int `json:"follower" mapstructure:"follower"` // 粉丝数
}

// "data"
type DynamicAll struct {
	HasMore        bool            `json:"has_more" mapstructure:"has_more"`
	Items          []DynamicDetail `json:"items" mapstructure:"items"`
	Offset         string          `json:"offset" mapstructure:"offset"`
	UpdateBaseline string          `json:"update_baseline" mapstructure:"update_baseline"`
	UpdateNum      int             `json:"update_num" mapstructure:"update_num"`
}

// "data"
type DynamicAllUpdate struct {
	UpdateNum int `json:"update_num" mapstructure:"update_num"`
}

// "basic"
type dynamicBasic struct {
	CommentIdStr string `json:"comment_id_str" mapstructure:"comment_id_str"`
	CommentType  int    `json:"comment_type" mapstructure:"comment_type"`
	IsOnlyFans   bool   `json:"is_only_fans" mapstructure:"is_only_fans"`
	// LikeIcon struct{} `json:"like_icon" mapstructure:"like_icon"`
	RidStr string `json:"rid_str" mapstructure:"rid_str"` // == CommentIdStr
}

type (
	// "emoji"
	dynamicRichTextNodeEmoji struct {
		IconUrl string `json:"icon_url" mapstructure:"icon_url"`
		Size    int    `json:"size" mapstructure:"size"`
		Text    string `json:"text" mapstructure:"text"`
		Type    int    `json:"type" mapstructure:"type"`
	}
	// "rich_text_nodes"
	dynamicRichTextNode struct {
		Emoji *dynamicRichTextNodeEmoji `json:"emoji" mapstructure:"emoji"` // "RICH_TEXT_NODE_TYPE_EMOJI"

		OrigText string `json:"orig_text" mapstructure:"orig_text"`
		Text     string `json:"text" mapstructure:"text"`
		Type     string `json:"type" mapstructure:"type"` // "RICH_TEXT_NODE_TYPE_xxx"
		// "RICH_TEXT_NODE_TYPE_VOTE" 投票 id
		// "RICH_TEXT_NODE_TYPE_AT" 用户 uid
		Rid string `json:"rid" mapstructure:"rid"`

		// Goods *struct{} `json:"goods" mapstructure:"goods"` // "RICH_TEXT_NODE_TYPE_GOODS"

		// "RICH_TEXT_NODE_TYPE_TOPIC" 跳转搜索
		// "RICH_TEXT_NODE_TYPE_GOODS" 商品链接
		JumpUrl string `json:"jump_url" mapstructure:"jump_url"`
	}
)

type (
	// "reserve"
	dynamicAddiReserve struct {
		// Button struct{} `json:"button" mapstructure:"button"`
		Desc1 struct {
			Style int    `json:"style" mapstructure:"style"`
			Text  string `json:"text" mapstructure:"text"` // "预计..." / "...直播"
		} `json:"desc1" mapstructure:"desc1"`
		Desc2 struct {
			Style   int    `json:"style" mapstructure:"style"`
			Text    string `json:"text" mapstructure:"text"` // "...万观看" / "...万人预约"
			Visible bool   `json:"visible" mapstructure:"visible"`
		} `json:"desc2" mapstructure:"desc2"`
		Desc3 struct {
			JumpUrl string `json:"jump_url" mapstructure:"jump_url"`
			Style   int    `json:"style" mapstructure:"style"`
			Text    string `json:"text" mapstructure:"text"` // "预约有奖：..."
		} `json:"desc3" mapstructure:"desc3"`

		JumpUrl      string `json:"jump_url" mapstructure:"jump_url"`           // "https:" + ...
		ReserveTotal int    `json:"reserve_total" mapstructure:"reserve_total"` // 预约人数
		Rid          int    `json:"rid" mapstructure:"rid"`
		State        int    `json:"state" mapstructure:"state"`
		Stype        int    `json:"stype" mapstructure:"stype"`
		Title        string `json:"title" mapstructure:"title"`
		UpMid        int    `json:"up_mid" mapstructure:"up_mid"`
	}

	// "vote"
	dynamicAddiVote struct {
		ChoiceCnt    int    `json:"choice_cnt" mapstructure:"choice_cnt"`
		DefaultShare int    `json:"default_share" mapstructure:"default_share"`
		Desc         string `json:"desc" mapstructure:"desc"`
		EndTime      int    `json:"end_time" mapstructure:"end_time"` // (s)
		JoinNum      int    `json:"join_num" mapstructure:"join_num"` // 参与人数
		Status       int    `json:"status" mapstructure:"status"`
		Type         int    `json:"type" mapstructure:"type"`
		Uid          int    `json:"uid" mapstructure:"uid"`
		VoteId       int    `json:"vote_id" mapstructure:"vote_id"`
	}

	// "ugc" 视频跳转
	dynamicAddiUgc struct {
		Cover      string `json:"cover" mapstructure:"cover"`
		DescSecond string `json:"desc_second" mapstructure:"desc_second"` // "xx观看 xx弹幕"
		Duration   string `json:"duration" mapstructure:"duration"`       // "04:05"
		HeadText   string `json:"head_text" mapstructure:"head_text"`     // ""
		IdStr      string `json:"id_str" mapstructure:"id_str"`           // avid
		JumpUrl    string `json:"jump_url" mapstructure:"jump_url"`       // "https:" + ...
		MultiLine  bool   `json:"multi_line" mapstructure:"multi_line"`
		Title      string `json:"title" mapstructure:"title"`
	}

	// "common" 相关游戏, ...
	dynamicAddiCommon struct {
		// Button struct{} `json:"button" mapstructure:"button"`
		Cover    string `json:"cover" mapstructure:"cover"`
		Desc1    string `json:"desc1" mapstructure:"desc1"`         // "策略/架空文明/末世"
		Desc2    string `json:"desc2" mapstructure:"desc2"`         // "六周年庆典现已开启"
		HeadText string `json:"head_text" mapstructure:"head_text"` // "相关游戏"
		IdStr    string `json:"id_str" mapstructure:"id_str"`       // gameId
		JumpUrl  string `json:"jump_url" mapstructure:"jump_url"`   // biligame.com/detail?id=101772
		Style    int    `json:"style" mapstructure:"style"`         // 1
		SubType  string `json:"sub_type" mapstructure:"sub_type"`   // "game"
		Title    string `json:"title" mapstructure:"title"`         // "明日方舟"
	}
)

func (da *dynamicAddiReserve) DoTemplate() string {
	return DoTemplate(da)
}

// 获取投票的详细信息 [VoteInfo] 再填充模板
func (dav *dynamicAddiVote) DoTemplate() string {
	if dav.VoteId == 0 {
		return "投票ID为0"
	}
	vi, err := FetchVoteInfo(itoa(dav.VoteId))
	if err != nil {
		return err.Error()
	}
	return vi.DoTemplate()
}

func (da *dynamicAddiUgc) DoTemplate() string {
	return DoTemplate(da)
}

func (da *dynamicAddiCommon) DoTemplate() string {
	return DoTemplate(da)
}

type (
	// "none"
	dynamicMajorNone struct {
		// "源动态已被作者删除"
		Tips string `json:"tips" mapstructure:"tips"`
	}

	// "blocked"
	dynamicMajorBlocked struct {
		// BgImg struct{} `json:"bg_img" mapstructure:"bg_img"`
		BlockedType int `json:"blocked_type" mapstructure:"blocked_type"` // 1: 充电动态
		// Button struct{} `json:"button" mapstructure:"button"`
		HintMessage string `json:"hint_message" mapstructure:"hint_message"`
		// Icon struct{} `json:"icon" mapstructure:"icon"`
	}

	dynamicMajorDrawItem struct {
		Height int     `json:"height" mapstructure:"height"`
		Size   float64 `json:"size" mapstructure:"size"`
		Src    string  `json:"src" mapstructure:"src"`
		Tags   []any   `json:"tags" mapstructure:"tags"` // TODO: determine type
		Width  int     `json:"width" mapstructure:"width"`
	}
	// "draw"
	dynamicMajorDraw struct {
		Id    int                    `json:"id" mapstructure:"id"`
		Items []dynamicMajorDrawItem `json:"items" mapstructure:"items"`
	}

	// "archive"
	dynamicMajorArchive struct {
		Aid string `json:"aid" mapstructure:"aid"`
		// Badge struct{} `json:"badge" mapstructure:"badge"`
		Bvid           string `json:"bvid" mapstructure:"bvid"`
		Cover          string `json:"cover" mapstructure:"cover"`
		Desc           string `json:"desc" mapstructure:"desc"`
		DisablePreview int    `json:"disable_preview" mapstructure:"disable_preview"`
		DurationText   string `json:"duration_text" mapstructure:"duration_text"`
		JumpUrl        string `json:"jump_url" mapstructure:"jump_url"` // `https:` + ...
		Stat           struct {
			Danmaku string `json:"danmaku" mapstructure:"danmaku"`
			Play    string `json:"play" mapstructure:"play"`
		} `json:"stat" mapstructure:"stat"`
		Title string `json:"title" mapstructure:"title"`
		Type  int    `json:"type" mapstructure:"type"`
	}

	// "article"
	dynamicMajorArticle struct {
		Covers  []string `json:"covers" mapstructure:"covers"`
		Desc    string   `json:"desc" mapstructure:"desc"`
		Id      int      `json:"id" mapstructure:"id"`
		JumpUrl string   `json:"jump_url" mapstructure:"jump_url"` // "https:" + ...
		Label   string   `json:"label" mapstructure:"label"`
		Title   string   `json:"title" mapstructure:"title"`
	}

	// "music"
	dynamicMajorMusic struct {
		Cover   string `json:"cover" mapstructure:"cover"`
		Id      int    `json:"id" mapstructure:"id"`             // auid
		JumpUrl string `json:"jump_url" mapstructure:"jump_url"` // "https:" + ...
		Label   string `json:"label" mapstructure:"label"`       // "音乐 · 人声演唱"
		Title   string `json:"title" mapstructure:"title"`
	}

	// "live"
	dynamicMajorLive struct {
		// Badge struct{} `json:"badge" mapstructure:"badge"`
		Cover       string `json:"cover" mapstructure:"cover"`
		DescFirst   string `json:"desc_first" mapstructure:"desc_first"`   // "虚拟日常"
		DescSecond  string `json:"desc_second" mapstructure:"desc_second"` // "2483人看过"
		Id          int    `json:"id" mapstructure:"id"`                   // 直播间 id
		JumpUrl     string `json:"jump_url" mapstructure:"jump_url"`       // "https:" + ...
		LiveState   int    `json:"live_state" mapstructure:"live_state"`
		ReserveType int    `json:"reserve_type" mapstructure:"reserve_type"`
		Title       string `json:"title" mapstructure:"title"`
	}

	liveRcmd struct {
		Type         int `json:"type" mapstructure:"type"`
		LivePlayInfo struct {
			Online         int    `json:"online" mapstructure:"online"`
			AreaId         int    `json:"area_id" mapstructure:"area_id"`
			LiveId         string `json:"live_id" mapstructure:"live_id"`
			RoomPaidType   int    `json:"room_paid_type" mapstructure:"room_paid_type"`
			Cover          string `json:"cover" mapstructure:"cover"`
			ParentAreaName string `json:"parent_area_name" mapstructure:"parent_area_name"` // "虚拟主播"
			LiveStartTime  int    `json:"live_start_time" mapstructure:"live_start_time"`   // (s)
			Link           string `json:"link" mapstructure:"link"`                         // "https:" + ...
			ParentAreaId   int    `json:"parent_area_id" mapstructure:"parent_area_id"`

			WatchedShow struct {
				Num       int    `json:"num" mapstructure:"num"`
				TextSmall string `json:"text_small" mapstructure:"text_small"` // == itoa(Num)
				TextLarge string `json:"text_large" mapstructure:"text_large"` // "itoa(Num)人看过"
				// Icon string `json:"icon" mapstructure:"icon"`
				// IconLocation string `json:"icon_location" mapstructure:"icon_location"`
				// IconWeb string `json:"icon_web" mapstructure:"icon_web"`
				Switch bool `json:"switch" mapstructure:"switch"`
			} `json:"watched_show" mapstructure:"watched_show"`

			LiveStatus int    `json:"live_status" mapstructure:"live_status"`
			RoomType   int    `json:"room_type" mapstructure:"room_type"`
			PlayType   int    `json:"play_type" mapstructure:"play_type"`
			AreaName   string `json:"area_name" mapstructure:"area_name"` // "虚拟日常"
			// Pendants struct{} `json:"pendants" mapstructure:"pendants"`
			RoomId         int    `json:"room_id" mapstructure:"room_id"`
			Uid            int    `json:"uid" mapstructure:"uid"`
			Title          string `json:"title" mapstructure:"title"`
			LiveScreenType int    `json:"live_screen_type" mapstructure:"live_screen_type"`
		} `json:"live_play_info" mapstructure:"live_play_info"`
		// LiveRecordInfo *struct{} `json:"live_record_info" mapstructure:"live_record_info"`
	}
	// "live_rcmd"
	dynamicMajorLiveRcmd struct {
		Content      string   `json:"content" mapstructure:"content"` // json
		ReserveType  int      `json:"reserve_type" mapstructure:"reserve_type"`
		Unmarshaled  liveRcmd `json:"-" mapstructure:"-"`
		UnmarshalErr error    `json:"-" mapstructure:"-"`
	}

	// "pgc"
	dynamicMajorPgc struct {
		// Badge struct{} `json:"badge" mapstructure:"badge"`
		Cover    string `json:"cover" mapstructure:"cover"`
		Epid     int    `json:"epid" mapstructure:"epid"`
		JumpUrl  string `json:"jump_url" mapstructure:"jump_url"`
		SeasonId int    `json:"season_id" mapstructure:"season_id"`
		Stat     struct {
			Danmaku string `json:"danmaku" mapstructure:"danmaku"`
			Play    string `json:"play" mapstructure:"play"`
		} `json:"stat" mapstructure:"stat"`
		SubType int    `json:"sub_type" mapstructure:"sub_type"`
		Title   string `json:"title" mapstructure:"title"`
		Type    int    `json:"type" mapstructure:"type"`
	}
)

// lazyload 导致必须要显式调用 [DoTemplate],
// 不能直接在模板中写{{template "dynamicMajorDraw" .Draw}},

func (dmd *dynamicMajorDraw) DoTemplate() string {
	return DoTemplate(dmd)
}

func (dma *dynamicMajorArchive) DoTemplate() string {
	return DoTemplate(dma)
}

func (dma *dynamicMajorArticle) DoTemplate() string {
	return DoTemplate(dma)
}

func (dmm *dynamicMajorMusic) DoTemplate() string {
	return DoTemplate(dmm)
}

func (dml *dynamicMajorLive) DoTemplate() string {
	return DoTemplate(dml)
}

func (dmlr *dynamicMajorLiveRcmd) DoTemplate() string {
	return DoTemplate(dmlr)
}

func (dmp *dynamicMajorPgc) DoTemplate() string {
	return DoTemplate(dmp)
}

// Unmarshal 将 Content 中的 rawJson 反序列化到 Unmarshaled 中
func (dmlr *dynamicMajorLiveRcmd) Unmarshal() (ok bool) {
	dmlr.UnmarshalErr = json.Unmarshal([]byte(dmlr.Content), &dmlr.Unmarshaled)
	return dmlr.UnmarshalErr == nil
}

type (
	// "additional"
	dynamicAdditional struct {
		Type string `json:"type" mapstructure:"type"` // "ADDITIONAL_TYPE_xxx"

		Reserve *dynamicAddiReserve `json:"reserve" mapstructure:"reserve"` // "ADDITIONAL_TYPE_RESERVE"
		Vote    *dynamicAddiVote    `json:"vote" mapstructure:"vote"`       // "ADDITIONAL_TYPE_VOTE"
		Ugc     *dynamicAddiUgc     `json:"ugc" mapstructure:"ugc"`         // "ADDITIONAL_TYPE_UGC"
		Common  *dynamicAddiCommon  `json:"common" mapstructure:"common"`   // "ADDITIONAL_TYPE_COMMON"
	}

	// "desc"
	dynamicDesc struct {
		Text string `json:"text" mapstructure:"text"`

		RichTextNodes []dynamicRichTextNode `json:"rich_text_nodes" mapstructure:"rich_text_nodes"`
	}

	// "major"
	dynamicMajor struct {
		Type string `json:"type" mapstructure:"type"` // "MAJOR_TYPE_xxx"

		None     *dynamicMajorNone     `json:"none" mapstructure:"none"`           // "MAJOR_TYPE_NONE"
		Blocked  *dynamicMajorBlocked  `json:"blocked" mapstructure:"blocked"`     // "MAJOR_TYPE_BLOCKED"
		Draw     *dynamicMajorDraw     `json:"draw" mapstructure:"draw"`           // "MAJOR_TYPE_DRAW"
		Archive  *dynamicMajorArchive  `json:"archive" mapstructure:"archive"`     // "MAJOR_TYPE_ARCHIVE"
		Article  *dynamicMajorArticle  `json:"article" mapstructure:"article"`     // "MAJOR_TYPE_ARTICLE"
		Music    *dynamicMajorMusic    `json:"music" mapstructure:"music"`         // "MAJOR_TYPE_MUSIC"
		Live     *dynamicMajorLive     `json:"live" mapstructure:"live"`           // "MAJOR_TYPE_LIVE"
		LiveRcmd *dynamicMajorLiveRcmd `json:"live_rcmd" mapstructure:"live_rcmd"` // "MAJOR_TYPE_LIVE_RCMD"
		Pgc      *dynamicMajorPgc      `json:"pgc" mapstructure:"pgc"`             // "MAJOR_TYPE_PGC"
	}

	// "topic"
	dynamicTopic struct {
		Id      int    `json:"id" mapstructure:"id"`
		JumpUrl string `json:"jump_url" mapstructure:"jump_url"`
		Name    string `json:"name" mapstructure:"name"`
	}
)

type (
	// "module_author"
	dynamicModuleAuthor struct {
		Face    string `json:"face" mapstructure:"face"`
		FaceNft bool   `json:"face_nft" mapstructure:"face_nft"`
		Mid     int    `json:"mid" mapstructure:"mid"`
		Name    string `json:"name" mapstructure:"name"`

		IconBadge struct {
			Text string `json:"text" mapstructure:"text"` // "充电专属"
		} `json:"icon_badge" mapstructure:"icon_badge"`

		JumpUrl string `json:"jump_url" mapstructure:"jump_url"` // "https:" + ...

		// "AUTHOR_TYPE_PGC" "番剧"
		Label string `json:"label" mapstructure:"label"`

		PubAction       string `json:"pub_action" mapstructure:"pub_action"` // "投稿了视频"/"发布了动态视频"/"投稿了文章"/"直播了"...
		PubLocationText string `json:"pub_location_text" mapstructure:"pub_location_text"`
		PubTime         string `json:"pub_time" mapstructure:"pub_time"` // 发布时间 "2006年01月02日 15:04"
		PubTs           int    `json:"pub_ts" mapstructure:"pub_ts"`     // (s)

		// "AUTHOR_TYPE_NORMAL" 用户
		// "AUTHOR_TYPE_PGC"    番剧
		Type string `json:"type" mapstructure:"type"`
	}

	// "module_dynamic"
	dynamicModuleDynamic struct {
		Additional *dynamicAdditional `json:"additional" mapstructure:"additional"`
		Desc       *dynamicDesc       `json:"desc" mapstructure:"desc"`
		Major      *dynamicMajor      `json:"major" mapstructure:"major"`
		Topic      *dynamicTopic      `json:"topic" mapstructure:"topic"`
	}

	// "module_more" 忽略
	// dynamicModuleMore struct{} // three_point_items, THREE_POINT_REPORT

	// "module_stat"
	dynamicModuleStat struct {
		Comment struct {
			Count     int  `json:"count" mapstructure:"count"`
			Forbidden bool `json:"forbidden" mapstructure:"forbidden"`
		} `json:"comment" mapstructure:"comment"`
		Forward struct {
			Count     int  `json:"count" mapstructure:"count"`
			Forbidden bool `json:"forbidden" mapstructure:"forbidden"`
		} `json:"forward" mapstructure:"forward"`
		Like struct {
			Count     int  `json:"count" mapstructure:"count"`
			Forbidden bool `json:"forbidden" mapstructure:"forbidden"`
			Status    bool `json:"status" mapstructure:"status"`
		} `json:"like" mapstructure:"like"`
	}
)

// "data.item" / "data.item.orig"
type DynamicDetail struct {
	Basic dynamicBasic `json:"basic" mapstructure:"basic"`
	IdStr string       `json:"id_str" mapstructure:"id_str"` // 动态 id

	Modules struct {
		Author  dynamicModuleAuthor  `json:"module_author" mapstructure:"module_author"`
		Dynamic dynamicModuleDynamic `json:"module_dynamic" mapstructure:"module_dynamic"`
		// More dynamicModuleMore `json:"module_more" mapstructure:"module_more"`
		Stat dynamicModuleStat `json:"module_stat" mapstructure:"module_stat"`
	} `json:"modules" mapstructure:"modules"`

	Orig *DynamicDetail `json:"orig" mapstructure:"orig"` // 转发的动态

	// "DYNAMIC_TYPE_FORWARD"       //
	// "DYNAMIC_TYPE_NONE"          // 只存在于转发的动态 ".orig" 中
	// "DYNAMIC_TYPE_WORD"          // 没有对应的 "MAJOR_TYPE_xxx"
	// "DYNAMIC_TYPE_DRAW"          -> "MAJOR_TYPE_DRAW"
	// "DYNAMIC_TYPE_AV"            -> "MAJOR_TYPE_ARCHIVE"
	// "DYNAMIC_TYPE_ARTICLE"       -> "MAJOR_TYPE_ARTICLE"
	// "DYNAMIC_TYPE_MUSIC"         -> "MAJOR_TYPE_MUSIC"
	// "DYNAMIC_TYPE_LIVE"          // 直播间分享
	// "DYNAMIC_TYPE_LIVE_RCMD"     -> "MAJOR_TYPE_LIVE_RCMD" // 直播开播(动态流拿不到更新)
	// "DYNAMIC_TYPE_COMMON_SQUARE" // 装扮 / 剧集点评 / 普通分享
	// "DYNAMIC_TYPE_PGC_UNION"     -> "MAJOR_TYPE_PGC" // 番剧更新
	Type string `json:"type" mapstructure:"type"` // "DYNAMIC_TYPE_xxx"
	// block 的动态也能看到 "DYNAMIC_TYPE_DRAW"

	Visible bool `json:"visible" mapstructure:"visible"`
}

func (dd *DynamicDetail) DoTemplate() string {
	return DoTemplate(dd)
}

func (dd *DynamicDetail) GetReplyList() (rl ReplyList, err error) {
	var typ CommentClass
	var id string

	if dd.Modules.Dynamic.Major != nil &&
		dd.Modules.Dynamic.Major.Draw != nil {
		typ = COMMENT_TYPE_ALBUM
		id = itoa(dd.Modules.Dynamic.Major.Draw.Id)
	} else {
		typ = COMMENT_TYPE_DYNAMIC
		id = dd.IdStr
	}

	return FetchReplyList(typ, id)
}

type voteOption struct {
	OptIdx  int    `json:"opt_idx" mapstructure:"opt_idx"`   // 选项索引
	OptDesc string `json:"opt_desc" mapstructure:"opt_desc"` // 选项描述
	ImgUrl  string `json:"img_url" mapstructure:"img_url"`   // 选项图片url
	Cnt     int    `json:"cnt" mapstructure:"cnt"`           // 选项投票数
}

// "data.vote_info"
type VoteInfo struct {
	VoteId        int          `json:"vote_id" mapstructure:"vote_id"`
	Title         string       `json:"title" mapstructure:"title"`
	Desc          string       `json:"desc" mapstructure:"desc"`
	JoinNum       int          `json:"join_num" mapstructure:"join_num"`
	Type          int          `json:"type" mapstructure:"type"`
	ChoiceCnt     int          `json:"choice_cnt" mapstructure:"choice_cnt"` // 可选数
	EndTime       int          `json:"end_time" mapstructure:"end_time"`     // 结束时间 (s)
	Status        int          `json:"status" mapstructure:"status"`
	VotePublisher int          `json:"vote_publisher" mapstructure:"vote_publisher"` // 投票发起者uid
	DefaultShare  int          `json:"default_share" mapstructure:"default_share"`
	CTime         int          `json:"ctime" mapstructure:"ctime"` // 创建时间 (s)
	Options       []voteOption `json:"options" mapstructure:"options"`
	OptionsCnt    int          `json:"options_cnt" mapstructure:"options_cnt"` // 选项数
	VoteLevel     int          `json:"vote_level" mapstructure:"vote_level"`   // 投票发起者账号等级
	Face          string       `json:"face" mapstructure:"face"`               // 投票发起者头像url
	Name          string       `json:"name" mapstructure:"name"`               // 投票发起者昵称
}

func (vi *VoteInfo) DoTemplate() string {
	return DoTemplate(vi)
}

/*
unders are without any template
*/

type searchPageInfo struct {
	Total int `json:"total" mapstructure:"total"`
	// NumResults int `json:"numResults" mapstructure:"numResults"` // == total
	Pages int `json:"pages" mapstructure:"pages"`
	// NumPages int `json:"numPages" mapstructure:"numPages"` // == pages
}

type searchBase struct {
	Seid           string `json:"seid" mapstructure:"seid"`
	Page           int    `json:"page" mapstructure:"page"`
	Pagesize       int    `json:"pagesize" mapstructure:"pagesize"`
	Next           int    `json:"next" mapstructure:"next"`
	NumResults     int    `json:"numResults" mapstructure:"numResults"`
	NumPages       int    `json:"numPages" mapstructure:"numPages"`
	SuggestKeyword string `json:"suggest_keyword" mapstructure:"suggest_keyword"`
	RqtType        string `json:"rqt_type" mapstructure:"rqt_type"`

	// CostTime struct{} `json:"cost_time" mapstructure:"cost_time"`
	// ExpList struct{...bool} `json:"exp_list" mapstructure:"exp_list"`

	EggHit int `json:"egg_hit" mapstructure:"egg_hit"`
}

// "data"
type SearchAll struct {
	searchBase

	// "video", "bangumi", ...
	PageInfo map[string]searchPageInfo `json:"page_info" mapstructure:"page_info"`
	Result   []struct {
		ResultType string           `json:"result_type" mapstructure:"result_type"`
		Data       []map[string]any `json:"data" mapstructure:"data"`
	} `json:"result" mapstructure:"result"`
}

// "data"
type SearchType struct {
	searchBase
	Result []map[string]any `json:"result" mapstructure:"result"`
}

// "data"
type SearchTypeLive struct {
	searchBase
	Pageinfo struct {
		LiveUser searchPageInfo `json:"live_user" mapstructure:"live_user"`
		LiveRoom searchPageInfo `json:"live_room" mapstructure:"live_room"`
	} `json:"pageinfo" mapstructure:"pageinfo"`
	Result struct {
		LiveUser map[string]any `json:"live_user" mapstructure:"live_user"`
		LiveRoom map[string]any `json:"live_room" mapstructure:"live_room"`
	} `json:"result" mapstructure:"result"`
}

type videoPlayurlDurl []struct {
	Order     int      `json:"order" mapstructure:"order"`
	Length    int      `json:"length" mapstructure:"length"`
	Size      int      `json:"size" mapstructure:"size"`
	Ahead     string   `json:"ahead" mapstructure:"ahead"`
	Vhead     string   `json:"vhead" mapstructure:"vhead"`
	Url       string   `json:"url" mapstructure:"url"`
	BackupUrl []string `json:"backup_url" mapstructure:"backup_url"`
}

type videoPlayurlDashInfo struct {
	Id        int      `json:"id" mapstructure:"id"` // 80
	BaseUrl   string   `json:"base_url" mapstructure:"base_url"`
	BackupUrl []string `json:"backup_url" mapstructure:"backup_url"`
	Bandwidth int      `json:"bandwidth" mapstructure:"bandwidth"`

	// "video/mp4", "audio/mp4"
	MimeType string `json:"mime_type" mapstructure:"mime_type"`

	// [VideoCodec]
	// "avc1.640032", "hev1.1.6.L150.90", "av01.0.00M.10.0.110.01.01.01.0"
	// "mp4a.40.2"
	Codecs string `json:"codecs" mapstructure:"codecs"`

	Width     int    `json:"width" mapstructure:"width"`           // 1920
	Height    int    `json:"height" mapstructure:"height"`         // 1080
	FrameRate string `json:"frame_rate" mapstructure:"frame_rate"` // "29.970"

	Sar string `json:"sar" mapstructure:"sar"` // "1:1"

	StartWithSap int `json:"start_with_sap" mapstructure:"start_with_sap"` // 1

	SegmentBase struct {
		Initialization string `json:"initialization" mapstructure:"initialization"`
		IndexRange     string `json:"index_range" mapstructure:"index_range"`
	} `json:"segment_base" mapstructure:"segment_base"`

	Codecid int `json:"codecid" mapstructure:"codecid"` // 7
}

type videoPlayurlDash struct {
	Duration      int     `json:"duration" mapstructure:"duration"`
	MinBufferTime float64 `json:"min_buffer_time" mapstructure:"min_buffer_time"`

	Video []videoPlayurlDashInfo `json:"video" mapstructure:"video"`
	Audio []videoPlayurlDashInfo `json:"audio" mapstructure:"audio"`

	Dolby struct {
		Type int `json:"type" mapstructure:"type"`
		// Audio *struct{} `json:"audio" mapstructure:"audio"`
	} `json:"dolby" mapstructure:"dolby"` // "dolby" : {}

	// Flac *struct{} `json:"flac" mapstructure:"flac"` // "flac" : {}
}

type videoPlayurlSupportFormat struct {
	Quality        int      `json:"quality" mapstructure:"quality"`                 // 80
	Format         string   `json:"format" mapstructure:"format"`                   // "mp4"
	NewDescription string   `json:"new_description" mapstructure:"new_description"` // "高清 720P"
	DisplayDesc    string   `json:"display_desc" mapstructure:"display_desc"`       // "高清 720P"
	Superscript    string   `json:"superscript" mapstructure:"superscript"`         // "4K"
	Codecs         []string `json:"codecs" mapstructure:"codecs"`                   // ["avc1.640032", "mp4a.40.2"]
}

// "data"
type VideoPlayurl struct {
	From              string   `json:"from" mapstructure:"from"`                             // "bili"
	Result            string   `json:"result" mapstructure:"result"`                         // "success"
	Message           string   `json:"message" mapstructure:"message"`                       // "ok"
	Quality           int      `json:"quality" mapstructure:"quality"`                       // 80
	Format            string   `json:"format" mapstructure:"format"`                         // "mp4"
	Timelength        int      `json:"timelength" mapstructure:"timelength"`                 // 0
	AcceptFormat      string   `json:"accept_format" mapstructure:"accept_format"`           // "mp4720,mp4"
	AcceptDescription []string `json:"accept_description" mapstructure:"accept_description"` // ["高清 720P", "流畅 360P"]
	AcceptQuality     []int    `json:"accept_quality" mapstructure:"accept_quality"`         // [64, 16]
	VideoCodecid      int      `json:"video_codecid" mapstructure:"video_codecid"`           // 7
	SeekParam         string   `json:"seek_param" mapstructure:"seek_param"`                 // "v"
	SeekType          string   `json:"seek_type" mapstructure:"seek_type"`                   // "play"

	Durl videoPlayurlDurl `json:"durl" mapstructure:"durl"` // mp4 流地址

	Dash *videoPlayurlDash `json:"dash" mapstructure:"dash"`

	SupportFormats []videoPlayurlSupportFormat `json:"support_formats" mapstructure:"support_formats"` // "support_formats" : {}

	// HighFormat *struct{} `json:"high_format" mapstructure:"high_format"`
	LastPlayTime int `json:"last_play_time" mapstructure:"last_play_time"` // (s)
	LastPlayCid  int `json:"last_play_cid" mapstructure:"last_play_cid"`
	// ViewInfo *struct{} `json:"view_info" mapstructure:"view_info"`
	PlayConf struct {
		IsNewDescription bool `json:"is_new_description" mapstructure:"is_new_description"` // 是否使用新描述
	}
}

// "data"
type ReplyList struct {
	Page struct {
		Num    int `json:"num" mapstructure:"num"`
		Size   int `json:"size" mapstructure:"size"`
		Count  int `json:"count" mapstructure:"count"`
		Acount int `json:"acount" mapstructure:"acount"`
	} `json:"page" mapstructure:"page"`
	// Config struct{} `json:"config" mapstructure:"config"`
	// Replies []struct{} `json:"replies" mapstructure:"replies"`
	// TopReplies []struct{} `json:"top_replies" mapstructure:"top_replies"`
	Upper struct {
		Mid int `json:"mid" mapstructure:"mid"` // uid
		Top struct {
			Member struct {
				Mid    string `json:"mid" mapstructure:"mid"`
				Uname  string `json:"uname" mapstructure:"uname"`
				Sex    string `json:"sex" mapstructure:"sex"`
				Sign   string `json:"sign" mapstructure:"sign"`
				Avatar string `json:"avatar" mapstructure:"avatar"`
			} `json:"member" mapstructure:"member"`
			Content struct {
				Message string `json:"message" mapstructure:"message"`
				// Members []any `json:"members" mapstructure:"members"`
				// Emote struct{} `json:"emote" mapstructure:"emote"`
				// JumpUrl struct{} `json:"jump_url" mapstructure:"jump_url"`
				MaxLine int `json:"max_line" mapstructure:"max_line"`
			} `json:"content" mapstructure:"content"`
		} `json:"top" mapstructure:"top"`
		Vote *VoteInfo `json:"vote" mapstructure:"vote"` // ?TODO: complete
	} `json:"upper" mapstructure:"upper"` // 置顶
}

// "data"
type QrcodeGenerate struct {
	Url       string `json:"url" mapstructure:"url"`               // 二维码内容
	QrcodeKey string `json:"qrcode_key" mapstructure:"qrcode_key"` // 轮询时 query
}

// "data"
type QrcodePoll struct {
	Url          string `json:"url" mapstructure:"url"`
	RefreshToken string `json:"refresh_token" mapstructure:"refresh_token"`
	Timestamp    int    `json:"timestamp" mapstructure:"timestamp"`
	Code         int    `json:"code" mapstructure:"code"` // [LoginCodeState]
	Message      string `json:"message" mapstructure:"message"`
}

// "data"
type Nav struct {
	IsLogin       bool   `json:"isLogin" mapstructure:"isLogin"`
	EmailVerified int    `json:"email_verified" mapstructure:"email_verified"`
	Face          string `json:"face" mapstructure:"face"`
	FaceNft       int    `json:"face_nft" mapstructure:"face_nft"`
	FaceNftType   int    `json:"face_nft_type" mapstructure:"face_nft_type"`
	LevelInfo     struct {
		CurrentLevel int    `json:"current_level" mapstructure:"current_level"`
		CurrentMin   int    `json:"current_min" mapstructure:"current_min"`
		CurrentExp   int    `json:"current_exp" mapstructure:"current_exp"`
		NextExp      string `json:"next_exp" mapstructure:"next_exp"`
	} `json:"level_info" mapstructure:"level_info"`
	Mid            int     `json:"mid" mapstructure:"mid"`
	MobileVerified int     `json:"mobile_verified" mapstructure:"mobile_verified"`
	Money          float64 `json:"money" mapstructure:"money"`
	Moral          int     `json:"moral" mapstructure:"moral"`
	Official       struct {
		Role  int    `json:"role" mapstructure:"role"`
		Title string `json:"title" mapstructure:"title"`
		Desc  string `json:"desc" mapstructure:"desc"`
		Type  int    `json:"type" mapstructure:"type"`
	} `json:"official" mapstructure:"official"`
	OfficialVerify struct {
		Type int    `json:"type" mapstructure:"type"`
		Desc string `json:"desc" mapstructure:"desc"`
	} `json:"officialVerify" mapstructure:"officialVerify"`
	Pendant struct {
		Pid               int    `json:"pid" mapstructure:"pid"`
		Name              string `json:"name" mapstructure:"name"`
		Image             string `json:"image" mapstructure:"image"`
		Expire            int    `json:"expire" mapstructure:"expire"`
		ImageEnhance      string `json:"image_enhance" mapstructure:"image_enhance"`
		ImageEnhanceFrame string `json:"image_enhance_frame" mapstructure:"image_enhance_frame"`
		NPid              int    `json:"n_pid" mapstructure:"n_pid"`
	} `json:"pendant" mapstructure:"pendant"`
	Scores       int    `json:"scores" mapstructure:"scores"`
	Uname        string `json:"uname" mapstructure:"uname"`
	VipDueDate   int64  `json:"vipDueDate" mapstructure:"vipDueDate"`
	VipStatus    int    `json:"vipStatus" mapstructure:"vipStatus"`
	VipType      int    `json:"vipType" mapstructure:"vipType"`
	VipPayType   int    `json:"vip_pay_type" mapstructure:"vip_pay_type"`
	VipThemeType int    `json:"vip_theme_type" mapstructure:"vip_theme_type"`
	VipLabel     struct {
		Path                  string `json:"path" mapstructure:"path"`
		Text                  string `json:"text" mapstructure:"text"`
		LabelTheme            string `json:"label_theme" mapstructure:"label_theme"`
		TextColor             string `json:"text_color" mapstructure:"text_color"`
		BgStyle               int    `json:"bg_style" mapstructure:"bg_style"`
		BgColor               string `json:"bg_color" mapstructure:"bg_color"`
		BorderColor           string `json:"border_color" mapstructure:"border_color"`
		UseImgLabel           bool   `json:"use_img_label" mapstructure:"use_img_label"`
		ImgLabelUriHans       string `json:"img_label_uri_hans" mapstructure:"img_label_uri_hans"`
		ImgLabelUriHant       string `json:"img_label_uri_hant" mapstructure:"img_label_uri_hant"`
		ImgLabelUriHansStatic string `json:"img_label_uri_hans_static" mapstructure:"img_label_uri_hans_static"`
		ImgLabelUriHantStatic string `json:"img_label_uri_hant_static" mapstructure:"img_label_uri_hant_static"`
	} `json:"vip_label" mapstructure:"vip_label"`
	VipAvatarSubscript int    `json:"vip_avatar_subscript" mapstructure:"vip_avatar_subscript"`
	VipNicknameColor   string `json:"vip_nickname_color" mapstructure:"vip_nickname_color"`
	Vip                struct {
		Type       int   `json:"type" mapstructure:"type"`
		Status     int   `json:"status" mapstructure:"status"`
		DueDate    int64 `json:"due_date" mapstructure:"due_date"`
		VipPayType int   `json:"vip_pay_type" mapstructure:"vip_pay_type"`
		ThemeType  int   `json:"theme_type" mapstructure:"theme_type"`
		Label      struct {
			Path                  string `json:"path" mapstructure:"path"`
			Text                  string `json:"text" mapstructure:"text"`
			LabelTheme            string `json:"label_theme" mapstructure:"label_theme"`
			TextColor             string `json:"text_color" mapstructure:"text_color"`
			BgStyle               int    `json:"bg_style" mapstructure:"bg_style"`
			BgColor               string `json:"bg_color" mapstructure:"bg_color"`
			BorderColor           string `json:"border_color" mapstructure:"border_color"`
			UseImgLabel           bool   `json:"use_img_label" mapstructure:"use_img_label"`
			ImgLabelUriHans       string `json:"img_label_uri_hans" mapstructure:"img_label_uri_hans"`
			ImgLabelUriHant       string `json:"img_label_uri_hant" mapstructure:"img_label_uri_hant"`
			ImgLabelUriHansStatic string `json:"img_label_uri_hans_static" mapstructure:"img_label_uri_hans_static"`
			ImgLabelUriHantStatic string `json:"img_label_uri_hant_static" mapstructure:"img_label_uri_hant_static"`
		} `json:"label" mapstructure:"label"`
		AvatarSubscript    int    `json:"avatar_subscript" mapstructure:"avatar_subscript"`
		NicknameColor      string `json:"nickname_color" mapstructure:"nickname_color"`
		Role               int    `json:"role" mapstructure:"role"`
		AvatarSubscriptUrl string `json:"avatar_subscript_url" mapstructure:"avatar_subscript_url"`
		TvVipStatus        int    `json:"tv_vip_status" mapstructure:"tv_vip_status"`
		TvVipPayType       int    `json:"tv_vip_pay_type" mapstructure:"tv_vip_pay_type"`
		TvDueDate          int    `json:"tv_due_date" mapstructure:"tv_due_date"`
	} `json:"vip" mapstructure:"vip"`
	Wallet struct {
		Mid           int `json:"mid" mapstructure:"mid"`
		BcoinBalance  int `json:"bcoin_balance" mapstructure:"bcoin_balance"`
		CouponBalance int `json:"coupon_balance" mapstructure:"coupon_balance"`
		CouponDueTime int `json:"coupon_due_time" mapstructure:"coupon_due_time"`
	} `json:"wallet" mapstructure:"wallet"`
	HasShop        bool   `json:"has_shop" mapstructure:"has_shop"`
	ShopUrl        string `json:"shop_url" mapstructure:"shop_url"`
	AllowanceCount int    `json:"allowance_count" mapstructure:"allowance_count"`
	AnswerStatus   int    `json:"answer_status" mapstructure:"answer_status"`
	IsSeniorMember int    `json:"is_senior_member" mapstructure:"is_senior_member"`
	WbiImg         struct {
		ImgUrl string `json:"img_url" mapstructure:"img_url"`
		SubUrl string `json:"sub_url" mapstructure:"sub_url"`
	} `json:"wbi_img" mapstructure:"wbi_img"`
	IsJury bool `json:"is_jury" mapstructure:"is_jury"`
}
