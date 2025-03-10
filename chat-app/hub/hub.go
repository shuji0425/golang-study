package hub

// 全てのクライアントを管理
// メッセージのブロードキャストを行う
type Hub struct {
	clients    map[*Client]bool // 登録されているクライアント一覧
	broadcast  chan []byte      // ブロードキャスト用のチャネル
	register   chan *Client     // クライアント登録用チャネル
	unregister chan *Client     // クライアント解除用のチェネル
}

// Hubを初期化して返す
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// メインループ
// 登録・解除・メッセージのブロードキャストを処理
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// 全クライアントにメッセージを送信
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// 送信できない場合は接続を切断
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
