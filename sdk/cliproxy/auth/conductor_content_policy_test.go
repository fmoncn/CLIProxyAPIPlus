package auth

import (
	"net/http"
	"testing"
	"time"
)

const codebuddyBlockBody = `{"code":11140,"msg":"request illegal"}`

func TestIsContentPolicyResultError(t *testing.T) {
	cases := []struct {
		name string
		err  *Error
		want bool
	}{
		{"nil", nil, false},
		{"codebuddy 11140", &Error{HTTPStatus: http.StatusForbidden, Message: codebuddyBlockBody}, true},
		{"spaced json", &Error{HTTPStatus: http.StatusForbidden, Message: `{"code": 11140,"msg":"request illegal"}`}, true},
		{"plain forbidden", &Error{HTTPStatus: http.StatusForbidden, Message: "permission denied"}, false},
		{"same body on 429", &Error{HTTPStatus: http.StatusTooManyRequests, Message: codebuddyBlockBody}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isContentPolicyResultError(tc.err); got != tc.want {
				t.Fatalf("isContentPolicyResultError = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestApplyAuthFailureState_ContentPolicyKeepsAuthUsable(t *testing.T) {
	auth := &Auth{ID: "cb-1", Provider: "codebuddy", Status: StatusActive}
	applyAuthFailureState(auth, &Error{HTTPStatus: http.StatusForbidden, Message: codebuddyBlockBody}, nil, time.Now())

	if auth.Unavailable {
		t.Fatal("内容风控 403 不应把凭据置为 Unavailable")
	}
	if auth.Status == StatusError {
		t.Fatal("内容风控 403 不应把凭据置为 StatusError")
	}
}

func TestApplyAuthFailureState_PlainForbiddenStillSuspends(t *testing.T) {
	auth := &Auth{ID: "cb-2", Provider: "codebuddy", Status: StatusActive}
	applyAuthFailureState(auth, &Error{HTTPStatus: http.StatusForbidden, Message: "permission denied"}, nil, time.Now())

	if !auth.Unavailable {
		t.Fatal("普通 403 仍应把凭据置为 Unavailable")
	}
	if auth.Status != StatusError {
		t.Fatalf("普通 403 状态应为 StatusError，实际 %v", auth.Status)
	}
}
