package services

import (
	"context"
	"encoding/json"
	"fmt"
	"pbkk-quizlit-backend/internal/models"
	"strings"
	"time"

	"bytes"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type AIService struct {
	client         *openai.Client
	senopatiClient *SenopatiClient
	logger         *logrus.Logger
	apiKey         string
	useOpenAI      bool
	useSenopati    bool
	// RAG components
	rag       *RAGService
	enableRAG bool
}

func NewAIService(apiKey string) *AIService {
	var client *openai.Client
	useOpenAI := false

	if apiKey != "" && apiKey != "your_openai_api_key_here" {
		client = openai.NewClient(apiKey)
		useOpenAI = true
	}

	// Initialize Senopati client
	senopatiClient := NewSenopatiClient()
	// Senopati is always available (no API key required)
	useSenopati := true

	logger := logrus.New()

	ai := &AIService{
		client:         client,
		senopatiClient: senopatiClient,
		logger:         logger,
		apiKey:         apiKey,
		useOpenAI:      useOpenAI,
		useSenopati:    useSenopati,
	}
	// initialize lightweight RAG with hash embedding (works without external deps)
	ai.rag = NewRAGService(HashEmbedding{})
	// Enable RAG to intelligently select relevant content chunks
	ai.enableRAG = true
	return ai
}

// NewAIServiceWithOptions allows configuring RAG usage via flag
func NewAIServiceWithOptions(apiKey string, enableRAG bool) *AIService {
	ai := NewAIService(apiKey)
	ai.enableRAG = enableRAG
	return ai
}

// GenerateQuizFromContent generates quiz questions using AI or free alternatives
func (ai *AIService) GenerateQuizFromContent(content string, req *models.QuizGenerationRequest) (*models.Quiz, error) {
	if req.QuestionCount == 0 {
		req.QuestionCount = 10 // Default to 10 questions
	}

	ai.logger.Infof("Generating quiz with %d questions for difficulty: %s", req.QuestionCount, req.Difficulty)

	var questions []models.Question
	var err error

	// Build RAG index and retrieve top context chunks to ground prompts
	if ai.rag != nil && ai.enableRAG {
		ai.logger.Info("Using RAG to select relevant content chunks")
		// Use title or a generated docID
		docID := req.Title
		if strings.TrimSpace(docID) == "" {
			docID = uuid.New().String()
		}
		if err := ai.rag.BuildIndex(docID, content); err == nil {
			// retrieve with query from description+difficulty for better intent
			query := strings.TrimSpace(req.Description + " " + req.Difficulty)
			if query == "" {
				query = "generate quiz key concepts"
			}
			top, _ := ai.rag.Retrieve(query, 8) // Get more chunks for better coverage
			if len(top) > 0 {
				// Assemble clean context without metadata prefixes
				var b strings.Builder
				totalLength := 0
				maxTotalLength := 8000 // Keep under 8KB for better performance

				for i, it := range top {
					chunkText := strings.TrimSpace(it.Text)

					// Skip if would exceed limit
					if totalLength+len(chunkText) > maxTotalLength {
						break
					}

					if i > 0 {
						b.WriteString("\n\n")
					}
					b.WriteString(chunkText)
					totalLength += len(chunkText) + 2
				}

				// Replace content with curated chunks
				content = b.String()
				ai.logger.Infof("RAG selected %d chunks, total length: %d chars", len(top), len(content))
			}
		} else {
			ai.logger.Warnf("RAG indexing failed: %v", err)
		}
	}

	// Try Senopati first (ITS local LLM)
	if ai.useSenopati {
		questions, err = ai.generateWithSenopati(content, req)
		if err != nil {
			ai.logger.Errorf("Senopati failed: %v", err)
			return nil, fmt.Errorf("quiz generation failed: %w", err)
		}
	} else {
		return nil, fmt.Errorf("no AI provider configured - Senopati is required")
	}

	// Create quiz object
	quiz := &models.Quiz{
		ID:             uuid.New().String(),
		Title:          req.Title,
		Description:    req.Description,
		Questions:      questions,
		Difficulty:     req.Difficulty,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		TotalQuestions: len(questions),
	}

	ai.logger.Infof("Successfully generated quiz with %d questions", len(questions))
	return quiz, nil
}

// generateWithOpenAI uses OpenAI GPT for quiz generation
func (ai *AIService) generateWithOpenAI(content string, req models.CreateQuizRequest) ([]models.Question, error) {
	prompt := ai.createPrompt(content, req)

	// Call OpenAI API
	resp, err := ai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert quiz generator. Generate high-quality multiple choice questions based on the provided content. Return ONLY valid JSON without any additional text or formatting.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   2000,
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse the response
	questions, err := ai.parseAIResponse(resp.Choices[0].Message.Content)
	if err != nil {
		ai.logger.Errorf("Failed to parse OpenAI response: %v", err)
		return nil, err
	}

	return questions, nil
}

// generateWithSenopati uses ITS Senopati local LLM for quiz generation
func (ai *AIService) generateWithSenopati(content string, req *models.QuizGenerationRequest) ([]models.Question, error) {
	ai.logger.Info("Using ITS Senopati LLM for quiz generation")

	// Truncate content if too large (max ~10KB for better performance)
	maxContentLength := 10000
	if len(content) > maxContentLength {
		ai.logger.Warnf("Content too large (%d chars), truncating to %d", len(content), maxContentLength)
		content = content[:maxContentLength] + "\\n[Content truncated due to size...]"
	}

	prompt := ai.buildPrompt(content, req)

	// Get available models from API
	modelsResp, err := ai.senopatiClient.ListModels()
	var model string
	if err == nil && len(modelsResp.Models) > 0 {
		// Prefer larger models for better quality and completion
		// Priority: qwen2.5:14b > llama3:latest > qwen2.5:7b > gemma:7b
		preferredModels := []string{"qwen2.5:14b", "llama3:latest", "qwen2.5:7b", "llama3", "qwen2.5"}

		for _, preferred := range preferredModels {
			for _, available := range modelsResp.Models {
				if strings.Contains(available, preferred) {
					model = available
					ai.logger.Infof("Using Senopati model: %s (selected from %d available models)", model, len(modelsResp.Models))
					goto modelSelected
				}
			}
		}

		// If no preferred model found, use first available
		model = modelsResp.Models[0]
		ai.logger.Infof("Using Senopati model: %s (default from %d available models)", model, len(modelsResp.Models))
	modelSelected:
	} else {
		// Fallback to default model from docs
		model = "qwen2.5:14b"
		ai.logger.Warnf("Could not fetch models (error: %v), using default: %s", err, model)
	}

	// Call Senopati Generate endpoint
	// Scale max tokens based on question count (each question ~300 tokens)
	maxTokens := req.QuestionCount * 400
	if maxTokens < 4000 {
		maxTokens = 4000 // Minimum for safety
	}
	if maxTokens > 8000 {
		maxTokens = 8000 // Cap at 8000
	}
	resp, err := ai.senopatiClient.GenerateText(model, prompt, 0.7, maxTokens)
	if err != nil {
		return nil, fmt.Errorf("senopati API error: %w", err)
	}

	// Parse the response using the same parser
	questions, err := ai.parseAIResponse(resp.Response)
	if err != nil {
		ai.logger.Errorf("Failed to parse Senopati response: %v", err)
		return nil, err
	}

	return questions, nil
}

// generateWithOllama uses Ollama API for free local AI processing
func (ai *AIService) generateWithOllama(content string, req *models.QuizGenerationRequest) ([]models.Question, error) {
	ai.logger.Info("Using Ollama for quiz generation")

	// Ollama API endpoint (default local installation)
	url := "http://localhost:11434/api/generate"

	prompt := ai.buildPrompt(content, req)

	requestBody := map[string]interface{}{
		"model":  "llama2", // Default model, can be made configurable
		"prompt": prompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		ai.logger.Errorf("Ollama API request failed: %v", err)
		return nil, fmt.Errorf("ollama API error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		ai.logger.Errorf("Ollama API error response: %s", string(body))
		return nil, fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode ollama response: %w", err)
	}

	// Parse the response using the same parser
	questions, err := ai.parseAIResponse(ollamaResp.Response)
	if err != nil {
		ai.logger.Errorf("Failed to parse Ollama response: %v", err)
		return nil, err
	}

	return questions, nil
}

// generateIntelligentQuestions creates quiz questions using rule-based generation
func (ai *AIService) generateIntelligentQuestions(content string, req *models.QuizGenerationRequest) []models.Question {
	ai.logger.Info("Using enhanced intelligent rule-based quiz generation")
	ai.logger.Infof("Content length: %d characters", len(content))

	// Enhanced content analysis
	sentences := ai.extractSentences(content)
	ai.logger.Infof("Extracted %d sentences", len(sentences))

	keywords := ai.extractKeywords(content)
	ai.logger.Infof("Extracted %d keywords", len(keywords))

	concepts := ai.extractConcepts(content)
	ai.logger.Infof("Extracted %d concepts", len(concepts))

	// Filter sentences to only informative ones
	validSentences := ai.filterInformativeSentences(sentences, keywords)
	ai.logger.Infof("Filtered to %d valid sentences", len(validSentences))

	if len(validSentences) == 0 {
		ai.logger.Warn("No valid sentences found for question generation")
		return []models.Question{}
	}

	var questions []models.Question
	targetCount := req.QuestionCount
	if targetCount <= 0 {
		targetCount = 5
	}

	usedSentences := make(map[string]bool)

	// Generate diverse, high-quality questions
	for i := 0; i < targetCount*2 && len(questions) < targetCount; i++ {
		// Pick sentence that hasn't been used yet
		var sentence string
		for _, s := range validSentences {
			if !usedSentences[s] {
				sentence = s
				usedSentences[s] = true
				break
			}
		}

		if sentence == "" && len(validSentences) > 0 {
			// All used, allow reuse
			sentence = validSentences[i%len(validSentences)]
		}

		if sentence == "" {
			continue
		}

		var question models.Question

		// Only generate multiple-choice and true/false questions
		// Fill-in-blank disabled because it has no options (causes UI issues)
		switch i % 2 {
		case 0: // Multiple choice based on key sentences
			question = ai.generateMultipleChoiceFromSentence(sentence, keywords, concepts)
		case 1: // True/false questions
			question = ai.generateTrueFalseFromSentence(sentence, keywords)
			// case 2: // Fill in the blank - DISABLED
			// 	question = ai.generateFillInTheBlank(sentence, keywords)
		}

		// Validate question quality before adding
		if ai.isValidQuestion(question, content) {
			question.ID = uuid.New().String()
			questions = append(questions, question)
		}
	}

	return questions
}

// Helper methods for rule-based generation
func (ai *AIService) extractSentences(content string) []string {
	// Better sentence extraction with length limits
	sentences := []string{}

	// Split by common sentence endings
	content = strings.ReplaceAll(content, "! ", ".|")
	content = strings.ReplaceAll(content, "? ", ".|")
	content = strings.ReplaceAll(content, ". ", ".|")

	parts := strings.Split(content, "|")

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		words := strings.Fields(trimmed)
		wordCount := len(words)

		// Filter: min 5 words, max 30 words, min 30 chars, max 200 chars
		if wordCount >= 5 && wordCount <= 30 && len(trimmed) >= 30 && len(trimmed) <= 200 {
			// Ensure it ends with proper punctuation
			if !strings.HasSuffix(trimmed, ".") && !strings.HasSuffix(trimmed, "!") && !strings.HasSuffix(trimmed, "?") {
				trimmed += "."
			}

			// Check if sentence is meaningful (not just numbers/symbols)
			letterCount := 0
			for _, r := range trimmed {
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
					letterCount++
				}
			}

			// At least 50% should be letters
			if letterCount >= len(trimmed)/2 {
				sentences = append(sentences, trimmed)
			}
		}
	}

	return sentences
}

