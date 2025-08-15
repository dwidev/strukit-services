package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strukit-services/pkg/logger"

	"github.com/sirupsen/logrus"
	"google.golang.org/genai"
)

func Run() *Manager {
	manager := new(Manager)
	manager.Context = context.Background()

	client, err := genai.NewClient(manager.Context, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		panic(fmt.Sprintf("panic when run llm Manager with error : %s", err))
	}

	manager.client = client
	return manager
}

type Manager struct {
	Context context.Context
	client  *genai.Client
}

func (m *Manager) ScanReceiptWithImage(image []byte) (*ReceiptResponse, error) {
	contents := []*genai.Content{
		genai.NewContentFromBytes(image, "image/jpeg", genai.RoleUser),
	}

	res, err := m.generateContent(contents)
	if err != nil {
		return nil, err
	}

	result := res.Text()
	receipt, err := m.parse(result)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (m *Manager) ScanReceiptFromOCR(ocrRaw string) (*ReceiptResponse, error) {
	ocrRaw = fmt.Sprintf(`### DATA OCR ###%s### END DATA OCR ###`, ocrRaw)
	contents := []*genai.Content{
		genai.NewContentFromText(ocrRaw, genai.RoleUser),
	}

	res, err := m.generateContent(contents)
	if err != nil {
		return nil, err
	}

	result := res.Text()
	receipt, err := m.parse(result)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (m *Manager) generateContent(contents []*genai.Content) (*genai.GenerateContentResponse, error) {
	ccfg := m.contentCfg()
	res, err := m.client.Models.GenerateContent(m.Context, "gemini-2.0-flash", contents, ccfg)
	if err != nil {
		logger.Log.LLM(m.Context).Errorf("llm error generate content, error : %s", err)
		return nil, err
	}

	logger.Log.LLM(m.Context).WithField(
		"tokenUsage", res.UsageMetadata,
	).Info("token usage")
	return res, nil
}

func (m *Manager) contentCfg() *genai.GenerateContentConfig {
	sys := SystemPrompt()

	structuredOutput := m.StructuredOutput()
	return &genai.GenerateContentConfig{
		ResponseMIMEType:  "application/json",
		SystemInstruction: genai.NewContentFromText(sys, genai.RoleModel),
		ResponseSchema:    structuredOutput,
	}
}

func (m *Manager) parse(result string) (*ReceiptResponse, error) {
	receipt := new(ReceiptResponse)
	err := json.Unmarshal([]byte(result), receipt)
	if err != nil {
		logger.Log.LLM(m.Context).WithField("raw", receipt).Errorf("llm error parse content result, error : %s", err)
		return nil, err
	}

	logger.Log.LLM(m.Context).WithFields(logrus.Fields{"data": receipt, "rawData": result}).Infof("success read the receipt wiht data")

	return receipt, nil
}
