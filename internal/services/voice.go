package services

import (
	"fmt"
	"mime/multipart"
)

// AudioToText 处理音频转文字 (ASR)
// 在这里调用百度语音或阿里云 ASR 的 SDK
func AudioToText(file *multipart.FileHeader) (string, error) {
	return "", fmt.Errorf("ASR 服务尚未配置，请直接传输文本指令")
}

// TextToAudio 将文字转为语音 URL (TTS)
func TextToAudio(text string) string {
	// 这里返回前端可调用的第三方免费 TTS 接口作为 MVP 替身
	return fmt.Sprintf("https://dict.youdao.com/dictvoice?audio=%s&le=zh", text)
}
