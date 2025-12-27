package biligo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	netUrl "net/url"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
)

var DefaultHeaders = map[string]string{
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
	"Dnt":                       "1",
	"Host":                      "www.bilibili.com",
	"Origin":                    "https://www.bilibili.com",
	"Referer":                   "https://www.bilibili.com/",
	"Priority":                  "u=1, i",
	"Sec-Ch-Ua":                 "\"Chromium\";v=\"141\", \"Microsoft Edge\";v=\"141\", \"Not.A/Brand\";v=\"99\"",
	"Sec-Ch-Ua-Mobile":          "?0",
	"Sec-Ch-Ua-Platform":        "\"Windows\"",
	"Sec-Fetch-Dest":            "document",
	"Sec-Fetch-Mode":            "navigate",
	"Sec-Fetch-Site":            "same-site",
	"Sec-Fetch-User":            "?1",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36 Edg/141.0.0.0",
}

var httpClient = newClient()

func newClient() client {
	c := client{
		Client: &http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar: cookie,
		},
		Headers: http.Header{},
	}
	for k, v := range DefaultHeaders {
		c.Headers.Add(k, v)
	}
	return c
}

func (c *client) Do(req *Request) (*Response, error) {
	if req == nil {
		return nil, wrapErr(ErrClientNilRequest, nil)
	}
	maps.Copy(req.Header, c.Headers)
	if len(*req.querys) > 0 {
		req.URL.RawQuery = req.querys.Encode()
	}

	resp, err := c.Client.Do(req.Std())
	if err != nil {
		return nil, err
	}
	return &Response{Response: resp}, nil
}

var identity Identity

type Identity struct {
	Cookie       string `json:"cookie"`
	RefreshToken string `json:"refresh_token"`
	Uid          int    `json:"uid"`
}

var cookie = &cookieManager{}

type cookieManager struct {
	cookie *http.Cookie
	// guestFetching atomic.Bool
	guestFetching sync.Mutex
}

func (c *cookieManager) fetchGuestCookie() error {
	// defer c.guestFetching.Store(false)
	defer c.guestFetching.Unlock()

	// 获取过程中卸下 cookie jar
	jar := httpClient.Client.Jar
	httpClient.Client.Jar = nil
	defer func() {
		httpClient.Client.Jar = jar
	}()

	resp, err := httpClient.Do(NewGet(URL_MAIN_PAGE))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = c.set(joinCookie(resp.Header["Set-Cookie"]))
	if err != nil {
		return err
	}

	return nil
}

func (c *cookieManager) set(setCookie string) (err error) {
	cookie, err := http.ParseSetCookie(setCookie)
	if err != nil {
		return err
	}
	c.cookie = cookie
	return nil
}

func (c *cookieManager) SetCookies(_ *netUrl.URL, _ []*http.Cookie) {
	// only for implementation of [http.CookieJar] currently
}

func (c *cookieManager) Cookies(u *netUrl.URL) []*http.Cookie {
	if c.cookie == nil {
		// 没有 cookie 时尝试获取访客 cookie
		// if c.guestFetching.CompareAndSwap(false, true) {
		if c.guestFetching.TryLock() {
			if c.fetchGuestCookie() != nil {
				return nil
			}
		}
	}

	parts := strings.Split(u.Host, ".")
	if len(parts) < 2 {
		return nil
	}
	if parts[len(parts)-2]+"."+parts[len(parts)-1] != URL_DOMAIN {
		return nil
	}

	return []*http.Cookie{c.cookie}
}

type client struct {
	*http.Client
	Headers http.Header
}

type Request struct {
	*http.Request
	querys *netUrl.Values

	wbi bool // wbi 签名
}

type Response struct {
	*http.Response
}

func fixUrl(u string) string {
	uu, err := netUrl.Parse(u)
	if err != nil {
		return u
	}
	if uu.Scheme == "" {
		uu.Scheme = "https"
	}
	return uu.String()
}

func NewHead(url string) *Request {
	req, err := http.NewRequest(http.MethodHead, fixUrl(url), nil)
	if err != nil {
		panic(err)
	}
	return &Request{
		Request: req,
		querys:  &netUrl.Values{},
	}
}

func NewHeadCtx(ctx context.Context, url string) *Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fixUrl(url), nil)
	if err != nil {
		panic(err)
	}
	return &Request{
		Request: req,
		querys:  &netUrl.Values{},
	}
}

func NewGet(url string) *Request {
	req, err := http.NewRequest(http.MethodGet, fixUrl(url), nil)
	if err != nil {
		panic(err)
	}
	return &Request{
		Request: req,
		querys:  &netUrl.Values{},
	}
}

func NewGetCtx(ctx context.Context, url string) *Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fixUrl(url), nil)
	if err != nil {
		panic(err)
	}
	return &Request{
		Request: req,
		querys:  &netUrl.Values{},
	}
}

func NewPost(url, contentType string, body io.Reader) *Request {
	req, err := http.NewRequest(http.MethodPost, fixUrl(url), body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", contentType)
	return &Request{
		Request: req,
		querys:  &netUrl.Values{},
	}
}

func NewPostCtx(ctx context.Context, url, contentType string, body io.Reader) *Request {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fixUrl(url), body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", contentType)
	return &Request{
		Request: req,
		querys:  &netUrl.Values{},
	}
}

func (req *Request) Std() *http.Request {
	if req.wbi {
		WbiSign(req.URL)
	}
	return req.Request
}

func (req *Request) WbiSign() *Request {
	req.wbi = true
	return req
}

func (req *Request) WithQuery(k, v string) *Request {
	req.querys.Add(k, v)
	return req
}

func (req *Request) WithQuerys(kv ...string) *Request {
	if len(kv) == 0 {
		return req
	}
	if len(kv)%2 != 0 {
		panic(fmt.Errorf("request.WithQuerys: odd kv: (%d)%v", len(kv), kv))
	}
	for i := 0; i < len(kv); i += 2 {
		req.querys.Add(kv[i], kv[i+1])
	}
	return req
}

func (req *Request) WithQueryList(k string, v ...string) *Request {
	if len(v) == 0 {
		panic(fmt.Errorf("request.WithQueryList: empty value for %s", k))
	}
	for _, i := range v {
		req.querys.Add(k, i)
	}
	return req
}

func (resp *Response) Std() *http.Response {
	return resp.Response
}

func (resp *Response) ToGjson() (j gjson.Result, err error) {
	b, err := io.ReadAll(resp.Response.Body)
	if err != nil {
		return
	}
	return gjson.Parse(toString(b)), nil
}

func (resp *Response) UnmarshalTo(v any) error {
	return json.NewDecoder(resp.Response.Body).Decode(v)
}
