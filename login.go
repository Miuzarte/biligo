package biligo

import (
	"iter"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Login() (qrcodeUrl string, it iter.Seq2[LoginCodeState, error], err error) {
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

	return qg.Url, loginIter(qg.QrcodeKey), nil
}

func loginIter(qrcodeKey string) iter.Seq2[LoginCodeState, error] {
	return func(yield func(LoginCodeState, error) bool) {
		for {
			time.Sleep(time.Second)
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
				setCookie := header.Get("Set-Cookie")
				refreshToken := qp.RefreshToken
				switch {
				case setCookie == "":
					err = wrapErr(ErrLoginEmptyCookie, qp)
				case refreshToken == "":
					err = wrapErr(ErrLoginEmptyRefreshToken, qp)
				default:
					cookie.store(setCookie)
					identity.Cookie = setCookie
					identity.RefreshToken = refreshToken
					identity.Uid, err = getUid()
				}

			case LOGIN_CODE_STATE_EXPIRED:
				err = wrapErr(ErrLoginQrcodeExpired, nil)

			case LOGIN_CODE_STATE_SCANED:
			case LOGIN_CODE_STATE_UNSCANED:

			default:
				err = wrapErr(ErrLoginUnknownCode, qp)
			}

			if !yield(code, err) {
				return
			}

		}
	}
}

func getUid() (int, error) {
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
