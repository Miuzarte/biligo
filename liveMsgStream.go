package biligo

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gorilla/websocket"
)

func SetWsDialerProxy(pf func(*http.Request) (*url.URL, error)) {
	wsDialer.Proxy = pf
}

var wsDialer = &websocket.Dialer{
	// Proxy: http.ProxyFromEnvironment,
	Proxy:            nil,
	HandshakeTimeout: 45 * time.Second,
	Jar:              cookie,
}

type LiveMsgStream struct {
	*websocket.Conn
	connMu sync.Mutex

	roomId   int
	sequence uint32

	heartBeatRunning atomic.Bool
	heartBeatStop    chan struct{}

	// cmdStatistics map[string]int
	// protoVerStatistics map[uint16]int
}

func NewLiveMsgStream(roomId int) *LiveMsgStream {
	lms := &LiveMsgStream{
		heartBeatStop: make(chan struct{}, 1),
		roomId:        roomId,
		// cmdStatistics: make(map[string]int),
		// protoVerStatistics: make(map[uint16]int),
	}
	return lms
}

func (s *LiveMsgStream) RunIter() iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		// retry until break by outside
		for {
			err := s.connect()
			if err != nil {
				if yield("", err) {
					continue
				}
				return
			}
			break
		}
		// connected
		defer s.Stop()
		s.heartBeatRunning.Store(true)
		go s.hearteatLoop()

		for {
			msgType, data, err := s.Conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseAbnormalClosure) {
					yield("", err)
					return
				}
				if yield("", err) {
					continue
				}
				return
			}
			if msgType != websocket.BinaryMessage {
				if yield("", wrapErr(ErrLmsPacketNotBinary, msgType)) {
					continue
				}
				return
			}

			pkt, err := s.parsePacket(data)
			if err != nil {
				if yield("", err) {
					continue
				}
				return
			}

			for pkt, err := range s.ParseIter(pkt) {
				if err != nil {
					if yield("", err) {
						continue
					}
					return
				}

				// s.cmdStatistics[gjson.Get(toString(pkt.Body), "cmd").String()]++
				if !yield(toString(pkt.Body), nil) {
					return
				}
			}
		}
	}
}

// "data"
type LiveDanmuInfo struct {
	Group            string  `json:"group" mapstructure:"group"`
	BusinessId       int     `json:"business_id" mapstructure:"business_id"`
	RefreshRowFactor float64 `json:"refresh_row_factor" mapstructure:"refresh_row_factor"`
	RefreshRate      int     `json:"refresh_rate" mapstructure:"refresh_rate"`
	MaxDelay         int     `json:"max_delay" mapstructure:"max_delay"`
	Token            string  `json:"token" mapstructure:"token"`
	HostList         []struct {
		Host    string `json:"host" mapstructure:"host"`
		Port    int    `json:"port" mapstructure:"port"`
		WsPort  int    `json:"ws_port" mapstructure:"ws_port"`
		WssPort int    `json:"wss_port" mapstructure:"wss_port"`
	} `json:"host_list" mapstructure:"host_list"`
}

// connect 发起连接并发送认证包
func (s *LiveMsgStream) connect() error {
	ldi, err := FetchLiveDanmuInfo(itoa(s.roomId))
	if err != nil {
		return err
	}
	if ldi.Token == "" {
		return wrapErr(ErrLmsFailedToGetToken, ldi)
	}
	hosts := make([]string, 0, len(ldi.HostList))
	for _, h := range ldi.HostList {
		hosts = append(hosts, "wss://"+h.Host+"/sub")
	}

	reqHeader := http.Header{}
	for k, v := range DefaultHeaders {
		reqHeader.Add(k, v)
	}
	var conn *websocket.Conn
	for _, h := range hosts {
		conn, _, err = wsDialer.Dial(h, reqHeader)
		if err != nil {
			continue
		}
		break
	}
	if err != nil {
		return err
	}

	// 认证包
	pkt, err := s.newAuthPacket(ldi.Token)
	if err != nil {
		conn.Close()
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, pkt)
	if err != nil {
		conn.Close()
		return err
	}

	s.Conn = conn
	return nil
}