func (ai *AIService) extractKeywords(content string) []string {
	// Better keyword extraction with frequency analysis
	words := strings.Fields(content)
	wordFreq := make(map[string]int)

	// Common stop words to ignore (Indonesian and English)
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "as": true, "is": true, "was": true,
		"are": true, "were": true, "been": true, "be": true, "have": true, "has": true,
		"had": true, "do": true, "does": true, "did": true, "will": true, "would": true,
		"could": true, "should": true, "may": true, "might": true, "must": true,
		"can": true, "this": true, "that": true, "these": true, "those": true,
		"yang": true, "dan": true, "atau": true, "adalah": true, "ini": true, "itu": true,
		"dari": true, "ke": true, "di": true, "untuk": true, "dengan": true, "pada": true,
	}

	for _, word := range words {
		cleaned := strings.ToLower(strings.Trim(word, ".,!?;:()[]{}\"'"))
		// Keep words that are substantial and not stop words
		if len(cleaned) > 3 && !stopWords[cleaned] {
			wordFreq[cleaned]++
		}
	}

	// Sort by frequency
	type keywordScore struct {
		word  string
		score int
	}
	var scored []keywordScore
	for word, freq := range wordFreq {
		scored = append(scored, keywordScore{word, freq})
	}

	// Simple sort by score
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Return top keywords
	var result []string
	for i := 0; i < len(scored) && i < 30; i++ {
		result = append(result, scored[i].word)
	}

	return result
}

