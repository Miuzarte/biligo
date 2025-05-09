package biligo

import (
	"testing"
)

func TestParseLink(t *testing.T) {
	const msg = `[CQ:json,data={"app":"com.tencent.tuwen.lua","bizsrc":"qqconnect.sdkshare","config":{"ctime":1738634984,"forward":1,"token":"422f382ad8dfd3c3048d395730ab3bc3","type":"normal"},"extra":{"app_type":1,"appid":100951776,"msg_seq":7467380391885968581,"uin":1428259430},"meta":{"news":{"app_type":1,"appid":100951776,"ctime":1738634984,"desc":"介 绍《战地风云》实验室 | 《战地风云》工作室","jumpUrl":"https://www.bilibili.com/video/BV1dQPdeJEh5?buvid=XU5F80C57B990E131F7997614188FD80D4ADE&from_spmid=creation.hot-tab.0.0&is_story_h5=false&mid=HMNjDlolKi4H02UHYbWVBA%3D%3D&p=1&plat_id=116&share_from=ugc&share_medium=android&share_plat=android&share_session_id=54aefd73-95fb-4625-b33a-7fb9336b17bf&share_source=QQ&share_tag=s_i&spmid=main.ugc-video-detail.0.0&timestamp=1738634981&unique_k=M7W9KyY&up_id=381399962&bbid=XU5F80C57B990E131F7997614188FD80D4ADE&ts=1738634980823","preview":"https://pic.ugcimg.cn/25f6929d1a76a1d01cf447b8d5bb6921/jpg1","tag":"哔哩哔哩","tagIcon":"https://open.gtimg.cn/open/app_icon/00/95/17/76/100951776_100_m.png?t=1737540523","title":"哔哩哔哩","uin":1428259430}},"prompt":"[分享]哔哩哔哩","ver":"0.0.0.1","view":"news"}]`
	results, err := ParseLink(msg)
	if err != nil {
		t.Error(err)
	}
	for _, r := range results {
		t.Log("type:", r.Type)
		t.Log("content:", r.Content)
	}
}