func (s *LiveMsgStream) Write(data []byte) error {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	if s.Conn == nil {
		return wrapErr(ErrLmsNilConn, nil)
	}
	return s.Conn.WriteMessage(websocket.BinaryMessage, data)
}

func (s *LiveMsgStream) Stop() *LiveMsgStream {
	if s.Conn == nil {
		return s
	}
	go s.tryStopHeartbeat()
	s.Write(websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return s
}

// func (s *LiveMsgStream) Statistics() (map[string]int, map[uint16]int) {
// 	return s.cmdStatistics, s.protoVerStatistics
// }

func (s *LiveMsgStream) tryStopHeartbeat() {
	if !s.heartBeatRunning.Load() {
		return
	}
	select {
	case s.heartBeatStop <- struct{}{}:
	case <-time.After(time.Second * 5):
	}
}

// hearteatLoop 每 30s 发送心跳包
func (s *LiveMsgStream) hearteatLoop() {
	defer s.heartBeatRunning.Store(false)

	var pkt liveMsgPacket
	var err error
mainLoop:
	for {
		select {
		case <-time.After(time.Second * 30):
			for range 4 { // 重试最多 4 次
				pkt = s.newPacket(PROTOCOL_PLAIN, OPERATION_HEARTBEAT, nil)
				err = s.Write(pkt.build())
				if err == nil {
					continue mainLoop
				}

				select {
				case <-time.After(time.Second * 5):
				case <-s.heartBeatStop:
					return
				}
			}
			s.Stop()
			return

		case <-s.heartBeatStop:
			return
		}
	}
}

const (
	PROTOCOL_PLAIN = iota
	PROTOCOL_POPULARITY
	PROTOCOL_ZLIB
	PROTOCOL_BROTLI
)

const (
	_ = iota
	_
	OPERATION_HEARTBEAT
	OPERATION_HEARTBEAT_RESP
	_
	OPERATION_NORMAL
	_
	OPERATION_AUTHENTICATION
	OPERATION_AUTHENTICATION_RESP
)

type liveMsgPacket struct {
	PacketLength    uint32 // 4
	HeaderLength    uint16 // 2
	ProtocolVersion uint16 // 2
	Operation       uint32 // 4
	Sequence        uint32 // 4
	// header 16
	Body []byte
}

func (p *liveMsgPacket) build() []byte {
	rawBuf := make([]byte, 16)
	binary.BigEndian.PutUint32(rawBuf[12:], p.Sequence)       // Sequence        uint32 // 4
	binary.BigEndian.PutUint32(rawBuf[8:], p.Operation)       // Operation       uint32 // 4
	binary.BigEndian.PutUint16(rawBuf[6:], p.ProtocolVersion) // ProtocolVersion uint16 // 2
	binary.BigEndian.PutUint16(rawBuf[4:], 0x0010)            // HeaderLength    uint16 // 2
	rawBuf = slices.Concat(rawBuf, p.Body)
	binary.BigEndian.PutUint32(rawBuf[0:], uint32(len(rawBuf))) // PacketLength    uint32 // 4
	return rawBuf
}

func (s *LiveMsgStream) newPacket(protocolVersion uint16, operation uint32, body []byte) liveMsgPacket {
	s.sequence++
	return liveMsgPacket{
		ProtocolVersion: protocolVersion,
		Operation:       operation,
		Sequence:        s.sequence,
		Body:            body,
	}
}

func (s *LiveMsgStream) parsePacket(data []byte) (pkt liveMsgPacket, err error) {
	pakLen := binary.BigEndian.Uint32(data[0:4])
	if int(pakLen) != len(data) {
		return pkt, wrapErr(ErrLmsInvalidPacket, fmt.Sprintf("int(pakLen) %d != len(data) %d", pakLen, len(data)))
	}
	headLen := binary.BigEndian.Uint16(data[4:6])
	if int(headLen) != 16 {
		return pkt, wrapErr(ErrLmsInvalidPacket, fmt.Sprintf("int(headLen) %d != 16", headLen))
	}
	protocol := binary.BigEndian.Uint16(data[6:8])
	operation := binary.BigEndian.Uint32(data[8:12])
	body := data[16:pakLen]
	pkt = s.newPacket(protocol, operation, body)
	return
}

type liveAuthenticate struct {
	Uid      int    `json:"uid"`
	RoomId   int    `json:"roomid"`
	ProtoVer int    `json:"protover"`
	Platform string `json:"platform"`
	Type     int    `json:"type"`
	Key      string `json:"key"`
}

func (s *LiveMsgStream) newAuthPacket(key string) ([]byte, error) {
	j, err := json.Marshal(
		liveAuthenticate{
			Uid:      identity.Uid, // uid 为 0 游客登录
			RoomId:   s.roomId,
			ProtoVer: 3,
			Platform: "web",
			Type:     2,
			Key:      key,
		},
	)
	if err != nil {
		panic(err)
	}

	pkt := s.newPacket(PROTOCOL_PLAIN, OPERATION_AUTHENTICATION, j)
	return pkt.build(), nil
}

func (s *LiveMsgStream) slicePacketsIter(data []byte) iter.Seq2[liveMsgPacket, error] {
	return func(yield func(liveMsgPacket, error) bool) {
		total := len(data)
		pktPos := 0
		pktLen := 0
		for pktPos < total {
			pktLen = int(binary.BigEndian.Uint32(data[pktPos : pktPos+4]))
			if pktLen > total {
				yield(liveMsgPacket{}, wrapErr(ErrLmsInvalidPacket, fmt.Sprintf("pktLen %d > total %d", pktLen, total)))
				return
			}

			if !yield(s.parsePacket(data[pktPos : pktPos+pktLen])) {
				return
			}
			pktPos += pktLen
		}
	}
}

func (s *LiveMsgStream) ParseIter(p liveMsgPacket) iter.Seq2[liveMsgPacket, error] {
	return func(yield func(liveMsgPacket, error) bool) {
		// s.protoVerStatistics[p.ProtocolVersion]++
		switch p.ProtocolVersion {
		case PROTOCOL_PLAIN, PROTOCOL_POPULARITY:
			yield(p, nil)

		case PROTOCOL_ZLIB:
			zr, err := zlib.NewReader(bytes.NewReader(p.Body))
			if err != nil {
				yield(p, err)
				return
			}
			buf, err := io.ReadAll(zr)
			if err != nil {
				yield(p, err)
				return
			}
			for pkt, err := range s.slicePacketsIter(buf) {
				if !yield(pkt, err) {
					return
				}
			}

		case PROTOCOL_BROTLI:
			br := brotli.NewReader(bytes.NewReader(p.Body))
			buf, err := io.ReadAll(br)
			if err != nil {
				yield(p, err)
				return
			}
			for pkt, err := range s.slicePacketsIter(buf) {
				if !yield(pkt, err) {
					return
				}
			}

		default:
			yield(p, wrapErr(ErrLmsUnknownProtocol, p.ProtocolVersion))
			return
		}
	}
}

var LiveMsgCmdMap = map[string]struct{}{
	"ANCHOR_LOT_AWARD":                     {},
	"ANCHOR_LOT_CHECKSTATUS":               {},
	"ANCHOR_LOT_END":                       {},
	"ANCHOR_LOT_NOTICE":                    {},
	"ANCHOR_LOT_START":                     {},
	"AREA_RANK_CHANGED":                    {},
	"CHANGE_ROOM_INFO":                     {},
	"COMBO_SEND":                           {},
	"COMBO_END":                            {},
	"COMMON_ANIMATION":                     {},
	"COMMON_NOTICE_DANMAKU":                {},
	"CUT_OFF":                              {}, // 切断 .Get("msg").String()
	"DANMU_AGGREGATION":                    {},
	"DANMU_MSG":                            {},
	"DM_INTERACTION":                       {}, // 弹幕
	"ENTRY_EFFECT":                         {},
	"ENTRY_EFFECT_MUST_RECEIVE":            {},
	"FULL_SCREEN_SPECIAL_EFFECT":           {},
	"GIFT_STAR_PROCESS":                    {},
	"GOTO_BUY_FLOW":                        {},
	"GUARD_BUY":                            {},
	"GUARD_HONOR_THOUSAND":                 {},
	"HOT_RANK_CHANGED":                     {},
	"HOT_RANK_CHANGED_V2":                  {},
	"HOT_RANK_SETTLEMENT":                  {},
	"HOT_RANK_SETTLEMENT_V2":               {},
	"INTERACT_WORD":                        {},
	"LIKE_GUIDE_USER":                      {},
	"LIKE_INFO_V3_CLICK":                   {},
	"LIKE_INFO_V3_UPDATE":                  {},
	"LIVE":                                 {}, // 直播开始
	"LIVE_PANEL_CHANGE":                    {},
	"LIVE_ANI_RES_UPDATE":                  {},
	"LOG_IN_NOTICE":                        {},
	"MESSAGEBOX_USER_GAIN_MEDAL":           {},
	"NOTICE_MSG":                           {},
	"ONLINE_RANK_COUNT":                    {},
	"ONLINE_RANK_V2":                       {},
	"ONLINE_RANK_TOP3":                     {},
	"PK_BATTLE_END":                        {},
	"PK_BATTLE_FINAL_PROCESS":              {},
	"PK_BATTLE_PRE":                        {},
	"PK_BATTLE_PRE_NEW":                    {},
	"PK_BATTLE_START":                      {},
	"PK_BATTLE_START_NEW":                  {},
	"PK_BATTLE_PROCESS":                    {},
	"PK_BATTLE_PROCESS_NEW":                {},
	"PK_BATTLE_SETTLE":                     {},
	"PK_BATTLE_SETTLE_USER":                {},
	"PK_BATTLE_SETTLE_V2":                  {},
	"PLAY_TOGETHER":                        {},
	"POPULAR_RANK_CHANGED":                 {},
	"POPULARITY_RED_POCKET_NEW":            {},
	"POPULARITY_RED_POCKET_START":          {},
	"POPULARITY_RED_POCKET_V2_NEW":         {},
	"POPULARITY_RED_POCKET_V2_START":       {},
	"POPULARITY_RED_POCKET_V2_WINNER_LIST": {},
	"POPULARITY_RED_POCKET_WINNER_LIST":    {},
	"PREPARING":                            {}, // 直播准备中 (结束)
	"RANK_CHANGED":                         {},
	"RANK_REM":                             {},
	"RECALL_DANMU_MSG":                     {},
	"RECOMMEND_CARD":                       {},
	"REENTER_LIVE_ROOM":                    {},
	"REVENUE_RANK_CHANGED":                 {},
	"RING_STATUS_CHANGE":                   {},
	"RING_STATUS_CHANGE_V2":                {},
	"room_admin_entrance":                  {},
	"ROOM_ADMIN_REVOKE":                    {},
	"ROOM_ADMINS":                          {},
	"ROOM_CHANGE":                          {}, // 房间信息变更
	"ROOM_CONTENT_AUDIT_REPORT":            {},
	"ROOM_REAL_TIME_MESSAGE_UPDATE":        {},
	"ROOM_SKIN_MSG":                        {},
	"ROOM_SILENT_OFF":                      {},
	"ROOM_SILENT_ON":                       {},
	"SEND_GIFT":                            {},
	"SHOPPING_CART_SHOW":                   {},
	"SPECIAL_GIFT":                         {},
	"SPREAD_SHOW_FEET_V2":                  {},
	"STOP_LIVE_ROOM_LIST":                  {},
	"SUPER_CHAT_ENTRANCE":                  {},
	"SUPER_CHAT_MESSAGE":                   {},
	"SUPER_CHAT_MESSAGE_DELETE":            {},
	"SUPER_CHAT_MESSAGE_JPN":               {},
	"SYS_MSG":                              {},
	"TRADING_SCORE":                        {},
	"USER_TOAST_MSG":                       {},
	"USER_TOAST_MSG_V2":                    {},
	"VIDEO_CONNECTION_JOIN_END":            {},
	"VIDEO_CONNECTION_JOIN_START":          {},
	"VIDEO_CONNECTION_MSG":                 {},
	"WARNING":                              {}, // 警告 .Get("msg").String()
	"WATCHED_CHANGE":                       {},
	"WIDGET_BANNER":                        {},
	"WIDGET_GIFT_STAR_PROCESS":             {},
	"WIDGET_WISH_LIST":                     {},
}
