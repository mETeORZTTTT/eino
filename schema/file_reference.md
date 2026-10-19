# Uploaded file references

`MessagePartCommon.FileID` and the `UserInputImage`, `UserInputAudio`,
`UserInputVideo` and `UserInputFile` fields named `FileID` represent an opaque
identifier returned by a model service's file upload API. The value belongs to
that service and account. It is not a portable identifier across providers.

```go
fileID := "file-example"
message := &schema.Message{
    Role: schema.User,
    UserInputMultiContent: []schema.MessageInputPart{{
        Type: schema.ChatMessagePartTypeVideoURL,
        Video: &schema.MessageInputVideo{
            MessagePartCommon: schema.MessagePartCommon{FileID: &fileID},
        },
    }},
}

block := schema.NewContentBlock(&schema.UserInputVideo{FileID: fileID})
```

For each input, choose FileID, URL, or Base64Data. A selected FileID must be
nonempty and must not accompany URL or Base64Data. Model adapters are responsible
for checking this rule and encoding the reference in their API's request shape.
Existing URL/Base64 behavior is unchanged. File URIs such as `gs://...` continue
to use URL; Eino does not synthesize a URI from an identifier.

The common struct also appears in tool results and media outputs. Adding a
field to the struct does not imply that every adapter can produce or consume
that reference. Use an adapter that explicitly supports FileID for the selected
API and media type. Older adapters may reject the input or ignore unknown fields;
a schema dependency upgrade alone does not enable provider support. The deprecated
`Message.MultiContent` representation is not extended by this change.

Eino preserves the opaque value through JSON serialization, message formatting,
and tool-result conversion, including ADK/ToolsNode wrapping and the
patchtoolcalls/reduction middleware conversions. Provider support for a tool
result still depends on its API and media type. Prompt formatting does not substitute variables in
FileID. Uploading, polling file readiness, selecting account-scoped references,
checking expiry, and deleting or renewing files remain application responsibilities.

This is a schema addition. Provider conversions are developed separately in
eino-ext. The initial adapter work targets Ark Responses and OpenAI Responses;
Ark Chat, Anthropic and other APIs require their own explicit support. No new
network or filesystem operations are added to message processing.
