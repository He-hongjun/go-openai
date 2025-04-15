package openai

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// Chat message role defined by the OpenAI API.
const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
	ChatMessageRoleFunction  = "function"
	ChatMessageRoleTool      = "tool"
	ChatMessageRoleDeveloper = "developer"
)

const chatCompletionsSuffix = "/chat/completions"

var (
	ErrChatCompletionInvalidModel       = errors.New("this model is not supported with this method, please use CreateCompletion client method instead") //nolint:lll
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream")              //nolint:lll
	ErrContentFieldsMisused             = errors.New("can't use both Content and MultiContent properties simultaneously")
)

type Hate struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}
type SelfHarm struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}
type Sexual struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}
type Violence struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}

type JailBreak struct {
	Filtered bool `json:"filtered"`
	Detected bool `json:"detected"`
}

type Profanity struct {
	Filtered bool `json:"filtered"`
	Detected bool `json:"detected"`
}

type ContentFilterResults struct {
	Hate      Hate      `json:"hate,omitempty"`
	SelfHarm  SelfHarm  `json:"self_harm,omitempty"`
	Sexual    Sexual    `json:"sexual,omitempty"`
	Violence  Violence  `json:"violence,omitempty"`
	JailBreak JailBreak `json:"jailbreak,omitempty"`
	Profanity Profanity `json:"profanity,omitempty"`
}

type PromptAnnotation struct {
	PromptIndex          int                  `json:"prompt_index,omitempty"`
	ContentFilterResults ContentFilterResults `json:"content_filter_results,omitempty"`
}

type ImageURLDetail string

const (
	ImageURLDetailHigh ImageURLDetail = "high"
	ImageURLDetailLow  ImageURLDetail = "low"
	ImageURLDetailAuto ImageURLDetail = "auto"
)

type ChatMessageImageURL struct {
	URL    string         `json:"url,omitempty"`
	Detail ImageURLDetail `json:"detail,omitempty"`
}

type AudioVoice string

const (
	AudioVoiceAlloy   AudioVoice = "alloy"
	AudioVoiceAsh     AudioVoice = "ash"
	AudioVoiceBallad  AudioVoice = "ballad"
	AudioVoiceCoral   AudioVoice = "coral"
	AudioVoiceEcho    AudioVoice = "echo"
	AudioVoiceSage    AudioVoice = "sage"
	AudioVoiceShimmer AudioVoice = "shimmer"
	AudioVoiceVerse   AudioVoice = "verse"
)

type AudioFormat string

const (
	AudioFormatWAV   AudioFormat = "wav"
	AudioFormatMP3   AudioFormat = "mp3"
	AudioFormatFLAC  AudioFormat = "flac"
	AudioFormatOPUS  AudioFormat = "opus"
	AudioFormatPCM16 AudioFormat = "pcm16"
)

type ChatMessageAudio struct {
	// Base64 encoded audio data.
	Data string `json:"data,omitempty"`
	// The format of the encoded audio data. Currently supports "wav" and "mp3".
	Format AudioFormat `json:"format,omitempty"`
}

type Modality string

const (
	ModalityAudio Modality = "audio"
	ModalityText  Modality = "text"
	ModalityImage Modality = "image"
)

func IsMultiOutPut(modalities []Modality) bool {
	for _, modality := range modalities {
		if modality == ModalityAudio || modality == ModalityImage {
			return true
		}
	}
	return false
}

type AudioOutput struct {
	// The voice the model uses to respond. Supported voices are alloy, ash, ballad, coral, echo, sage, shimmer, and verse.
	Voice AudioVoice `json:"voice"`
	// Specifies the output audio format. Must be one of wav, mp3, flac, opus, or pcm16.
	Format AudioFormat `json:"format"`
}

type ChatMessageFile struct {
	ID   string `json:"file_id,omitempty"`
	Name string `json:"filename,omitempty"`
	Data string `json:"file_data,omitempty"`
}

type ChatMessagePartType string

const (
	ChatMessagePartTypeText       ChatMessagePartType = "text"
	ChatMessagePartTypeImageURL   ChatMessagePartType = "image_url"
	ChatMessagePartTypeInputAudio ChatMessagePartType = "input_audio"
	ChatMessagePartTypeAudio      ChatMessagePartType = "audio"
	ChatMessagePartTypeFile       ChatMessagePartType = "file"
)