func (ai *AIService) extractConcepts(content string) []string {
	// Extract multi-word concepts and phrases
	concepts := []string{}
	sentences := strings.Split(content, ".")

	for _, sentence := range sentences {
		words := strings.Fields(sentence)
		// Look for 2-3 word phrases (likely concepts)
		for i := 0; i < len(words)-1; i++ {
			w1 := strings.Trim(words[i], ".,!?;:")
			w2 := strings.Trim(words[i+1], ".,!?;:")

			if len(w1) > 3 && len(w2) > 3 {
				phrase := w1 + " " + w2
				if len(phrase) > 8 {
					concepts = append(concepts, phrase)
				}
			}
		}
	}

	// Limit and deduplicate
	seen := make(map[string]bool)
	var result []string
	for _, concept := range concepts {
		conceptLower := strings.ToLower(concept)
		if !seen[conceptLower] && len(result) < 20 {
			seen[conceptLower] = true
			result = append(result, concept)
		}
	}

	return result
}

func (ai *AIService) filterInformativeSentences(sentences []string, keywords []string) []string {
	// Filter sentences that contain important keywords
	var filtered []string

	for _, sentence := range sentences {
		sentenceLower := strings.ToLower(sentence)
		keywordCount := 0

		// Count how many keywords appear in this sentence
		for _, keyword := range keywords {
			if strings.Contains(sentenceLower, strings.ToLower(keyword)) {
				keywordCount++
			}
		}

		// Keep sentences with at least 2 keywords and substantial length
		if keywordCount >= 2 && len(sentence) > 40 {
			filtered = append(filtered, sentence)
		}
	}

	// If too strict, lower the bar
	if len(filtered) < 5 && len(sentences) > 0 {
		filtered = sentences
	}

	return filtered
}

