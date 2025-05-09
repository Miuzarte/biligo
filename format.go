package biligo

import (
	"errors"
	"slices"

	"github.com/tidwall/gjson"
)

var (
	errSearchFailedVideo   = errors.New("视频信息获取失败")
	errSearchFailedMedia   = errors.New("剧集信息获取失败")
	errSearchFailedLive    = errors.New("直播间信息获取失败")
	errSearchFailedArticle = errors.New("文章信息获取失败")
	errSearchFailedSpace   = errors.New("空间信息获取失败")
	errSearchUnknownType   = errors.New("未知类型")
)

func FormatSearch(results []gjson.Result) []Templatable {
	ts := make([]Templatable, 0, len(results))
	for _, r := range results {
		switch typ := SearchClass(r.Get("type").String()); typ {
		case SEARCH_TYPE_VIDEO:
			vi, _, err := FetchVideoInfo2(r.Get("id").String())
			if err == nil {
				ts = append(ts, &vi)
			} else {
				ts = append(ts, wrapErr(errSearchFailedVideo, err).(*Error))
			}

		case SEARCH_TYPE_BANGUMI, SEARCH_TYPE_FT:
			m, err := FetchMediaInfoSsid(r.Get("season_id").String())
			if err == nil {
				ts = append(ts, &m)
			} else {
				ts = append(ts, wrapErr(errSearchFailedMedia, err).(*Error))
			}

		case SEARCH_TYPE_LIVE: // 实际结果中不存在这个类型, 仅在请求中使用
			fallthrough
		case SEARCH_TYPE_LIVE_USER, SEARCH_TYPE_LIVE_ROOM:
			uid := r.Get("uid").String()
			l, err := FetchLiveStatus(uid)
			if err == nil {
				ts = append(ts, l.Get(uid))
			} else {
				ts = append(ts, wrapErr(errSearchFailedLive, err).(*Error))
			}

		case SEARCH_TYPE_ARTICLE:
			a, err := FetchArticleInfo(r.Get("id").String())
			if err == nil {
				ts = append(ts, &a)
			} else {
				ts = append(ts, wrapErr(errSearchFailedArticle, err).(*Error))
			}

		case SEARCH_TYPE_BILI_USER:
			s, err := FetchSpaceCard(r.Get("mid").String())
			if err != nil {
				ts = append(ts, &s)
			} else {
				ts = append(ts, wrapErr(errSearchFailedSpace, err).(*Error))
			}

		default:
			ts = append(ts, wrapErr(errSearchUnknownType, typ).(*Error))
		}
	}
	return ts
}

func FormatSearchAll(j gjson.Result) []Templatable {
	var results []gjson.Result
	for _, r := range j.Array() {
		// 扁平化所有子类中的结果
		results = append(results, r.Get("data").Array()...)
	}
	return FormatSearch(results)
}

func FormatSearchType(j gjson.Result) []Templatable {
	if !j.IsArray() {
		// 可能请求的是 [SEARCH_TYPE_LIVE],
		// result 下存在 live_room 与 live_user
		return FormatSearch(slices.Concat(
			j.Get("live_user").Array(),
			j.Get("live_room").Array(),
		))
	}
	return FormatSearch(j.Array())
}