type ChatMessagePart struct {
	Type       ChatMessagePartType  `json:"type,omitempty"`
	Text       string               `json:"text,omitempty"`
	ImageURL   *ChatMessageImageURL `json:"image_url,omitempty"`
	InputAudio *ChatMessageAudio    `json:"input_audio,omitempty"`
	File       *ChatMessageFile     `json:"file,omitempty"`
}

type ChatCompletionMessage struct {
	Role         string `json:"role"`
	Content      string `json:"content,omitempty"`
	Refusal      string `json:"refusal,omitempty"`
	MultiContent []ChatMessagePart

	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name,omitempty"`

	FunctionCall *FunctionCall `json:"function_call,omitempty"`

	// For Role=assistant prompts this may be set to the tool calls generated by the model, such as function calls.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// For Role=tool prompts this should be set to the ID given in the assistant's prior request to call a tool.
	ToolCallID       string               `json:"tool_call_id,omitempty"`
	ReasoningContent string               `json:"reasoning_content,omitempty"` // only available in some reasoning models like deepseek-r1
	Audio            *ChatCompletionAudio `json:"audio,omitempty"`             // 音频内容
	Image            *ChatCompletionImage `json:"image,omitempty"`             // 图片内容
}

type chatCompletionMessageMultiContent struct {
	Role             string               `json:"role"`
	Content          string               `json:"-"`
	Refusal          string               `json:"refusal,omitempty"`
	MultiContent     []ChatMessagePart    `json:"content,omitempty"`
	Name             string               `json:"name,omitempty"`
	FunctionCall     *FunctionCall        `json:"function_call,omitempty"`
	ToolCalls        []ToolCall           `json:"tool_calls,omitempty"`
	ToolCallID       string               `json:"tool_call_id,omitempty"`
	ReasoningContent string               `json:"reasoning_content,omitempty"`
	Audio            *ChatCompletionAudio `json:"audio,omitempty"`
	Image            *ChatCompletionImage `json:"image,omitempty"`
}

type chatCompletionMessageSingleContent struct {
	Role             string               `json:"role"`
	Content          string               `json:"content,omitempty"`
	Refusal          string               `json:"refusal,omitempty"`
	MultiContent     []ChatMessagePart    `json:"-"`
	Name             string               `json:"name,omitempty"`
	FunctionCall     *FunctionCall        `json:"function_call,omitempty"`
	ToolCalls        []ToolCall           `json:"tool_calls,omitempty"`
	ToolCallID       string               `json:"tool_call_id,omitempty"`
	ReasoningContent string               `json:"reasoning_content,omitempty"`
	Audio            *ChatCompletionAudio `json:"audio,omitempty"`
	Image            *ChatCompletionImage `json:"image,omitempty"`
}

func (m ChatCompletionMessage) MarshalJSON() ([]byte, error) {
	if m.Content != "" && m.MultiContent != nil {
		return nil, errors.New("can't use both Content and MultiContent properties simultaneously")
	}
	if len(m.MultiContent) > 0 {
		msg := chatCompletionMessageMultiContent(m)
		return json.Marshal(msg)
	}

	msg := chatCompletionMessageSingleContent(m)
	return json.Marshal(msg)
}

func (m *ChatCompletionMessage) UnmarshalJSON(bs []byte) error {
	msg := chatCompletionMessageSingleContent{}

	if err := json.Unmarshal(bs, &msg); err == nil {
		*m = ChatCompletionMessage(msg)
		return nil
	}
	multiMsg := chatCompletionMessageMultiContent{}
	if err := json.Unmarshal(bs, &multiMsg); err != nil {
		return err
	}
	*m = ChatCompletionMessage(multiMsg)
	return nil
}