func (ai *AIService) isValidQuestion(question models.Question, content string) bool {
	// Validate question quality
	contentLower := strings.ToLower(content)
	questionLower := strings.ToLower(question.Text)

	// Question should have substantial length
	if len(question.Text) < 15 {
		return false
	}

	// For multiple choice, ensure we have options
	if question.Type == "multiple-choice" {
		if len(question.Options) < 4 {
			return false
		}

		// Check if correct answer appears in content
		if len(question.Correct) > 0 {
			correctLower := strings.ToLower(question.Correct)
			if !strings.Contains(contentLower, correctLower) {
				return false
			}
		}
	}

	// Avoid generic questions and options
	genericPhrases := []string{"concept a", "concept b", "option 1", "option 2", "option 3", "option 4"}
	for _, phrase := range genericPhrases {
		if strings.Contains(questionLower, phrase) {
			return false
		}
		// Check options too
		for _, option := range question.Options {
			if strings.Contains(strings.ToLower(option), phrase) {
				return false
			}
		}
	}

	return true
}

func (ai *AIService) generateMultipleChoiceFromSentence(sentence string, keywords []string, concepts []string) models.Question {
	// Create a question by identifying and replacing a key term
	words := strings.Fields(sentence)
	var targetWord string
	var questionText string

	// First try to find a keyword in the sentence
	for i, word := range words {
		cleaned := strings.Trim(word, ".,!?;:")
		for _, keyword := range keywords {
			if strings.EqualFold(cleaned, keyword) && len(cleaned) > 4 {
				targetWord = cleaned
				// Replace in original sentence
				wordsCopy := make([]string, len(words))
				copy(wordsCopy, words)
				wordsCopy[i] = "____"
				questionText = strings.Join(wordsCopy, " ")
				break
			}
		}
		if targetWord != "" {
			break
		}
	}

	// If no keyword found, try concepts
	if targetWord == "" {
		for _, concept := range concepts {
			if strings.Contains(sentence, concept) {
				targetWord = concept
				questionText = strings.Replace(sentence, concept, "____", 1)
				break
			}
		}
	}

	// Fallback: pick a meaningful word
	if targetWord == "" && len(words) > 5 {
		for i := 2; i < len(words)-2; i++ {
			cleaned := strings.Trim(words[i], ".,!?;:")
			if len(cleaned) > 4 && !strings.Contains(strings.ToLower(cleaned), "yang") && !strings.Contains(strings.ToLower(cleaned), "adalah") {
				targetWord = cleaned
				wordsCopy := make([]string, len(words))
				copy(wordsCopy, words)
				wordsCopy[i] = "____"
				questionText = strings.Join(wordsCopy, " ")
				break
			}
		}
	}

	if questionText == "" {
		questionText = sentence
		targetWord = "unknown"
	}

	// Generate realistic distractors from content
	options := []string{targetWord}
	usedOptions := make(map[string]bool)
	usedOptions[strings.ToLower(targetWord)] = true

	// Add similar keywords as distractors
	for _, keyword := range keywords {
		if len(options) >= 4 {
			break
		}
		keywordLower := strings.ToLower(keyword)
		if !usedOptions[keywordLower] && len(keyword) > 3 {
			options = append(options, keyword)
			usedOptions[keywordLower] = true
		}
	}

	// Add concepts as distractors if needed
	for _, concept := range concepts {
		if len(options) >= 4 {
			break
		}
		conceptLower := strings.ToLower(concept)
		if !usedOptions[conceptLower] && !strings.Contains(conceptLower, strings.ToLower(targetWord)) {
			options = append(options, concept)
			usedOptions[conceptLower] = true
		}
	}

	// If still not enough options, add generic but plausible ones
	genericOptions := []string{"None of the above", "All of the above", "Cannot be determined", "Not specified"}
	for _, generic := range genericOptions {
		if len(options) >= 4 {
			break
		}
		options = append(options, generic)
	}

	// Ensure we have exactly 4 options
	if len(options) < 4 {
		for len(options) < 4 {
			options = append(options, fmt.Sprintf("Option %d", len(options)+1))
		}
	} else if len(options) > 4 {
		options = options[:4]
	}

	// Trim question text if too long
	if len(questionText) > 150 {
		questionText = questionText[:147] + "..."
	}

	return models.Question{
		Type:     "multiple-choice",
		Text:     "Complete the sentence: " + questionText,
		Options:  options,
		Correct:  targetWord,
		Points:   1,
		Metadata: map[string]interface{}{"source": "rule-based-enhanced"},
	}
}

