package biligo

import (
	"context"
	"time"
)

// RegisterVideoConclusion 注册视频总结轮询回调,
// 视频无总结时返回 [ErrPollNoSummary] 错误
func RegisterVideoConclusion(ctx context.Context, id, cid string, f func(VideoConclusion, error)) (err error) {
	if f == nil {
		panic("RegisterVideoConclusion: f == nil")
	}
	if cid == "" || cid == "0" {
		cid, err = FetchVideoCidTryCache(id)
		if err != nil {
			return
		}
	}

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

		req := Chain{Req: ReqVideoConclusion(id, cid)}
		var vc VideoConclusion
		var err error

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			err = req.Do()
			if err != nil {
				f(vc, err)
				return
			}
			vc, err = req.ToVideoConclusion()
			if err != nil {
				f(vc, err)
				return
			}

			switch {
			case vc.Code == 0: // finished
				if vc.ModelResult.Summary == "" {
					f(vc, wrapErr(ErrPollNoSummary, nil))
					return
				}
				f(vc, nil)
				return

			case vc.Code == 1 && vc.Stid == "0": // in queue

			case vc.Code == 1 && vc.Stid == "": // no voice recognized
				f(vc, wrapErr(ErrPollNoVoiceRecognized, nil))
				return
			default: // unknown code
				f(vc, wrapErr(ErrPollUnknownCode, req.Body))
				return
			}
		}
	}()

	return nil
}