type ToolCall struct {
	// Index is not nil only in chat completion chunk object
	Index    *int         `json:"index,omitempty"`
	ID       string       `json:"id,omitempty"`
	Type     ToolType     `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name string `json:"name,omitempty"`
	// call function with arguments in JSON format
	Arguments string `json:"arguments,omitempty"`
}

type ChatCompletionAudio struct {
	ID         string `json:"id,omitempty"`         // 音频id
	Data       string `json:"data,omitempty"`       // 音频base6数据
	ExpiresAt  int64  `json:"expires_at,omitempty"` // 音频过期时间
	Transcript string `json:"transcript,omitempty"` // 音频转文字
}

type ChatCompletionImage struct {
	ID          string `json:"id,omitempty"`
	Created     int64  `json:"created,omitempty"`
	URL         string `json:"url,omitempty"`
	B64Data     string `json:"b64_data,omitempty"`
	Description string `json:"description,omitempty"`
}

type ChatCompletionResponseFormatType string

const (
	ChatCompletionResponseFormatTypeJSONObject ChatCompletionResponseFormatType = "json_object"
	ChatCompletionResponseFormatTypeJSONSchema ChatCompletionResponseFormatType = "json_schema"
	ChatCompletionResponseFormatTypeText       ChatCompletionResponseFormatType = "text"
)

type ChatCompletionResponseFormat struct {
	Type       ChatCompletionResponseFormatType        `json:"type,omitempty"`
	JSONSchema *ChatCompletionResponseFormatJSONSchema `json:"json_schema,omitempty"`
}

type ChatCompletionResponseFormatJSONSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Schema      json.Marshaler `json:"schema"`
	Strict      bool           `json:"strict"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
	// MaxTokens The maximum number of tokens that can be generated in the chat completion.
	// This value can be used to control costs for text generated via API.
	// This value is now deprecated in favor of max_completion_tokens, and is not compatible with o1 series models.
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat-create-max_tokens
	MaxTokens int `json:"max_tokens,omitempty"`
	// MaxCompletionTokens An upper bound for the number of tokens that can be generated for a completion,
	// including visible output tokens and reasoning tokens https://platform.openai.com/docs/guides/reasoning
	MaxCompletionTokens int                           `json:"max_completion_tokens,omitempty"`
	Temperature         float32                       `json:"temperature,omitempty"`
	TopP                float32                       `json:"top_p,omitempty"`
	N                   int                           `json:"n,omitempty"`
	Stream              bool                          `json:"stream,omitempty"`
	Stop                []string                      `json:"stop,omitempty"`
	PresencePenalty     float32                       `json:"presence_penalty,omitempty"`
	ResponseFormat      *ChatCompletionResponseFormat `json:"response_format,omitempty"`
	Seed                *int                          `json:"seed,omitempty"`
	FrequencyPenalty    float32                       `json:"frequency_penalty,omitempty"`
	// LogitBias is must be a token id string (specified by their token ID in the tokenizer), not a word string.
	// incorrect: `"logit_bias":{"You": 6}`, correct: `"logit_bias":{"1639": 6}`
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat/create-logit_bias
	LogitBias map[string]int `json:"logit_bias,omitempty"`
	// LogProbs indicates whether to return log probabilities of the output tokens or not.
	// If true, returns the log probabilities of each output token returned in the content of message.
	// This option is currently not available on the gpt-4-vision-preview model.
	LogProbs bool `json:"logprobs,omitempty"`
	// TopLogProbs is an integer between 0 and 5 specifying the number of most likely tokens to return at each
	// token position, each with an associated log probability.
	// logprobs must be set to true if this parameter is used.
	TopLogProbs int    `json:"top_logprobs,omitempty"`
	User        string `json:"user,omitempty"`
	// Deprecated: use Tools instead.
	Functions []FunctionDefinition `json:"functions,omitempty"`
	// Deprecated: use ToolChoice instead.
	FunctionCall any    `json:"function_call,omitempty"`
	Tools        []Tool `json:"tools,omitempty"`
	// This can be either a string or an ToolChoice object.
	ToolChoice any `json:"tool_choice,omitempty"`
	// Options for streaming response. Only set this when you set stream: true.
	StreamOptions *StreamOptions `json:"stream_options,omitempty"`
	// Disable the default behavior of parallel tool calls by setting it: false.
	ParallelToolCalls any `json:"parallel_tool_calls,omitempty"`
	// Store can be set to true to store the output of this completion request for use in distillations and evals.
	// https://platform.openai.com/docs/api-reference/chat/create#chat-create-store
	Store bool `json:"store,omitempty"`
	// Controls effort on reasoning for reasoning models. It can be set to "low", "medium", or "high".
	ReasoningEffort string `json:"reasoning_effort,omitempty"`
	// Metadata to store with the completion.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Output types that you would like the model to generate for this request.
	// Most models are capable of generating text, which is the default: ["text"]
	// The gpt-4o-audio-preview model can also be used to generate audio.
	// To request that this model generate both text and audio responses, you can use: ["text", "audio"]
	Modalities []Modality `json:"modalities,omitempty"`
	// Parameters for audio output. Required when audio output is requested with modalities: ["audio"]
	Audio *AudioOutput `json:"audio,omitempty"`
}