func (ai *AIService) generateTrueFalseFromSentence(sentence string, keywords []string) models.Question {
	// Create true/false by modifying factual statements
	questionText := sentence
	correct := "True"

	// 50% chance to make it false by intelligent negation
	words := strings.Fields(sentence)
	if len(words) > 5 && len(keywords) > 0 {
		// Strategy 1: Replace a keyword with another keyword (makes it false)
		for i, word := range words {
			cleaned := strings.Trim(word, ".,!?;:")
			// If this word is a keyword, replace with different keyword
			for j, kw := range keywords {
				if strings.EqualFold(cleaned, kw) && len(keywords) > j+1 {
					// Replace with a different keyword
					replacement := keywords[(j+1)%len(keywords)]
					if !strings.EqualFold(replacement, cleaned) {
						wordsCopy := make([]string, len(words))
						copy(wordsCopy, words)
						wordsCopy[i] = replacement
						questionText = strings.Join(wordsCopy, " ")
						correct = "False"
						break
					}
				}
			}
			if correct == "False" {
				break
			}
		}

		// Strategy 2: Add "not" to negate statement
		if correct == "True" && len(words) > 3 {
			// Find verb position (simplified - look for common verbs)
			for i := 1; i < len(words)-1; i++ {
				word := strings.ToLower(strings.Trim(words[i], ".,!?;:"))
				if word == "is" || word == "are" || word == "was" || word == "were" ||
					word == "dapat" || word == "adalah" || word == "merupakan" {
					wordsCopy := make([]string, len(words))
					copy(wordsCopy, words)
					wordsCopy[i] = words[i] + " not"
					questionText = strings.Join(wordsCopy, " ")
					correct = "False"
					break
				}
			}
		}
	}

	// Trim if too long
	if len(questionText) > 150 {
		questionText = questionText[:147] + "..."
	}

	return models.Question{
		Type:     "true-false",
		Text:     "True or False: " + questionText,
		Options:  []string{"True", "False"},
		Correct:  correct,
		Points:   1,
		Metadata: map[string]interface{}{"source": "rule-based-enhanced"},
	}
}

