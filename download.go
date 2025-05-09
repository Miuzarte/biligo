package biligo

import (
	"context"
	"io"
	"sync"
)

type VideoDownload struct {
	ctx context.Context
	id  string
	cid string
	qn  VideoQuality
	job *dlJob
}

func NewDownloadVideoMp4(ctx context.Context, id, cid string, qn VideoQuality) *VideoDownload {
	return &VideoDownload{ctx, id, cid, qn, nil}
}

// Init 返回文件大小以便提前申请内存
func (vd *VideoDownload) Init() (size int64, err error) {
	vp, err := FetchVideoPlayurl(vd.id, vd.cid, VIDEO_FNVAL_MP4, vd.qn)
	if err != nil {
		return 0, err
	}

	if len(vp.Durl) == 0 {
		return 0, wrapErr(ErrDownloadEmptyDUrl, vp)
	}
	durl := vp.Durl[0]

	var dlUrl string
	if len(durl.BackupUrl) != 0 { // 优先使用 cdn 备链
		dlUrl = durl.BackupUrl[0]
	} else {
		dlUrl = durl.Url // mcdn
	}
	if dlUrl == "" {
		return 0, wrapErr(ErrDownloadNoUrlFound, vp)
	}

	vd.job = newDljob(vd.ctx, dlUrl, nil)
	err = vd.job.init()
	if err != nil {
		return 0, err
	}
	return vd.job.size, nil
}

func (vd *VideoDownload) Start(w io.Writer) (written int64, err error) {
	vd.job.w = w
	return vd.job.start()
}

// func DownloadVideoDash(ctx context.Context, id, cid string, 优先清晰度, (优先)编码(不存在对应编码时返回错误/更换编码并返回下载的编码))

var (
	threadNum       = 8
	blockSize int64 = 1 << 24 // 16 MiB
)

type dlJob struct {
	url    string
	ctx    context.Context
	cancel context.CancelFunc
	w      io.Writer
	size   int64
	blocks dlBlocks
}

type dlBlocks []*dlBlock

type dlBlock struct {
	start  int64
	end    int64
	err    chan error
	body   io.ReadCloser
	copied chan struct{}
}

func newDljob(ctx context.Context, url string, w io.Writer) *dlJob {
	ctx, cancel := context.WithCancel(ctx)
	return &dlJob{
		url:    url,
		ctx:    ctx,
		cancel: cancel,
		w:      w,
	}
}

func (j *dlJob) init() error {
	if j.size > 0 { // inited
		return nil
	}

	// 获取大小
	resp, err := httpClient.Do(NewHead(j.url))
	if err != nil {
		return err
	}
	defer resp.Body.Close() // HEAD 请求也需要关闭
	j.size = resp.ContentLength
	if j.size <= 0 {
		return wrapErr(ErrDownloadInvaliedFileSize, resp.Header)
	}

	// 初始化分块
	numBlocks := max((j.size+blockSize-1)/blockSize, 1) // >=1
	j.blocks = make(dlBlocks, numBlocks)
	for i := range numBlocks {
		start := blockSize * i     // 左闭
		end := blockSize*(i+1) - 1 // 右闭
		if i == numBlocks-1 {
			end = j.size - 1
		}
		j.blocks[i] = &dlBlock{
			start:  start,
			end:    end,
			err:    make(chan error, 1),
			copied: make(chan struct{}),
		}
	}

	return nil
}

func (j *dlJob) start() (written int64, err error) {
	go func() {
		limiter := newLimiter()
		defer limiter.close()

		for _, block := range j.blocks {
			select {
			case <-j.ctx.Done():
				return
			case limiter.acquireChannel() <- struct{}{}:
				// fmt.Printf("%d downloading...\n", i)
			}

			go func(block *dlBlock) {
				defer limiter.release()
				req := NewGetCtx(j.ctx, j.url)
				// req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", block.start, block.end))
				req.Header.Set("Range", "bytes="+itoa(block.start)+"-"+itoa(block.end))
				resp, err := httpClient.Do(req)
				if err != nil {
					block.err <- err
					return
				}
				block.body = resp.Body
				block.err <- nil
				select {
				case <-j.ctx.Done():
					return
				case <-block.copied:
					// 当前块下载完成后再释放锁
				}
			}(block)
		}
	}()

	defer j.cancel()
	for _, block := range j.blocks {
		select {
		case <-j.ctx.Done():
			return
		case err = <-block.err:
			// fmt.Printf("%d: %d - %d\n", i, block.start, block.end)
		}
		if err != nil {
			return
		}

		var n int64
		n, err = io.Copy(j.w, block.body)
		block.body.Close()
		written += n
		if err != nil {
			return
		}

		select {
		case <-j.ctx.Done():
			return
		case block.copied <- struct{}{}:
		}
	}
	return written, nil
}

type limiter struct {
	sem  chan struct{}
	once sync.Once
}

func newLimiter() *limiter {
	return &limiter{
		sem: make(chan struct{}, threadNum),
	}
}

func (l *limiter) acquireDirect() {
	l.sem <- struct{}{}
}

func (l *limiter) acquireChannel() chan<- struct{} {
	return l.sem
}

func (l *limiter) release() {
	<-l.sem
}

func (l *limiter) close() {
	l.once.Do(func() {
		close(l.sem)
	})
}