type StreamOptions struct {
	// If set, an additional chunk will be streamed before the data: [DONE] message.
	// The usage field on this chunk shows the token usage statistics for the entire request,
	// and the choices field will always be an empty array.
	// All other chunks will also include a usage field, but with a null value.
	IncludeUsage bool `json:"include_usage,omitempty"`
}

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type Tool struct {
	Type     ToolType            `json:"type"`
	Function *FunctionDefinition `json:"function,omitempty"`
}

type ToolChoice struct {
	Type     ToolType     `json:"type"`
	Function ToolFunction `json:"function,omitempty"`
}

type ToolFunction struct {
	Name string `json:"name"`
}

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
	// Parameters is an object describing the function.
	// You can pass json.RawMessage to describe the schema,
	// or you can pass in a struct which serializes to the proper JSON schema.
	// The jsonschema package is provided for convenience, but you should
	// consider another specialized library if you require more complex schemas.
	Parameters any `json:"parameters"`
}

// Deprecated: use FunctionDefinition instead.
type FunctionDefine = FunctionDefinition

type TopLogProbs struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"`
}

// LogProb represents the probability information for a token.
type LogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"` // Omitting the field if it is null
	// TopLogProbs is a list of the most likely tokens and their log probability, at this token position.
	// In rare cases, there may be fewer than the number of requested top_logprobs returned.
	TopLogProbs []TopLogProbs `json:"top_logprobs"`
}

// LogProbs is the top-level structure containing the log probability information.
type LogProbs struct {
	// Content is a list of message content tokens with log probability information.
	Content []LogProb `json:"content"`
}

type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"
	FinishReasonLength        FinishReason = "length"
	FinishReasonFunctionCall  FinishReason = "function_call"
	FinishReasonToolCalls     FinishReason = "tool_calls"
	FinishReasonContentFilter FinishReason = "content_filter"
	FinishReasonNull          FinishReason = "null"
)

func (r FinishReason) MarshalJSON() ([]byte, error) {
	if r == FinishReasonNull || r == "" {
		return []byte("null"), nil
	}
	return []byte(`"` + string(r) + `"`), nil // best effort to not break future API changes
}

type ChatCompletionChoice struct {
	Index   int                   `json:"index"`
	Message ChatCompletionMessage `json:"message"`
	// FinishReason
	// stop: API returned complete message,
	// or a message terminated by one of the stop sequences provided via the stop parameter
	// length: Incomplete model output due to max_tokens parameter or token limit
	// function_call: The model decided to call a function
	// content_filter: Omitted content due to a flag from our content filters
	// null: API response still in progress or incomplete
	FinishReason         FinishReason         `json:"finish_reason"`
	LogProbs             *LogProbs            `json:"logprobs,omitempty"`
	ContentFilterResults ContentFilterResults `json:"content_filter_results,omitempty"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID                  string                 `json:"id"`
	Object              string                 `json:"object"`
	Created             int64                  `json:"created"`
	Model               string                 `json:"model"`
	Choices             []ChatCompletionChoice `json:"choices"`
	Usage               Usage                  `json:"usage"`
	SystemFingerprint   string                 `json:"system_fingerprint"`
	PromptFilterResults []PromptFilterResult   `json:"prompt_filter_results,omitempty"`

	httpHeader
}

// CreateChatCompletion — API call to Create a completion for the chat message.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response ChatCompletionResponse, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}

	urlSuffix := chatCompletionsSuffix
	if !checkEndpointSupportsModel(urlSuffix, request.Model) {
		err = ErrChatCompletionInvalidModel
		return
	}

	reasoningValidator := NewReasoningValidator()
	if err = reasoningValidator.Validate(request); err != nil {
		return
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix, withModel(request.Model)),
		withBody(request),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