func (ai *AIService) generateFillInTheBlank(sentence string, keywords []string) models.Question {
	words := strings.Fields(sentence)
	var blank string
	var questionText string

	// Try to find a keyword to blank out
	for i, word := range words {
		cleaned := strings.Trim(word, ".,!?;:")
		for _, keyword := range keywords {
			if strings.EqualFold(cleaned, keyword) && len(cleaned) > 4 {
				blank = cleaned
				wordsCopy := make([]string, len(words))
				copy(wordsCopy, words)
				wordsCopy[i] = "______"
				questionText = strings.Join(wordsCopy, " ")
				break
			}
		}
		if blank != "" {
			break
		}
	}

	// Fallback: remove an important word
	if blank == "" && len(words) > 5 {
		for i := 2; i < len(words)-2; i++ {
			cleaned := strings.Trim(words[i], ".,!?;:")
			if len(cleaned) > 5 {
				blank = cleaned
				wordsCopy := make([]string, len(words))
				copy(wordsCopy, words)
				wordsCopy[i] = "______"
				questionText = strings.Join(wordsCopy, " ")
				break
			}
		}
	}

	if questionText == "" {
		questionText = sentence
		blank = "unknown"
	}

	// Trim if too long
	if len(questionText) > 150 {
		questionText = questionText[:147] + "..."
	}

	return models.Question{
		Type:     "fill-blank",
		Text:     "Fill in the blank: " + questionText,
		Options:  []string{}, // No options for fill-in-the-blank
		Correct:  blank,
		Points:   1,
		Metadata: map[string]interface{}{"source": "rule-based-enhanced"},
	}
}

