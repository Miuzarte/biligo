package biligo

import (
	"context"
	"iter"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Login(ctx context.Context) (qrcodeUrl string, it iter.Seq2[LoginCodeState, error], err error) {
	qg, err := FetchLoginQrcodeGenerate()
	if err != nil {
		return "", nil, err
	}
	if qg.Url == "" {
		return "", nil, wrapErr(ErrLoginEmptyQrcodeUrl, qg)
	}
	if qg.QrcodeKey == "" {
		return "", nil, wrapErr(ErrLoginEmptyQrcodeKey, qg)
	}

	return qg.Url, loginIter(ctx, qg.QrcodeKey), nil
}

func loginIter(ctx context.Context, qrcodeKey string) iter.Seq2[LoginCodeState, error] {
	return func(yield func(LoginCodeState, error) bool) {
		for {
			select {
			case <-ctx.Done():
				yield(-1, ctx.Err())
				return
			case <-time.After(time.Second):
			}

			qp, header, err := FetchLoginQrcodePoll(qrcodeKey)
			if err != nil {
				if yield(-1, err) {
					continue
				}
				return
			}
			code := LoginCodeState(qp.Code)

			switch code {
			case LOGIN_CODE_STATE_SUCCESS:
				setCookie := joinCookie(header["Set-Cookie"])
				refreshToken := qp.RefreshToken
				switch {
				case setCookie == "":
					err = wrapErr(ErrLoginEmptyCookie, qp)
				case refreshToken == "":
					err = wrapErr(ErrLoginEmptyRefreshToken, qp)
				default:
					cookie.set(setCookie)
					identity.Cookie = setCookie
					identity.RefreshToken = refreshToken
					keyI := strings.Index(setCookie, "DedeUserID=")
					if keyI == -1 {
						err = wrapErr(ErrLoginNoDedeUserID, setCookie)
						break
					}
					keyI += len("DedeUserID=")
					endI := strings.IndexByte(setCookie[keyI:], ';')
					var uidStr string
					if endI == -1 {
						uidStr = setCookie[keyI:]
					} else {
						uidStr = setCookie[keyI : keyI+endI]
					}
					identity.Uid, err = strconv.Atoi(uidStr)
					if err != nil {
						err = wrapErr(ErrLoginInvalidDedeUserID, uidStr)
						break
					}
					err = CookieUpdateBuvid34()
					if err != nil {
						break
					}
				}

			case LOGIN_CODE_STATE_EXPIRED:
				err = wrapErr(ErrLoginQrcodeExpired, nil)

			case LOGIN_CODE_STATE_SCANNED:
			case LOGIN_CODE_STATE_UNSCANNED:

			default:
				err = wrapErr(ErrLoginUnknownCode, qp)
			}

			if !yield(code, err) {
				return
			}

		}
	}
}

func joinCookie(cookies []string) (cookie string) {
	if len(cookies) == 0 {
		return ""
	}
	for _, c := range cookies {
		semicolon := strings.IndexByte(c, ';')
		if semicolon != -1 {
			c = c[:semicolon]
		}
		if cookie != "" {
			cookie += "; "
		}
		cookie += c
	}
	return cookie
}

func CookieUpdateBuvid34() error {
	if identity.Cookie == "" {
		return wrapErr(ErrLoginEmptyCookie, nil)
	}
	b34, err := FetchBuvid34()
	if err != nil {
		return err
	}
	c := identity.Cookie
	b3Index := strings.Index(c, "buvid3=")
	if b3Index != -1 {
		endIndex := strings.IndexByte(c[b3Index:], ';')
		if endIndex == -1 {
			c = c[:b3Index]
		} else {
			c = c[:b3Index+endIndex]
		}
	}
	b4Index := strings.Index(c, "buvid4=")
	if b4Index != -1 {
		endIndex := strings.IndexByte(c[b4Index:], ';')
		if endIndex == -1 {
			c = c[:b4Index]
		} else {
			c = c[:b4Index+endIndex]
		}
	}
	if c != "" && !strings.HasSuffix(c, ";") {
		c += "; "
	}
	c += "buvid3=" + b34.B_3 + "; buvid4=" + b34.B_4
	cookie.set(c)
	identity.Cookie = c
	return nil
}

func getSpaceUid() (int, error) {
	userSpace, err := fetchLocation("https://space.bilibili.com")
	if err != nil {
		return 0, err
	}
	if strings.Contains(userSpace, "passport.bilibili.com") {
		return 0, wrapErr(ErrLoginFailed, nil)
	}
	u, err := url.Parse(userSpace)
	if err != nil {
		panic(err)
	}
	if u.Host != "space.bilibili.com" {
		return 0, wrapErr(ErrLoginUnexpectedHost, u.Host)
	}
	id := strings.TrimPrefix(u.Path, "/")
	if id == "" {
		return 0, wrapErr(ErrLoginEmptyUserId, userSpace)
	}
	return strconv.Atoi(id)
}
