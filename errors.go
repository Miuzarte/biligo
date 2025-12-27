package biligo

import (
	"errors"
	"fmt"
	"io"
)

func wrapErr(err error, detail any) error {
	if err == nil {
		return nil
	}
	return &Error{raw: err, detail: detail}
}

func UnwrapErr(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{raw: err}
}

// Error 附带了更详细的信息
//
// biliErr := biligo.UnwarpErr(err)
//
// isXxx := biliErr.Is(ErrXXX)
type Error struct {
	raw    error
	detail any
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.detail != nil {
		return fmt.Sprintf("%s: %v", e.raw.Error(), e.detail)
	}
	return e.raw.Error()
}

// implement [Templatable] for convenient string output
func (e *Error) DoTemplate() string {
	return e.Error()
}

// implement [Templatable] for convenient string output
func (e *Error) DoTemplateTo(w io.Writer) error {
	w.Write([]byte(e.Error()))
	return e
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.raw
}

func (e *Error) Is(target error) bool {
	if e == nil {
		return target == nil
	}
	return e.raw == target
}

var ErrInvalidId = errors.New("invalid id")

var (
	ErrChainRespCodeNotZero = errors.New("response code not zero")
	ErrChainNilResp         = errors.New("response is nil")
	ErrChainNilBody         = errors.New("body is nil")
	ErrChainPathNotExists   = errors.New("path not exists")
)

var ErrClientNilRequest = errors.New("nil request provided")

var (
	ErrDownloadEmptyDUrl        = errors.New("empty durl")
	ErrDownloadNoUrlFound       = errors.New("no url found")
	ErrDownloadInvaliedFileSize = errors.New("invalid file size")
)

var (
	ErrFetchNoLocation      = errors.New("no location")
	ErrFetchRespCodeNotZero = errors.New("response code not zero")
	ErrFetchPathNotExists   = errors.New("path not exists")
	ErrFetchCidNotExists    = errors.New("cid not exists")
)

var ErrFetchDynamicSpaceItemsNotEqual = errors.New("dynamic space items length not equal")

var (
	ErrLmsPacketNotBinary  = errors.New("packet not binary")
	ErrLmsFailedToGetToken = errors.New("failed to get token")
	ErrLmsNilConn          = errors.New("conn is nil")
	ErrLmsInvalidPacket    = errors.New("invalid packet")
	ErrLmsUnknownProtocol  = errors.New("unknown protocol")
)

var (
	ErrLoginEmptyQrcodeUrl = errors.New("empty qrcode url")
	ErrLoginEmptyQrcodeKey = errors.New("empty qrcode key")

	ErrLoginEmptyCookie       = errors.New("empty cookie")
	ErrLoginEmptyRefreshToken = errors.New("empty refresh token")
	ErrLoginNoDedeUserID      = errors.New("no DedeUserID in cookie")
	ErrLoginInvalidDedeUserID = errors.New("invalid DedeUserID in cookie")
	ErrLoginQrcodeExpired     = errors.New("qrcode expired")
	ErrLoginUnknownCode       = errors.New("unknown code")

	ErrLoginFailed         = errors.New("failed to login")
	ErrLoginUnexpectedHost = errors.New("unexpected host")
	ErrLoginEmptyUserId    = errors.New("empty user id")
)

var (
	ErrPollNoSummary         = errors.New("no summary") // 本视频暂无AI总结内容
	ErrPollNoVoiceRecognized = errors.New("no voice recognized")
	ErrPollUnknownCode       = errors.New("unknown code")
)

var (
	ErrParseMediaNoSsid = errors.New("no ssid in result")
	ErrParseLiveNoUid   = errors.New("no uid in result")
)

var (
	ErrSearchUnknownType = errors.New("unknown search type")
	ErrSearchEmptyType   = errors.New("empty search type")
)

var ErrWbiEmptyUrls = errors.New("empty image or sub url")
