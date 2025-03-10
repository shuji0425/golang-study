package hub

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 定数の定義
const (
	writeWait      = 10 * time.Second    // 書き込みのタイムアウト
	pongWait       = 60 * time.Second    // Pongの待機時間
	pingPeriod     = (pongWait * 9) / 10 // Pingの送信間隔(PongWaitの90%)
	maxMessageSize = 512                 // 最大メッセージサイズ
)

// WebSocketアップグレーダーの設定
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 同一オリジンポリシーのチェックを簡略化（本番環境では注意が必要）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Clientは各WebSocket接続と送信用チェネルを保持
type Client struct {
	hub  *Hub
	conn *websocket.Conn // WebSocket接続
	send chan []byte     // メッセージ送信用チャネル
}

// クライアントからのメッセージ受信処理
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		// メッセージ一覧
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("readPumpエラー: %v", err)
			}
			break
		}
		// 受信したメッセージをHubにブロードキャスト
		c.hub.broadcast <- message
	}
}

// クライアントへのメッセージ送信処理
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hubからチャネルを閉じられた場合は、WebSocket接続を終了
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// 送信メッセージをWebSocketに書き込む
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			// 定期的にPingメッセージを送信して接続が生存しているか確認
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// HTTPリクエストをWebSocket接続にアップグレードし、新しいClientを生成・登録する
func ServerWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("serveWsアップグレードエラー:", err)
		return
	}
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	// 読み込みと書き込みをそれぞれゴルーチンで実行
	go client.writePump()
	go client.readPump()
}
