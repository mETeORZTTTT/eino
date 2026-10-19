/*
 * Copyright 2026 CloudWeGo Authors
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

package compose

import (
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFileIDToolResultPreserved(t *testing.T) {
	id := "file-tool-result"
	common := schema.MessagePartCommon{FileID: &id}
	result := &schema.ToolResult{Parts: []schema.ToolOutputPart{
		{Type: schema.ToolPartTypeImage, Image: &schema.ToolOutputImage{MessagePartCommon: common}},
		{Type: schema.ToolPartTypeAudio, Audio: &schema.ToolOutputAudio{MessagePartCommon: common}},
		{Type: schema.ToolPartTypeVideo, Video: &schema.ToolOutputVideo{MessagePartCommon: common}},
		{Type: schema.ToolPartTypeFile, File: &schema.ToolOutputFile{MessagePartCommon: common}},
	}}
	expected := []*schema.FunctionToolResultContentBlock{
		{Type: schema.FunctionToolResultContentBlockTypeImage, Image: &schema.UserInputImage{FileID: id}},
		{Type: schema.FunctionToolResultContentBlockTypeAudio, Audio: &schema.UserInputAudio{FileID: id}},
		{Type: schema.FunctionToolResultContentBlockTypeVideo, Video: &schema.UserInputVideo{FileID: id}},
		{Type: schema.FunctionToolResultContentBlockTypeFile, File: &schema.UserInputFile{FileID: id}},
	}
	parts, err := result.ToMessageInputParts()
	require.NoError(t, err)
	messages := []*schema.Message{{Role: schema.Tool, ToolCallID: "call-1", ToolName: "read_media", UserInputMultiContent: parts}}
	converted := toolMessageToAgenticMessage(messages)
	require.Equal(t, "call-1", converted[0].ContentBlocks[0].FunctionToolResult.CallID)
	require.Equal(t, expected, converted[0].ContentBlocks[0].FunctionToolResult.Content)
	stream := streamToolMessageToAgenticMessage(schema.StreamReaderFromArray([][]*schema.Message{messages}))
	defer stream.Close()
	chunk, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, expected, chunk[0].ContentBlocks[0].FunctionToolResult.Content)
}
