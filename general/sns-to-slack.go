package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type SNSMessage struct {
	Message struct {
		NewStateReason  string `json:"NewStateReason"`
		StateChangeTime string `json:"StateChangeTime"`
		Trigger         struct {
			MetricName string `json:"MetricName"`
			Dimensions []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"Dimensions"`
		} `json:"Trigger"`
	} `json:"Message"`
}

type DynamoDBMetric struct {
	ComparisonOperator string `json:"ComparisonOperator"`
	Dimensions         []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"Dimensions"`
	EvaluateLowSampleCountPercentile string      `json:"EvaluateLowSampleCountPercentile"`
	EvaluationPeriods                int         `json:"EvaluationPeriods"`
	MetricName                       string      `json:"MetricName"`
	Namespace                        string      `json:"Namespace"`
	Period                           int         `json:"Period"`
	Statistic                        string      `json:"Statistic"`
	StatisticType                    string      `json:"StatisticType"`
	Threshold                        int         `json:"Threshold"`
	TreatMissingData                 string      `json:"TreatMissingData"`
	Unit                             interface{} `json:"Unit"`
}

type webhookHandler struct {
}

type WebhookHandler interface {
	WebhookSNSAws() echo.HandlerFunc
}

func NewWebhookHandler() WebhookHandler {
	return &webhookHandler{}
}

func (d *webhookHandler) WebhookSNSAws() echo.HandlerFunc {
	return func(c echo.Context) error {
		// ctx := c.Request().Context()

		// var testVariable string

		// fmt.Println("aaaa")

		getJsonBody := c.Request().Body

		defer getJsonBody.Close()

		jsonBody, _ := io.ReadAll(getJsonBody)

		// var snsMessage SNSMessage

		// err := json.Unmarshal(jsonBody, &snsMessage)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	return c.JSON(http.StatusBadRequest, err)
		// }

		// // Extract AlarmName and Dimension value
		// newStateReason := snsMessage.Message.NewStateReason
		// dimensionValue := ""
		// for _, dimension := range snsMessage.Message.Trigger.Dimensions {
		// 	if dimension.Name == "TableName" {
		// 		dimensionValue = dimension.Value
		// 		break
		// 	}
		// }
		var snsMessage map[string]interface{}

		err := json.Unmarshal([]byte(jsonBody), &snsMessage)
		if err != nil {
			fmt.Println("Error:", err)
			return c.JSON(http.StatusBadRequest, err)
		}

		message, ok := snsMessage["Message"].(string)
		if !ok {
			fmt.Println("Error: Message field not found or not a string")
			return c.JSON(http.StatusBadRequest, ("Message field not found or not a string"))
		}

		var parsedMessage map[string]interface{}
		err = json.Unmarshal([]byte(message), &parsedMessage)
		if err != nil {
			fmt.Println("Error parsing Message JSON:", err)
			return c.JSON(http.StatusBadRequest, err)
		}

		alarmName, ok := parsedMessage["AlarmName"].(string)
		if !ok {
			fmt.Println("Error: AlarmName field not found or not a string")
			return c.JSON(http.StatusBadRequest, ("AlarmName field not found or not a string"))
		}

		newStateReason, ok := parsedMessage["NewStateReason"].(string)
		if !ok {
			fmt.Println("Error: NewStateReason field not found or not a string")
			return c.JSON(http.StatusBadRequest, ("NewStateReason field not found or not a string"))
		}

		stateChangeTime, ok := parsedMessage["StateChangeTime"].(string)
		if !ok {
			fmt.Println("Error: StateChangeTime field not found or not a string")
			return c.JSON(http.StatusBadRequest, ("StateChangeTime field not found or not a string"))
		}

		trigger, ok := parsedMessage["Trigger"].(map[string]interface{})
		if !ok {
			fmt.Println("Error: Trigger field not found or not a string")
			return c.JSON(http.StatusBadRequest, ("Message field not found or not a string"))
		}

		jsonString, err := json.Marshal(trigger)
		if err != nil {
			fmt.Println("Error:", err)
		}

		var parsedTrigger DynamoDBMetric
		err = json.Unmarshal([]byte(jsonString), &parsedTrigger)
		if err != nil {
			fmt.Println("Error parsing Message JSON:", err)
			return c.JSON(http.StatusBadRequest, err)
		}

		var tableName string
		if len(parsedTrigger.Dimensions) > 0 {
			tableName = parsedTrigger.Dimensions[0].Value
		} else {
			tableName = "No Table Name Found!"
		}

		// fmt.Println("AlarmName:", alarmName)
		// fmt.Println("NewStateReason:", newStateReason)
		// fmt.Println("StateChangeTime:", stateChangeTime)
		// fmt.Println(tableName)

		payload := map[string]interface{}{
			"blocks": []map[string]interface{}{
				{
					"type": "section",
					"text": map[string]interface{}{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Alert Name:*\n*%s*", alarmName),
					},
				},
				{
					"type": "section",
					"fields": []map[string]interface{}{
						{
							"type": "mrkdwn",
							"text": fmt.Sprintf("*Triggered Time:*\n%s", stateChangeTime),
						},
						{
							"type": "mrkdwn",
							"text": fmt.Sprintf("*Table Name:*\n%s", tableName),
						},
					},
				},
				{
					"type": "section",
					"text": map[string]interface{}{
						"type": "mrkdwn",
						"text": fmt.Sprintf("*Cause*:\n%s", newStateReason),
					},
				},
			},
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error Convert Slack Message to JSON:", err)
			return c.JSON(http.StatusBadRequest, err)
		}

		webhookURL := os.Getenv("CLOUDWATCH_SLACK_WEBHOOK")

		resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Println("Error Sending Message to Slack", err)
			return c.JSON(http.StatusBadRequest, err)
		}
		defer resp.Body.Close()

		return c.JSON(http.StatusOK, payload)
	}
}