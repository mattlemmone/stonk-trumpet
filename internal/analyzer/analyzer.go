package analyzer

import (
	"context"
	"fmt"
	"log"
	"stonk-trumpet/pkg/truthsocial"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// Sentiment represents the analyzed sentiment of a status.
type Sentiment string

const (
	Positive Sentiment = "positive"
	Negative Sentiment = "negative"
	Neutral  Sentiment = "neutral"
	Error    Sentiment = "error"
)

// AnalysisResult holds the sentiment and potential market impact.
type AnalysisResult struct {
	Sentiment  Sentiment
	IsRelevant bool // Indicates if the tweet is relevant to stock market analysis
}

// Analyzer defines the interface for sentiment analysis.
type Analyzer interface {
	AnalyzeSentiment(status *truthsocial.Status) (AnalysisResult, error)
}

// openaiAnalyzer implements the Analyzer interface using OpenAI.
type openaiAnalyzer struct {
	client   *openai.Client
	mockMode bool
}

// NewOpenAIAnalyzer creates a new openaiAnalyzer instance.
func NewOpenAIAnalyzer(apiKey string) Analyzer {
	client := openai.NewClient(apiKey)
	return &openaiAnalyzer{
		client:   client,
		mockMode: false,
	}
}

// NewMockAnalyzer creates an analyzer in mock mode for testing.
func NewMockAnalyzer() Analyzer {
	return &openaiAnalyzer{
		mockMode: true,
	}
}

// AnalyzeSentiment analyzes the sentiment of a given status content.
func (a *openaiAnalyzer) AnalyzeSentiment(status *truthsocial.Status) (AnalysisResult, error) {
	log.Printf("Analyzing sentiment for status %s...", status.ID)

	// If in mock mode, return mock results based on content
	if a.mockMode {
		return a.generateMockAnalysis(status), nil
	}

	// Strip HTML from content
	content := stripHTML(status.Content)

	// Prepare prompt for OpenAI
	prompt := fmt.Sprintf(`
Analyze the following social media post for its sentiment regarding the stock market or business/economic implications.
Determine if it is:
1) Relevant to stock market, business, or economic matters
2) If relevant, whether the sentiment is positive, negative, or neutral

Tweet: "%s"

Reply in the following JSON format:
{
  "isRelevant": true/false,
  "sentiment": "positive"/"negative"/"neutral"/"error"
}
`, content)

	// Create chat completion request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a financial sentiment analyzer. You analyze social media posts to determine if they are relevant to the stock market and their sentiment.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.1, // Lower temperature for more consistent results
		},
	)

	if err != nil {
		log.Printf("Error from OpenAI API: %v", err)
		return AnalysisResult{Sentiment: Error, IsRelevant: false}, fmt.Errorf("OpenAI API error: %w", err)
	}

	// Parse the response
	response := resp.Choices[0].Message.Content
	log.Printf("OpenAI response: %s", response)

	// Simple parsing for demonstration (in production, use proper JSON parsing)
	var result AnalysisResult

	// Check for relevance
	if strings.Contains(response, `"isRelevant": true`) {
		result.IsRelevant = true
	} else {
		result.IsRelevant = false
	}

	// Check for sentiment
	if strings.Contains(response, `"sentiment": "positive"`) {
		result.Sentiment = Positive
	} else if strings.Contains(response, `"sentiment": "negative"`) {
		result.Sentiment = Negative
	} else if strings.Contains(response, `"sentiment": "error"`) {
		result.Sentiment = Error
	} else {
		result.Sentiment = Neutral
	}

	log.Printf("Analysis result for %s: relevant=%t, sentiment=%s",
		status.ID, result.IsRelevant, result.Sentiment)

	return result, nil
}

// generateMockAnalysis returns mock analysis results for testing
func (a *openaiAnalyzer) generateMockAnalysis(status *truthsocial.Status) AnalysisResult {
	log.Println("Using mock mode: Generating fake analysis result")

	content := strings.ToLower(stripHTML(status.Content))

	// Check for keywords to determine relevance
	isRelevant := false
	keywords := []string{"stock", "market", "economy", "business", "trade", "finance", "economic"}
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			isRelevant = true
			break
		}
	}

	// Determine sentiment based on keywords
	sentiment := Neutral
	if isRelevant {
		positiveWords := []string{"strong", "up", "gain", "boom", "great", "wonderful", "positive"}
		negativeWords := []string{"down", "fall", "crash", "weak", "trouble", "bad", "negative"}

		for _, word := range positiveWords {
			if strings.Contains(content, word) {
				sentiment = Positive
				break
			}
		}

		for _, word := range negativeWords {
			if strings.Contains(content, word) {
				sentiment = Negative
				break
			}
		}
	}

	result := AnalysisResult{
		IsRelevant: isRelevant,
		Sentiment:  sentiment,
	}

	log.Printf("Mock analysis result for %s: relevant=%t, sentiment=%s",
		status.ID, result.IsRelevant, result.Sentiment)
	return result
}

// stripHTML removes HTML tags from a string.
// TODO: Consider using a proper HTML parsing library (e.g., htmlquery or goquery) for robustness.
func stripHTML(s string) string {
	dst := make([]byte, 0, len(s))
	var inTag bool
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			inTag = true
			continue
		}
		if s[i] == '>' {
			inTag = false
			continue
		}
		if !inTag {
			dst = append(dst, s[i])
		}
	}
	return string(dst)
}

// TODO: Add imports for strings
