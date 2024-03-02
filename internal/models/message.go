package models

type Message struct {
	TopicName string `json:"topic_name"`
	Message   struct {
		EventType      string `json:"event_type"`
		UserId         int    `json:"user_id"`
		Action         string `json:"action"`
		Timestamp      string `json:"timestamp"`
		AdditionalData struct {
			ButtonId string `json:"button_id"`
			Page     string `json:"page"`
		} `json:"additional_data"`
	} `json:"message"`
}
