// touched for cleanup
package Generation

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/joho/godotenv"
	"github.com/nigdanil/Quality_Control_System_Generation/model"
)

var apiKey, apiURL string

func init() {
	// Загрузка переменных из .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️  Не удалось загрузить .env файл")
	}

	apiKey = os.Getenv("API_KEY")
	apiURL = os.Getenv("API_URL")

	if apiKey == "" || apiURL == "" {
		fmt.Println("❌ Отсутствуют переменные API_KEY или API_URL в .env")
		os.Exit(1)
	}
}

// encodeImageBase64 кодирует изображение в base64 (PNG).
func encodeImageBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// loadImageFromURL загружает изображение по URL.
func loadImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки URL: %v", err)
	}
	defer resp.Body.Close()

	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования изображения: %v", err)
	}
	return img, nil
}

// GenerateImage выполняет генерацию изображения по данным лога.
func GenerateImage(log *model.GenerationLog) (string, float64, error) {
	// Загружаем изображения
	refImg1, err := loadImageFromURL(log.RefImg1URL)
	if err != nil {
		return "", 0, fmt.Errorf("ref_img1 error: %v", err)
	}
	refImg2, err := loadImageFromURL(log.RefImg2URL)
	if err != nil {
		return "", 0, fmt.Errorf("ref_img2 error: %v", err)
	}

	refImg1Base64, _ := encodeImageBase64(refImg1)
	refImg2Base64, _ := encodeImageBase64(refImg2)

	// Формируем JSON payload
	payload := map[string]interface{}{
		"inputs": map[string]interface{}{
			"ref_image1":          refImg1Base64,
			"ref_image2":          refImg2Base64,
			"ref_task1":           log.RefTask1,
			"ref_task2":           log.RefTask2,
			"ref_weight1":         log.RefWeight1,
			"ref_weight2":         log.RefWeight2,
			"prompt":              log.Prompt,
			"neg_prompt":          log.NegPrompt,
			"seed":                log.Seed,
			"width":               log.Width,
			"height":              log.Height,
			"ref_res":             log.RefRes,
			"num_steps":           log.Steps,
			"guidance":            log.Guidance,
			"true_cfg":            log.TrueCFG,
			"cfg_start_step":      log.CFGStartStep,
			"cfg_end_step":        log.CFGEndStep,
			"neg_guidance":        log.NegGuidance,
			"first_step_guidance": log.FirstStepGuidance,
		},
	}

	jsonData, _ := json.Marshal(payload)

	// Отправка запроса
	start := time.Now()
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	duration := time.Since(start).Seconds()

	if err != nil {
		return "", duration, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	// Обработка ответа
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", duration, fmt.Errorf("ошибка ответа [%d]: %s", resp.StatusCode, string(bodyBytes))
	}

	var respData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		return "", duration, fmt.Errorf("ошибка JSON-ответа: %v", err)
	}

	base64Img, ok := respData["image_base64"].(string)
	if !ok || base64Img == "" {
		return "", duration, fmt.Errorf("не найдено поле image_base64 в ответе")
	}

	// Декодирование base64 и сохранение
	imgData, _ := base64.StdEncoding.DecodeString(base64Img)
	imgReader := bytes.NewReader(imgData)
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return "", duration, fmt.Errorf("ошибка декодирования изображения: %v", err)
	}

	// Сохраняем файл
	now := time.Now().Format("02_01_2006_15_04_05")
	fileName := fmt.Sprintf("img_result_%s.png", now)
	savePath := filepath.Join("images", fileName)

	outFile, _ := os.Create(savePath)
	defer outFile.Close()
	png.Encode(outFile, img)

	return fileName, duration, nil
}
