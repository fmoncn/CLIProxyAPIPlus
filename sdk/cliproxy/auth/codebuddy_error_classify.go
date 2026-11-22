package auth

import (
	"net/http"
	"strings"
)

// isContentPolicyMessage reports whether a 403 body is an upstream content
// moderation verdict on the prompt rather than a credential problem. CodeBuddy
// answers `{"code":11140,"msg":"request illegal"}` for prompts its filter
// rejects; the credential itself stays perfectly usable for other prompts.
func isContentPolicyMessage(message string) bool {
	lower := strings.ToLower(message)
	if lower == "" {
		return false
	}
	return strings.Contains(lower, "\"code\":11140") ||
		strings.Contains(lower, "\"code\": 11140") ||
		strings.Contains(lower, "request illegal")
}

// isModelNotProvisionedMessage reports whether a 400 body says the upstream has
// no service behind the requested model id. CodeBuddy answers
// `{"code":11102,"msg":"model [x] service info not found"}` after retiring a
// model; every credential answers the same way, so it is a property of the
// request, not of the credential.
func isModelNotProvisionedMessage(message string) bool {
	lower := strings.ToLower(message)
	if lower == "" {
		return false
	}
	return strings.Contains(lower, "\"code\":11102") ||
		strings.Contains(lower, "\"code\": 11102") ||
		strings.Contains(lower, "service info not found")
}

func isModelNotProvisionedResultError(err *Error) bool {
	if err == nil || statusCodeFromResult(err) != http.StatusBadRequest {
		return false
	}
	return isModelNotProvisionedMessage(err.Message)
}

// isContentPolicyResultError reports whether the failure is request-scoped
// content moderation. Such failures must not suspend the credential: doing so
// lets a single rejected prompt take every credential of the provider offline
// and starve unrelated clients with auth_unavailable.
func isContentPolicyResultError(err *Error) bool {
	if err == nil || statusCodeFromResult(err) != http.StatusForbidden {
		return false
	}
	return isContentPolicyMessage(err.Message)
}
