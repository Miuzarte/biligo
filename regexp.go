package biligo

import (
	"regexp"
	"strings"
)

var (
	reg_short   = regexp.MustCompile(REGEXP_SHORT)
	reg_video   = regexp.MustCompile(REGEXP_VIDEO)
	reg_bangumi = regexp.MustCompile(REGEXP_BANGUMI)
	reg_live    = regexp.MustCompile(REGEXP_LIVE)
	reg_article = regexp.MustCompile(REGEXP_ARTICLE)
	reg_audio   = regexp.MustCompile(REGEXP_AUDIO)
	reg_space   = regexp.MustCompile(REGEXP_SPACE)
	reg_dynamic = regexp.MustCompile(REGEXP_DYNAMIC)
)

type ParseResult struct {
	Type    LinkType
	Content string
}

func (r ParseResult) FetchFormat() (string, error) {
	switch r.Type {
	case LINK_TYPE_ARCHIVE:
		video, _, err := FetchVideoInfo2(r.Content)
		if err != nil {
			return "", err
		}
		return video.DoTemplate(), nil

	case LINK_TYPE_MEDIA:
		switch r.Content[:2] {
		case "md":
			// mdid 获取的信息太少
			// fallthrough 到 ssid
			mb, err := FetchMediaInfoBase(r.Content)
			if err != nil {
				return "", err
			}
			if mb.Media.SeasonId == 0 {
				return "", wrapErr(ErrParseMediaNoSsid, mb)
			}
			r.Content = "ss" + itoa(mb.Media.SeasonId)
			fallthrough

		case "ss":
			m, err := FetchMediaInfoSsid(r.Content)
			if err != nil {
				return "", err
			}
			return m.DoTemplate(), nil

		case "ep":
			m, err := FetchMediaInfoEpid(r.Content)
			if err != nil {
				return "", err
			}
			return m.DoTemplate(), nil

		default:
			panic("unreachable, fix regexp")
		}

	case LINK_TYPE_LIVE: // id 为直播间 id, 需要先获取 uid
		room, err := FetchLiveRoomInfo(r.Content)
		if err != nil {
			return "", err
		}
		if room.Uid == 0 {
			return "", wrapErr(ErrParseLiveNoUid, room)
		}
		uidStr := itoa(room.Uid)
		live, err := FetchLiveStatus(uidStr)
		if err != nil {
			return "", err
		}
		return live.Get(uidStr).DoTemplate(), nil

	case LINK_TYPE_ARTICLE:
		article, err := FetchArticleInfo(r.Content)
		if err != nil {
			return "", err
		}
		return article.DoTemplate(), nil

	case LINK_TYPE_SONG:
		song, err := FetchSong(r.Content)
		if err != nil {
			return "", err
		}
		return song.DoTemplate(), nil

	case LINK_TYPE_SPACE:
		space, err := FetchSpaceCard(r.Content)
		if err != nil {
			return "", err
		}
		return space.DoTemplate(), nil

	case LINK_TYPE_DYNAMIC:
		// 同时请求 desktop api 补充正文
		dd, err := FetchDynamicDetailFix(r.Content)
		if err != nil {
			return "", err
		}
		return dd.DoTemplate(), nil

	default:
		return "<unknonw type of link>", nil
	}
}

// ParseLink 解析多个链接并去重,
// 返回错误的原因只有解析短链时的网络问题
func ParseLink(s string) (results []ParseResult, err error) {
	if m := reg_short.FindAllStringSubmatch(s, -1); len(m) > 0 {
		var loc string
		var locs []string = make([]string, 0, len(m))
		var r []ParseResult

		for _, v := range m {
			loc, err = fetchLocation(v[REGEXP_INDEX_SHORT_URL])
			if err != nil {
				return nil, err
			}
			locs = append(locs, loc)
		}

		r, err = ParseLink(strings.Join(locs, "\n"))
		if err != nil {
			return nil, err
		}
		results = append(results, r...)
	}
	for _, v := range reg_video.FindAllStringSubmatch(s, -1) {
		id, err := AnyToAid(v[REGEXP_INDEX_ARCHIVE]) // 去重用
		if err != nil {
			id = v[REGEXP_INDEX_ARCHIVE] // fallback
		}
		results = append(results, ParseResult{LINK_TYPE_ARCHIVE, id})
	}
	for _, v := range reg_bangumi.FindAllStringSubmatch(s, -1) {
		results = append(results, ParseResult{LINK_TYPE_MEDIA, v[REGEXP_INDEX_BANGUMI_ID]})
	}
	for _, v := range reg_live.FindAllStringSubmatch(s, -1) {
		results = append(results, ParseResult{LINK_TYPE_LIVE, v[REGEXP_INDEX_LIVE]})
	}
	for _, v := range reg_article.FindAllStringSubmatch(s, -1) {
		results = append(results, ParseResult{LINK_TYPE_ARTICLE, v[REGEXP_INDEX_ARTICLE]})
	}
	for _, v := range reg_audio.FindAllStringSubmatch(s, -1) {
		results = append(results, ParseResult{LINK_TYPE_SONG, v[REGEXP_INDEX_AUDIO]})
	}
	for _, v := range reg_space.FindAllStringSubmatch(s, -1) {
		results = append(results, ParseResult{LINK_TYPE_SPACE, v[REGEXP_INDEX_SPACE]})
	}
	for _, v := range reg_dynamic.FindAllStringSubmatch(s, -1) {
		results = append(results, ParseResult{LINK_TYPE_DYNAMIC, v[REGEXP_INDEX_DYNAMIC]})
	}

	set := make(set[ParseResult])
	return set.clean(results), nil
}
