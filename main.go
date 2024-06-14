package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type Role string

const (
	USER      Role = "USER"
	ASSISTANT Role = "ASSISTANT"
)

type ChatMessage struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type BuildEntitiesFilterFromUserQueryTaskInput struct {
	PreviousMessages []ChatMessage `json:"previous_messages"`
	UserQuery        string        `json:"user_query"`
}

type Entity string

const (
	EntityAll        Entity = "ALL"
	EntityCalls      Entity = "CALLS"
	EntityVoiceNotes Entity = "VOICE_NOTES"
	EntityMessages   Entity = "MESSAGES"
)

type BuildEntitiesFilterFromUserQueryTaskOutput struct {
	Entities []Entity `json:"entities"`
}

type BuildEntitiesFilterFromUserQueryMetadata struct {
	EntityID string `json:"entity_id"`
	ChunkID  string `json:"chunk_id"`
}

type TaskGroupProperties struct {
	Model       string  `json:"model,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
}

type TaskGroupReference struct {
	ID         string `json:"id,omitempty"`
	Iteration  int    `json:"iteration,omitempty"`
	Properties any    `json:"properties,omitempty"`
}

type RunRequest[TInput any, TOutput any] struct {
	TaskInput TInput             `json:"task_input"`
	Group     TaskGroupReference `json:"group"`
	Labels    []string           `json:"labels,omitempty"`
	Metadata  any                `json:"metadata,omitempty"`
}

type TaskGroup struct {
	ID         string              `json:"id"`
	Iteration  int                 `json:"iteration"`
	Properties TaskGroupProperties `json:"properties"`
}

type TaskRun[TInput any, TOutput any] struct {
	TaskInput  TInput     `json:"task_input"`
	TaskOutput TOutput    `json:"task_output"`
	Group      TaskGroup  `json:"group"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	Labels     []string   `json:"labels,omitempty"`
	Metadata   any        `json:"metadata,omitempty"`
	CostUSD    *float64   `json:"cost_usd,omitempty"`
}

func BuildEntitiesFilterFromUserQuery(client *resty.Client, input BuildEntitiesFilterFromUserQueryTaskInput, entityID string, chunkID string) (*BuildEntitiesFilterFromUserQueryTaskOutput, error) {
	request := RunRequest[BuildEntitiesFilterFromUserQueryTaskInput, BuildEntitiesFilterFromUserQueryTaskOutput]{
		TaskInput: input,
		Group: TaskGroupReference{
			Iteration: 20,
		},
		Metadata: BuildEntitiesFilterFromUserQueryMetadata{
			EntityID: entityID,
			ChunkID:  chunkID,
		},
	}

	taskRun := TaskRun[BuildEntitiesFilterFromUserQueryTaskInput, BuildEntitiesFilterFromUserQueryTaskOutput]{}

	resp, err := client.R().
		SetBody(request).
		SetResult(&taskRun).
		SetPathParam("task_id", "buildentitiesfilterfromuserquery").
		SetPathParam("schema_id", "4").
		Post("/tasks/{task_id}/schemas/{schema_id}/run")

	if err != nil {
		return nil, fmt.Errorf("unexpected error from resty: %+v", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return &taskRun.TaskOutput, nil
}

func main() {
	client := resty.New().
		SetAuthScheme("Bearer").
		SetAuthToken(os.Getenv("WORKFLOWAI_TOKEN")).
		SetBaseURL("https://api.workflowai.ai")
	input := BuildEntitiesFilterFromUserQueryTaskInput{
		PreviousMessages: []ChatMessage{
			{
				Role:    USER,
				Content: "Hello !",
			},
			{
				Role:    ASSISTANT,
				Content: "Hello ! How can I help you ?",
			},
		},
		UserQuery: "Summarize my latest call with Matthias",
	}

	output, err := BuildEntitiesFilterFromUserQuery(client, input, "", "")
	if err != nil {
		panic(err)
	}
	println(output.Entities)
}
