package biligo

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// FormatDynamic 格式化动态,
//
// .Get("data.item")
//
// 转发原动态 .Get("orig")
//
// Deprecated: use [DynamicDetail.DoTemplate] instead
func FormatDynamic(j gjson.Result) string {
	dynamic := j.Get("modules.module_dynamic")           // 动态主体
	id := j.Get("id_str").String()                       // 动态 id
	name := j.Get("modules.module_author.name").String() // 发布者用户名

	if dynamic.Get("major.blocked").Exists() { // 充电专属动态
		hintMsg := dynamic.Get("major.blocked.hint_message").String() // 提示信息
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：
	<访问受限>
	%s`,
			id,
			name,
			hintMsg,
		)
	}

	action := j.Get("modules.module_author.pub_action").String() // "投稿了视频"/"发布了动态视频"/"投稿了文章"/"直播了"...

	topic := "" // 话题, 可能为空
	if dynamic.Get("topic.name").Exists() {
		topic = "\n#️⃣" + dynamic.Get("topic.name").String() + "#️⃣"
	}

	text := "" // 正文, 可能为空
	if dynamic.Get("desc.text").String() != "" {
		text = "\n" + dynamic.Get("desc.text").String()
	}

	additional := "" // 附加子内容, 可能为空
	addiType := dynamic.Get("additional.type").String()
	switch addiType {
	case "ADDITIONAL_TYPE_RESERVE": // 预约
		reserveMain := dynamic.Get("additional.reserve")
		additional = fmt.Sprintf(`

%s
%s
%s`,
			reserveMain.Get("title").String(),      // "视频预约："/"直播预约："
			reserveMain.Get("desc1.text").String(), // "预计xxx发布"/"01-25 19:55 直播"
			reserveMain.Get("desc2.text").String(), // "xx人预约"/"xx观看
		)

	case "ADDITIONAL_TYPE_VOTE": // 投票
		vote, err := FetchVoteInfo(dynamic.Get("additional.vote.vote_id").String())
		if err != nil {
			additional = fmt.Sprintf("投票信息获取失败：%v", err)
			break
		}
		additional = vote.DoTemplate()

	case "ADDITIONAL_TYPE_UGC": // 视频跳转
		url := dynamic.Get("additional.ugc.jump_url").String()
		r, err := ParseLink(url)
		if err != nil {
			additional = fmt.Sprintf("视频信息获取失败：%v", err)
		} else if len(r) == 0 {
			additional = "视频信息获取失败: no result"
		} else {
			video, _, err := FetchVideoInfo2(id)
			if err != nil {
				additional = fmt.Sprintf("视频信息获取失败：%v", err)
			} else {
				additional = video.DoTemplate()
			}
		}
	}

	dynamicType := j.Get("type").String() // 动态类型
	switch dynamicType {
	case "DYNAMIC_TYPE_FORWARD": // 转发
		action = "转发动态"                             // 手动挡
		origDynamic := FormatDynamic(j.Get("orig")) // 原动态
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s
	
	%s`,
			id,
			name, action, topic, text,
			origDynamic,
		)

	case "DYNAMIC_TYPE_NONE": // 转发的动态已被删除
		return dynamic.Get("major.none.tips").String() // "源动态已被作者删除"

	case "DYNAMIC_TYPE_WORD": // 纯文字
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s%s`,
			id,
			name, action, topic, text, additional,
		)

	case "DYNAMIC_TYPE_DRAW": // 图文
		draw := dynamic.Get("major.draw")
		images := strings.Builder{}
		for _, item := range draw.Get("items").Array() {
			images.WriteString("[CQ:image,file=")
			images.WriteString(item.Get("src").String())
			images.WriteString("]")
		}
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s
	%s%s`,
			id,
			name, action, topic, text,
			images.String(), additional,
		)

	case "DYNAMIC_TYPE_AV": // 视频
		archive := dynamic.Get("major.archive")
		if dynamic.Get("desc.text").Exists() && text[1:] == archive.Get("desc").String() {
			text = "" // 如果正文和简介相同, 不显示正文
		}
		aid := archive.Get("aid").String() // av号数字
		var content string
		video, _, err := FetchVideoInfo2(aid)
		if err != nil {
			content = err.Error()
		} else {
			content = video.DoTemplate()
		}
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s
	
	%s%s`,
			id,
			name, action, topic, text,
			content, additional,
		)

	case "DYNAMIC_TYPE_ARTICLE": // 文章
		cvid := dynamic.Get("major.article.id").String() // cv号数字
		content := ""
		article, err := FetchArticleInfo(cvid)
		if err != nil {
			content = fmt.Sprintf("文章信息获取失败：%v", err)
		} else {
			content = article.DoTemplate()
		}
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s
	
	%s%s`,
			id,
			name, action, topic, text,
			content, additional,
		)

	case "DYNAMIC_TYPE_MUSIC": // 音频
		sid := dynamic.Get("major.music").Get("id").String()
		content := ""
		song, err := FetchSong(sid)
		if err != nil {
			content = fmt.Sprintf("音频信息获取失败：%v", err)
		} else {
			content = song.DoTemplate()
		}
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s
	
	%s%s`,
			id,
			name, action, topic, text,
			content, additional,
		)

	case "DYNAMIC_TYPE_LIVE", "DYNAMIC_TYPE_LIVE_RCMD": // 直播间分享, 直播开播(动态流拿不到更新)
		uid := j.Get("modules.module_author.mid").String() // 发布者 uid
		content := ""
		live, err := FetchLiveStatus(uid)
		if err != nil {
			content = fmt.Sprintf("直播间信息获取失败：%v", err)
		} else {
			content = live.Get(uid).DoTemplate()
		}
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s
	
	%s%s`,
			id,
			name, action, topic, text,
			content, additional,
		)

	case "DYNAMIC_TYPE_COMMON_SQUARE": // 装扮 / 剧集点评 / 普通分享
		text := dynamic.Get("desc.text").String()       // 正文
		decCover := dynamic.Get("major.cover").String() // 装扮封面
		decTitle := dynamic.Get("major.title").String() // 装扮标题
		decDesc := dynamic.Get("major.desc").String()   // 装扮描述
		decUrl := dynamic.Get("major.url").String()     // 装扮链接
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s
	
	%s
	%s
	%s
	%s`,
			id,
			name, action, topic, text,
			decCover,
			decTitle,
			decDesc,
			decUrl,
		)

	default:
		return fmt.Sprintf(
			`t.bilibili.com/%s
	%s：%s%s%s
	
	未知的动态类型：%s`,
			id,
			name, action, topic, text,
			dynamicType,
		)
	}
}