func (ai *AIService) buildPrompt(content string, req *models.QuizGenerationRequest) string {
	prompt := fmt.Sprintf(`Create a quiz with EXACTLY %d questions based on the following content. 

Content:
%s

Requirements:
- Title: %s
- Description: %s  
- Generate EXACTLY %d questions - NO MORE, NO LESS
- Count your questions carefully and ensure you generate precisely %d questions
- ONLY generate multiple-choice questions with 4 options each
- Each question must have exactly 4 answer options
- Indicate the correct answer as index (0-3)
- DO NOT include fill-in-the-blank or incomplete questions
- LANGUAGE: Generate ALL questions and options in Bahasa Indonesia ONLY
- Keep language consistent across all questions and answer options

QUALITY GUIDELINES:
- Make questions clear, specific, and directly related to the content
- Ensure all 4 options are plausible but only one is correct
- Create distractors (wrong answers) that are reasonable but clearly incorrect
- Avoid obvious patterns (e.g., correct answer always being option A)
- Vary question difficulty and types (definition, application, analysis)
- Each question should test different concepts from the material
- Write concise explanations that clarify why the answer is correct
- Ensure questions are unambiguous and have only one correct answer

Format as JSON ARRAY with this EXACT structure:
[
  {
    "question": "What is the complete question text here?",
    "options": ["Option A text", "Option B text", "Option C text", "Option D text"],
    "correctAnswer": 0,
    "explanation": "Brief explanation"
  }
]

CRITICAL REQUIREMENTS: 
- You MUST generate exactly %d questions in the JSON array
- Question text must be complete sentences, not fill-in-the-blank format
- Do not use underscores (____) in questions
- Each question must have exactly 4 distinct options
- correctAnswer must be 0, 1, 2, or 3 (array index)
- Return ONLY a JSON array, no additional text or wrapper object

Return ONLY valid JSON array, no markdown formatting.`,
		req.QuestionCount, content, req.Title, req.Description, req.QuestionCount, req.QuestionCount, req.QuestionCount)

	return prompt
}

func (ai *AIService) createPrompt(content string, req models.CreateQuizRequest) string {
	// Truncate content if too long
	maxContentLength := 3000
	if len(content) > maxContentLength {
		content = content[:maxContentLength] + "..."
	}

	difficultyInstructions := map[string]string{
		"easy":   "Create simple, straightforward questions that test basic understanding.",
		"medium": "Create moderately challenging questions that require analysis and comprehension.",
		"hard":   "Create complex questions that require deep understanding and critical thinking.",
	}

	instruction := difficultyInstructions[req.Difficulty]
	if instruction == "" {
		instruction = difficultyInstructions["medium"]
	}

	return fmt.Sprintf(`Based on the following content, generate %d multiple choice questions. %s

Content:
%s

Requirements:
- Generate exactly %d questions
- Each question should have 4 options (A, B, C, D)
- Indicate the correct answer (0-3 index)
- Provide a brief explanation for each answer
- Return the response as a JSON array of questions

JSON Format:
[
  {
    "question": "Question text here?",
    "options": ["Option A", "Option B", "Option C", "Option D"],
    "correctAnswer": 0,
    "explanation": "Brief explanation of why this is correct"
  }
]

Return ONLY the JSON array, no additional text.`,
		req.QuestionCount, instruction, content, req.QuestionCount)
}

