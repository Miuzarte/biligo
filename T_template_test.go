package biligo

import (
	"os"
	"testing"
)

func TestTemplateAll(t *testing.T) {
	var vi VideoInfo
	var vc VideoConclusion
	var mb MediaBase
	var m Media
	var ls LiveStatus
	var lri LiveRoomInfo
	var ai ArticleInfo
	var s Song
	var sc SpaceCard

	var dar dynamicAddiReserve
	var dav dynamicAddiVote
	var dau dynamicAddiUgc
	var dac dynamicAddiCommon
	var dmd dynamicMajorDraw
	var dma dynamicMajorArchive
	var dmar dynamicMajorArticle
	var dmm dynamicMajorMusic
	var dml dynamicMajorLive
	var dmlr dynamicMajorLiveRcmd
	dmlr.Content = "{\"type\":1,\"live_play_info\":{\"online\":65235,\"area_id\":371,\"live_id\":\"608106002646663537\",\"room_paid_type\":0,\"cover\":\"https://i0.hdslb.com/bfs/live/new_room_cover/d314aa7293443915960df60abaa4cc7d276f67c9.jpg\",\"parent_area_name\":\"虚拟主播\",\"live_start_time\":1746531658,\"link\":\"//live.bilibili.com/25788785?live_from=85002\",\"parent_area_id\":9,\"watched_show\":{\"num\":2471,\"text_small\":\"2471\",\"text_large\":\"2471人看过\",\"icon\":\"https://i0.hdslb.com/bfs/live/a725a9e61242ef44d764ac911691a7ce07f36c1d.png\",\"icon_location\":\"\",\"icon_web\":\"https://i0.hdslb.com/bfs/live/8d9d0f33ef8bf6f308742752d13dd0df731df19c.png\",\"switch\":true},\"live_status\":1,\"room_type\":0,\"play_type\":0,\"area_name\":\"虚拟日常\",\"pendants\":{\"list\":{\"mobile_index_badge\":{\"list\":{\"1\":{\"name\":\"百人成就\",\"position\":1,\"text\":\"\",\"bg_color\":\"#FB9E60\",\"bg_pic\":\"https://i0.hdslb.com/bfs/live/539ce26c45cd4019f55b64cfbcedc3c01820e539.png\",\"pendant_id\":426,\"type\":\"mobile_index_badge\"},\"2\":{\"type\":\"mobile_index_badge\",\"name\":\"长红计划-SSS\",\"position\":2,\"text\":\"璀璨之星\",\"bg_color\":\"#FB9E60\",\"bg_pic\":\"http://i0.hdslb.com/bfs/live/65ef92c9e9f8fafacd29e58efd34481993aa0e26.png\",\"pendant_id\":1728}}},\"index_badge\":{\"list\":{\"1\":{\"bg_pic\":\"https://i0.hdslb.com/bfs/live/539ce26c45cd4019f55b64cfbcedc3c01820e539.png\",\"pendant_id\":425,\"type\":\"index_badge\",\"name\":\"百人成就\",\"position\":1,\"text\":\"\",\"bg_color\":\"#FB9E60\"}}}}},\"room_id\":25788785,\"uid\":1954091502,\"title\":\"小岁来也！\",\"live_screen_type\":0},\"live_record_info\":null}"
	var dmp dynamicMajorPgc
	var dd DynamicDetail
	var vti VoteInfo

	// var sa SearchAll
	// var st SearchType
	// var stl SearchTypeLive

	t.Log(vi.DoTemplate())
	t.Log(vc.DoTemplate())
	t.Log(mb.DoTemplate())
	t.Log(m.DoTemplate())
	t.Log(ls.DoTemplate())
	t.Log(lri.DoTemplate())
	t.Log(ai.DoTemplate())
	t.Log(s.DoTemplate())
	t.Log(sc.DoTemplate())

	t.Log(dar.DoTemplate())
	t.Log(dav.DoTemplate())
	t.Log(dau.DoTemplate())
	t.Log(dac.DoTemplate())
	t.Log(dmd.DoTemplate())
	t.Log(dma.DoTemplate())
	t.Log(dmar.DoTemplate())
	t.Log(dmm.DoTemplate())
	t.Log(dml.DoTemplate())
	t.Log(dmlr.DoTemplate())
	t.Log(dmp.DoTemplate())
	t.Log(dd.DoTemplate())
	t.Log(vti.DoTemplate())

	// t.Log(sa.DoTemplate())
	// t.Log(st.DoTemplate())
	// t.Log(stl.DoTemplate())
}

func TestTemplateDefined(t *testing.T) {
	t.Log(typeTemplate.DefinedTemplates())
	typeTemplate.Execute(os.Stdout, nil)
}

func TestTemplateEmbed(t *testing.T) {
	t.Log(typeTemplates.ReadDir("."))
	t.Log(typeTemplates.ReadDir("template"))
	t.Log(typeTemplates.ReadDir("template/"))
	t.Log(typeTemplates.ReadDir("template/type"))
	t.Log(typeTemplates.ReadDir("template/type/"))
}
