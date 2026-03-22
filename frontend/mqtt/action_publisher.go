package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

// ActionPayload 动作下发Payload结构
type ActionPayload struct {
	ActionID      string                 `json:"action_id"`
	ActionName    string                 `json:"action_name"`
	DurationMs    int                    `json:"duration_ms"`
	Priority      int                    `json:"priority"`
	Parameters    map[string]interface{} `json:"parameters"`
	MotorCommands string                 `json:"motor_commands"`
	Timestamp     string                 `json:"timestamp"`
}

// PublishAction 发布动作指令到设备
// @param client MQTT客户端
// @param deviceID 设备ID
// @param payload 动作载荷
// @return error 错误信息
func PublishAction(client mqtt.Client, deviceID string, payload interface{}) error {
	topic := fmt.Sprintf("/miniclaw/%s/down/action", deviceID)

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化动作载荷失败: %w", err)
	}

	token := client.Publish(topic, 0, false, data)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("发布MQTT消息失败: %w", token.Error())
	}

	log.Printf("[MQTT] 动作已下发到设备 %s: %s", deviceID, string(data))
	return nil
}

// PublishCommand 发布通用指令到设备
// @param client MQTT客户端
// @param deviceID 设备ID
// @param cmdType 指令类型
// @param payload 指令载荷
// @return error 错误信息
func PublishCommand(client mqtt.Client, deviceID, cmdType string, payload interface{}) error {
	topic := fmt.Sprintf("/miniclaw/%s/down/cmd", deviceID)

	cmd := map[string]interface{}{
		"cmd_id":    fmt.Sprintf("%s-%s-%d", cmdType, deviceID, time.Now().Unix()),
		"cmd_type":  cmdType,
		"timestamp": time.Now().Format(time.RFC3339),
		"payload":   payload,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("序列化指令失败: %w", err)
	}

	token := client.Publish(topic, 0, false, data)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("发布MQTT消息失败: %w", token.Error())
	}

	log.Printf("[MQTT] 指令已下发到设备 %s: %s", deviceID, string(data))
	return nil
}

// PublishBoost 发布心情激励到设备
// @param client MQTT客户端
// @param deviceID 设备ID
// @param boostType 激励类型: food, play, praise, music
// @param mood 当前心情值
// @param expression 当前表情
// @return error 错误信息
func PublishBoost(client mqtt.Client, deviceID, boostType string, mood int, expression string) error {
	payload := map[string]interface{}{
		"type":       "boost",
		"boost_type": boostType,
		"mood":       mood,
		"expression": expression,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	return PublishAction(client, deviceID, payload)
}

// PublishExpressionChange 发布表情变化到设备
// @param client MQTT客户端
// @param deviceID 设备ID
// @param expression 表情名称
// @return error 错误信息
func PublishExpressionChange(client mqtt.Client, deviceID, expression string) error {
	payload := map[string]interface{}{
		"type":       "expression",
		"expression": expression,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	return PublishAction(client, deviceID, payload)
}

// PublishMoodUpdate 发布心情更新到设备
// @param client MQTT客户端
// @param deviceID 设备ID
// @param mood 心情值
// @param energy 能量值
// @param hunger 饥饿值
// @return error 错误信息
func PublishMoodUpdate(client mqtt.Client, deviceID string, mood, energy, hunger int) error {
	payload := map[string]interface{}{
		"type":      "mood_update",
		"mood":      mood,
		"energy":    energy,
		"hunger":    hunger,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	return PublishAction(client, deviceID, payload)
}
