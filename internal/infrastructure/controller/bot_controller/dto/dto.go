package dto

type Chat struct {
	ID int64 `json:"id"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

type Animation struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

type Document struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

type Video struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

type VideoNote struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
}

type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	UserID      int64  `json:"user_id,omitempty"`
}

type Sticker struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	IsAnimated   bool   `json:"is_animated,omitempty"`
	Emoji        string `json:"emoji,omitempty"`
	SetName      string `json:"set_name,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
}

type InlineKeyboardButton struct {
	Text                         string  `json:"text"`
	URL                          *string `json:"url,omitempty"`
	CallbackData                 *string `json:"callback_data,omitempty"`
	SwitchInlineQuery            *string `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat *string `json:"switch_inline_query_current_chat,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type MessageDTO struct {
	MessageID         int64                 `json:"message_id"`
	Date              int                   `json:"date"`
	Chat              *Chat                 `json:"chat"`
	ForwardDate       int                   `json:"forward_date,omitempty"`
	Text              string                `json:"text,omitempty"`
	Animation         *Animation            `json:"animation,omitempty"`
	Audio             *Audio                `json:"audio,omitempty"`
	Document          *Document             `json:"document,omitempty"`
	Photo             []PhotoSize           `json:"photo,omitempty"`
	Sticker           *Sticker              `json:"sticker,omitempty"`
	Video             *Video                `json:"video,omitempty"`
	VideoNote         *VideoNote            `json:"video_note,omitempty"`
	Voice             *Voice                `json:"voice,omitempty"`
	Contact           *Contact              `json:"contact,omitempty"`
	MigrateToChatID   int64                 `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID int64                 `json:"migrate_from_chat_id,omitempty"`
	ReplyMarkup       *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type TgUserDTO struct {
	TgID                  int64  `json:"id"`
	IsBot                 bool   `json:"is_bot,omitempty"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name,omitempty"`
	UserName              string `json:"username,omitempty"`
	LanguageCode          string `json:"language_code,omitempty"`
	SupportsInlineQueries bool   `json:"supports_inline_queries,omitempty"`
}
