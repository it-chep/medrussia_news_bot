package dto

type MediaType int8

const (
	// Unknown Текст или ничего
	Unknown MediaType = iota
	// Photo фотография
	Photo
	// Video Видео
	Video
	// VideoNote Кружок
	VideoNote
	// Voice Голосовое сообщение
	Voice
	// Document документ
	Document
	// Audio mp3
	Audio
)
