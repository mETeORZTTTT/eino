/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package schema

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileIDMessageRoundTrip(t *testing.T) {
	id := "file-{opaque}"
	common := MessagePartCommon{FileID: &id}
	parts := []MessageInputPart{
		{Type: ChatMessagePartTypeImageURL, Image: &MessageInputImage{MessagePartCommon: common}},
		{Type: ChatMessagePartTypeAudioURL, Audio: &MessageInputAudio{MessagePartCommon: common}},
		{Type: ChatMessagePartTypeVideoURL, Video: &MessageInputVideo{MessagePartCommon: common}},
		{Type: ChatMessagePartTypeFileURL, File: &MessageInputFile{MessagePartCommon: common}},
	}
	original := &Message{Role: User, Content: "inspect {subject}", UserInputMultiContent: parts}
	raw, err := json.Marshal(original)
	require.NoError(t, err)
	var restored Message
	require.NoError(t, json.Unmarshal(raw, &restored))
	require.Equal(t, original, &restored)
	formatted, err := restored.Format(context.Background(), map[string]any{"subject": "attachment", "opaque": "must-not-replace"}, FString)
	require.NoError(t, err)
	require.Equal(t, "inspect attachment", formatted[0].Content)
	require.Equal(t, parts, formatted[0].UserInputMultiContent)
	require.Contains(t, formatted[0].String(), "file_id=file-{opaque}")
	joined, err := ConcatMessages([]*Message{formatted[0], {Role: User, Content: " please"}})
	require.NoError(t, err)
	require.Equal(t, parts, joined.UserInputMultiContent)
}

func TestFileIDAgenticMessageRoundTrip(t *testing.T) {
	original := &AgenticMessage{Role: AgenticRoleTypeUser, ContentBlocks: []*ContentBlock{
		NewContentBlock(&UserInputText{Text: "inspect {subject}"}),
		NewContentBlock(&UserInputImage{FileID: "image-{opaque}"}),
		NewContentBlock(&UserInputAudio{FileID: "audio-{opaque}"}),
		NewContentBlock(&UserInputVideo{FileID: "video-{opaque}"}),
		NewContentBlock(&UserInputFile{FileID: "document-{opaque}", Name: "report.pdf"}),
	}}
	raw, err := json.Marshal(original)
	require.NoError(t, err)
	var restored AgenticMessage
	require.NoError(t, json.Unmarshal(raw, &restored))
	require.Equal(t, original, &restored)
	formatted, err := restored.Format(context.Background(), map[string]any{"subject": "attachment", "opaque": "must-not-replace"}, FString)
	require.NoError(t, err)
	require.Equal(t, original.ContentBlocks[1:], formatted[0].ContentBlocks[1:])
	require.Contains(t, formatted[0].String(), "file_id: video-{opaque}")
	joined, err := ConcatAgenticMessages([]*AgenticMessage{formatted[0], {Role: AgenticRoleTypeUser, ContentBlocks: []*ContentBlock{NewContentBlock(&UserInputText{Text: "please"})}}})
	require.NoError(t, err)
	require.Equal(t, original.ContentBlocks[1:], joined.ContentBlocks[1:5])
}

func TestFileIDToolResultConversion(t *testing.T) {
	id := "file-tool-result"
	common := MessagePartCommon{FileID: &id}
	parts := []ToolOutputPart{
		{Type: ToolPartTypeImage, Image: &ToolOutputImage{MessagePartCommon: common}},
		{Type: ToolPartTypeAudio, Audio: &ToolOutputAudio{MessagePartCommon: common}},
		{Type: ToolPartTypeVideo, Video: &ToolOutputVideo{MessagePartCommon: common}},
		{Type: ToolPartTypeFile, File: &ToolOutputFile{MessagePartCommon: common}},
	}
	for _, p := range parts {
		converted, err := convToolOutputPartToMessageInputPart(p)
		require.NoError(t, err)
		raw, err := json.Marshal(converted)
		require.NoError(t, err)
		require.Contains(t, string(raw), `"file_id":"file-tool-result"`)
		require.NotContains(t, string(raw), `"url"`)
	}
	agentic := &FunctionToolResultContentBlock{Type: FunctionToolResultContentBlockTypeImage, Image: &UserInputImage{FileID: id}}
	raw, err := json.Marshal(agentic)
	require.NoError(t, err)
	var restored FunctionToolResultContentBlock
	require.NoError(t, json.Unmarshal(raw, &restored))
	require.Equal(t, agentic, &restored)
}

func TestFileIDOmittedForExistingSources(t *testing.T) {
	url := "gs://bucket/image.png"
	for _, v := range []any{MessagePartCommon{URL: &url}, UserInputImage{URL: url}} {
		raw, err := json.Marshal(v)
		require.NoError(t, err)
		require.NotContains(t, string(raw), "file_id")
		require.Contains(t, string(raw), url)
	}
}
