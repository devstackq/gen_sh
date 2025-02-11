package speech

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/devstackq/gen_sh/internal/logger"
)

// GenerateAudio - генерирует аудиофайл на основе текста с помощью Google TTS или espeak.
func GenerateAudio(text string) (string, error) {
	// Определяем путь для сохранения аудиофайла
	audioPath := filepath.Join("/tmp", fmt.Sprintf("audio_%d.mp3", time.Now().Unix()))

	// Попробуем использовать Google TTS (если установлен gtts-cli)
	err := generateWithGoogleTTS(text, audioPath)
	if err != nil {
		logger.LogError(fmt.Sprint("Google TTS недоступен, переключаемся на espeak", err))
		// Если Google TTS недоступен, используем espeak
		err = generateWithEspeak(text, audioPath)
		if err != nil {
			logger.LogError(fmt.Sprint("Ошибка генерации аудио с espeak", err))
			return "", err
		}
	}

	logger.LogInfo(fmt.Sprint("Аудиофайл успешно сгенерирован path - ", audioPath))
	return audioPath, nil
}

// generateWithGoogleTTS - генерация аудио с помощью Google TTS (gtts-cli)
func generateWithGoogleTTS(text, audioPath string) error {
	cmd := exec.Command("gtts-cli", text, "--output", audioPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка при использовании Google TTS: %v, output: %s", err, string(output))
	}

	return nil
}

// generateWithEspeak - генерация аудио с помощью espeak (системного TTS)
func generateWithEspeak(text, audioPath string) error {
	tempWav := audioPath + ".wav"

	cmd := exec.Command("espeak", text, "-w", tempWav)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка при использовании espeak: %v, output: %s", err, string(output))
	}

	// Конвертируем WAV в MP3 с помощью ffmpeg
	convertCmd := exec.Command("ffmpeg", "-i", tempWav, "-q:a", "2", "-y", audioPath)
	convertOutput, convertErr := convertCmd.CombinedOutput()
	if convertErr != nil {
		return fmt.Errorf("ошибка при конвертации в MP3: %v, output: %s", convertErr, string(convertOutput))
	}

	// Удаляем временный WAV-файл
	_ = os.Remove(tempWav)

	return nil
}
