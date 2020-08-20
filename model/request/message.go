package request

type EventSource struct {
	UserID     string `json:"userId"`
	SourceType string `json:"type"`
	GroupID    string `json:"groupId"`
	RoomID     string `json:"roomId"`
}

type EmojiMessage struct {
	Index     int    `json:"index"`
	Length    int    `json:"lenght"`
	ProductID string `json:"productId"`
	EmojiID   string `json:"emojiId"`
}

type MessageContentProvider struct {
	ContentType        string `json:"type"`
	OriginalContentURL string `json:"originalContentUrl"`
	PreviewImageURL    string `json:"previewImageUrl"`
}

type MessageUnsent struct {
	MessageID string `json:"messageId"`
}

type EventMessage struct {
	MessageType         string                 `json:"type"`
	ID                  string                 `json:"id"`
	Text                string                 `json:"text"`
	Emojis              []EmojiMessage         `json:"emojis"`
	PacakgeID           string                 `json:"packageId"`
	StickerID           string                 `json:"stickerId"`
	StickerResourceType string                 `json:"stickerResourceType"`
	Duration            int                    `json:"duration"`
	ContentProvider     MessageContentProvider `json:"contentProvider"`
	FileName            string                 `json:"filename"`
	FileSize            int                    `json:"fileSize"`
	Title               string                 `json:"title"`
	Address             string                 `json:"address"`
	Latitude            float64                `json:"latitude"`
	Longtitude          float64                `json:"longtitude"`
}

type Event struct {
	MessageType   string        `json:"type"`
	ReplyToken    string        `json:"replyToken"`
	Source        EventSource   `json:"source"`
	Timestamp     int           `json:"timestamp"`
	Mode          string        `json:"mode"`
	Message       EventMessage  `json:"message"`
	MessageUnsent MessageUnsent `json:"unsend"`
	// checkpoint member join
}

type WebhookEventOjbect struct {
	Events      []Event `json:"events"`
	Destination string  `json:"destination"`
}