func (ai *AIService) parseAIResponse(response string) ([]models.Question, error) {
	// Clean up the response - remove any markdown formatting
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Try to fix truncated JSON - find the last complete question object
	if !strings.HasSuffix(response, "]") {
		ai.logger.Warn("Response appears truncated, attempting to recover valid questions")

		// Find the last complete question by looking for patterns
		// A complete question ends with: }\n  },\n  { or }\n  }\n]
		lastCompletePattern := strings.LastIndex(response, "},")
		if lastCompletePattern > 0 {
			// Found a complete object followed by comma - truncate there and close array
			response = response[:lastCompletePattern+1] + "\n]"
			ai.logger.Infof("Recovered truncated JSON by finding last complete object (pattern: }), new length: %d", len(response))
		} else {
			// Fallback: look for last complete "}" that's part of a full object
			lastBrace := strings.LastIndex(response, "}")
			if lastBrace > 0 {
				// Check if there's a comma after it (incomplete next object)
				afterBrace := response[lastBrace:]
				if !strings.Contains(afterBrace, "{") {
					// No new object started, safe to close
					response = response[:lastBrace+1] + "\n]"
					ai.logger.Infof("Recovered truncated JSON by finding last }, new length: %d", len(response))
				}
			}
		}
	}

	// Log the cleaned response for debugging
	ai.logger.Infof("Parsing AI response (first 500 chars): %s", response[:min(500, len(response))])

	var rawQuestions []struct {
		Question      string   `json:"question"`
		Options       []string `json:"options"`
		CorrectAnswer int      `json:"correctAnswer"`
		Explanation   string   `json:"explanation"`
	}

	if err := json.Unmarshal([]byte(response), &rawQuestions); err != nil {
		ai.logger.Errorf("JSON parse error: %v\nResponse was: %s", err, response[:min(1000, len(response))])
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	ai.logger.Infof("LLM generated %d questions in total", len(rawQuestions))

	var questions []models.Question
	for i, rq := range rawQuestions {
		if len(rq.Options) != 4 {
			ai.logger.Warnf("Skipping question %d with %d options: %s", i+1, len(rq.Options), rq.Question)
			continue // Skip malformed questions
		}

		questions = append(questions, models.Question{
			ID:            uuid.New().String(),
			Text:          rq.Question, // Use Text field for database
			Question:      rq.Question, // Keep Question for backward compatibility
			Options:       rq.Options,
			CorrectAnswer: rq.CorrectAnswer,
			Explanation:   rq.Explanation,
		})

		// Limit to prevent excessive questions (max 15 to match user input)
		if len(questions) >= 15 {
			break
		}
	}

	if len(questions) == 0 {
		ai.logger.Errorf("No valid questions found. Raw response had %d items", len(rawQuestions))
		return nil, fmt.Errorf("no valid questions found in response")
	}

	ai.logger.Infof("Successfully parsed %d questions from AI response", len(questions))
	return questions, nil
}

func (ai *AIService) generateFallbackQuestions(content string, req models.CreateQuizRequest) []models.Question {
	ai.logger.Warn("Using fallback question generation")

	// Generate a limited number of sample questions as fallback
	questionCount := req.QuestionCount
	if questionCount > 5 {
		questionCount = 5 // Limit fallback questions
	}

	var questions []models.Question
	contentSnippet := content
	if len(content) > 100 {
		contentSnippet = content[:100] + "..."
	}

	for i := 0; i < questionCount; i++ {
		questions = append(questions, models.Question{
			ID:       uuid.New().String(),
			Question: fmt.Sprintf("Based on the content provided, which statement is most accurate about the topic discussed? (Question %d)", i+1),
			Options: []string{
				fmt.Sprintf("The content discusses %s", contentSnippet[:min(30, len(contentSnippet))]),
				"This topic is not covered in the material",
				"The information is outdated",
				"None of the above",
			},
			CorrectAnswer: 0,
			Explanation:   "This question is based on the content analysis of the uploaded material.",
		})
	}

	return questions
}
