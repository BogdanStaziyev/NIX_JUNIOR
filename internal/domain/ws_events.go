package domain

const (
	ActionRegister              string = "register"
	ActionLeaveChat             string = "leave-chat"
	ActionJoinChat              string = "join-chat"
	ActionSandMessageToChat     string = "send-message-chat"
	ActionSandMessageToAllUsers string = "send-message-all"
	ActionSendPrivate           string = "send-private"
)

type Base struct {
	Action string `json:"action"`
}

type Register struct {
	Base        // actionRegister
	Name string `json:"name"`
}

type JoinChatroom struct {
	Base         // actionJoinChat
	ChatID int64 `json:"chat_id"`
}

type LeaveChat struct {
	Base         // actionLeaveChat
	ChatID int64 `json:"chat_id"`
}

type SendMessageToChat struct {
	Base           // ActionSandMessageToChat
	Message string `json:"message"`
}

type SendMessageToAll struct {
	Base           // ActionSandMessageToAllUsers
	Message string `json:"message"`
}

type SendMessageToOne struct {
	Base           // ActionSendPrivate
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}
