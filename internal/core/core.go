package core

type User struct {
	UserID string
	Name   string
}

const (
	ErrorAnswer               = "Не удалось выполнить запрос"
	NewFlatAnswerCallback     = "Введите номер квартиры. Пример ввода: 1А"
	NotFoundCommand           = "Команда не задана или введена неправильно"
	TaskCompletedSuccessfully = "Задача успешно выполнена"
	RepeatingMeaning          = "Это значение уже существует"
	NewCountAnswerCallback    = "Введите номер квартиры и показание счетчика. Пример ввода: 1A 123456"
	ErrorInputData            = "Неверный тип ввода данных"
)

type Message struct {
	Text            string
	UserID          int64
	UserName        string
	UserDisplayName string
	IsCallback      bool
	CallbackMsgID   string
	CallbackData    string
	MessageChatID   int64
	CallbackChatID  int64
}
