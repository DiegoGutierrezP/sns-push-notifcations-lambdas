package models

import (
	"time"
)

type TopicModel struct {
	ID          string    `dynamodbav:"pk"`
	TopicArn    string    `dynamodbav:"topicArn"`
	Name        string    `dynamodbav:"name"`
	Description string    `dynamodbav:"description"`
	Status      int16     `dynamodbav:"status"`
	CreatedAt   time.Time `dynamodbav:"createdAt"`
	UpdatedAt   time.Time `dynamodbav:"updatedAt"`
}

func (TopicModel) TableName() string {
	return "topics"
}
