package biligo

import (
	"strings"
)

func ReqVideoInfo(aid string) *Request {
	return NewGet(URL_VIDEO_INFO_WBI).WbiSign().
		WithQuery("aid", strings.TrimPrefix(aid, "av"))
}

func ReqVideoOnline(aid, cid string) *Request {
	return NewGet(URL_VIDEO_ONLINE).
		WithQuerys("aid", strings.TrimPrefix(aid, "av"), "cid", cid)
}

func ReqVideoConclusion(aid, cid string) *Request {
	return NewGet(URL_VIDEO_CONCLUSION_WBI).WbiSign().
		WithQuerys("aid", strings.TrimPrefix(aid, "av"), "cid", cid)
}

func ReqMediaInfoBase(mdid string) *Request {
	return NewGet(URL_MEDIA_INFO_BASE).
		WithQuery("media_id", strings.TrimPrefix(mdid, "md"))
}

func ReqMediaInfoSsid(ssid string) *Request {
	return NewGet(URL_MEDIA_INFO_DETAIL).
		WithQuery("season_id", strings.TrimPrefix(ssid, "ss"))
}

func ReqMediaInfoEpid(epid string) *Request {
	return NewGet(URL_MEDIA_INFO_DETAIL).
		WithQuery("ep_id", strings.TrimPrefix(epid, "ep"))
}

func ReqMediaSection(ssid string) *Request {
	return NewGet(URL_MEDIA_SECTION).
		WithQuery("season_id", strings.TrimPrefix(ssid, "ss"))
}

func ReqLiveStatus(uid ...string) *Request {
	return NewGet(URL_LIVE_STATUS_UIDS).
		WithQueryList("uids[]", uid...)
}

func ReqLiveRoomInfo(roomId string) *Request {
	return NewGet(URL_LIVE_INFO_ROOMID).
		WithQuery("room_id", roomId)
}

func ReqArticleInfo(id string) *Request {
	return NewGet(URL_ARTICLE_INFO).
		WithQuery("id", strings.TrimPrefix(id, "cv"))
}

func ReqAudio(sid string) (info, tag, member, lyric *Request) {
	sid = strings.TrimPrefix(sid, "au")
	return NewGet(URL_SONG_INFO).WithQuery("sid", sid),
		NewGet(URL_SONG_TAG).WithQuery("sid", sid),
		NewGet(URL_SONG_MEMBER).WithQuery("sid", sid),
		NewGet(URL_SONG_LIRIC).WithQuery("sid", sid)
}

func ReqSpaceCard(uid string) *Request {
	return NewGet(URL_SPACE_CARD).
		WithQuerys("mid", uid, "photo", "true")
}

func ReqRelationStat(uid string) *Request {
	return NewGet(URL_RELATION_STAT).
		WithQuery("vmid", uid)
}

func ReqDynamicAll() *Request {
	return NewGet(URL_DYNAMIC_ALL)
}

func ReqDynamicAllUpdate(updateBaseline string) *Request {
	return NewGet(URL_DYNAMIC_ALL_UPDATE).
		WithQuery("update_baseline", updateBaseline)
}

func ReqDynamicDetail(id string) *Request {
	return NewGet(URL_DYNAMIC_DETAIL).
		WithQuery("id", id)
}

func ReqVoteInfo(voteId string) *Request {
	return NewGet(URL_VOTE_INFO).
		WithQuery("vote_id", voteId)
}

func ReqSearchAll(keyword string) *Request {
	return NewGet(URL_SEARCH_ALL_WBI).WbiSign().
		WithQuery("keyword", keyword)
}

func ReqSearchType(searchType SearchClass, keyword string) *Request {
	return NewGet(URL_SEARCH_TYPE_WBI).WbiSign().
		WithQuerys("search_type", string(searchType), "keyword", keyword)
}

func ReqLiveDanmuInfo(roomId string) *Request {
	return NewGet(URL_LIVE_DANMU_INFO).
		WithQuery("id", roomId)
}

func ReqReplyList(typ CommentClass, oid string) *Request {
	return NewGet(URL_REPLY_LIST).
		WithQuerys("type", itoa(typ), "oid", oid)
}

func ReqVideoPlayurl(q ...string) *Request {
	return NewGet(URL_VIDEO_PALYURL_WBI).WbiSign().
		WithQuerys(q...)
}

func ReqLoginQrcodeGenerate() *Request {
	return NewGet(URL_LOGIN_QRCODE_GENERATE)
}

func ReqLoginQrcodePoll(qrcodeKey string) *Request {
	return NewGet(URL_LOGIN_QRCODE_POLL).
		WithQuery("qrcode_key", qrcodeKey)
}

func ReqNav() *Request {
	return NewGet(URL_NAV)
}
