package biligo

import (
	"fmt"
)

// SearchFormatAuto 在 searchType 为空字符串时使用综合搜索,
// 函数的返回结果可被断言为: [*VideoInfo], [*Media], [*LiveStatusUid], [*ArticleInfo], [*SpaceCard]
func SearchFormatAuto(searchType SearchClass, keyword string) ([]Templatable, error) {
	if searchType == "" {
		return SearchFormatAll(keyword)
	} else {
		return SearchFormatType(searchType, keyword)
	}
}

func SearchFormatAll(keyword string) ([]Templatable, error) {
	s, err := FetchSearchAll(keyword)
	if err != nil {
		return nil, err
	}
	return FormatSearchAll(s), nil
}

func SearchFormatType(searchType SearchClass, keyword string) ([]Templatable, error) {
	if searchType != "" && searchType.String() == "" {
		return nil, fmt.Errorf("unknown search type: %s", searchType)
	} else if searchType == "" {
		return nil, wrapErr(ErrSearchEmptyType, nil)
	}
	s, err := FetchSearchType(searchType, keyword)
	if err != nil {
		return nil, err
	}
	return FormatSearchType(s), nil
}
