package views

import (
	"bytes"
	"context"
	"errors"
	"ldap-password-change/cmd/config"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

type failWriter struct {
	target int
	count  int
}

func (w *failWriter) Write(p []byte) (n int, err error) {
	if w.count >= w.target {
		return 0, errors.New("simulated custom error")
	}
	w.count++
	return len(p), nil
}

func TestSuccessfulPasswordChange(t *testing.T) {
	component := SuccessfulPasswordChange()
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if doc.Find(".lead").Text() != "Your password has been successfully updated." {
		t.Errorf("expected success message not found")
	}

	if doc.Find("a[href='/']").Length() == 0 {
		t.Errorf("expected Go back link not found")
	}
}

func TestErrorToastie(t *testing.T) {
	component := ErrorToastie("Test Error Message")
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if !strings.Contains(doc.Find(".toast-body").Text(), "Test Error Message") {
		t.Errorf("expected error message not found")
	}
}

func TestIndex(t *testing.T) {
	conf := config.Config{
		Validation: config.ValidationConfig{
			UsernamePattern: "^[a-z]+$",
			PasswordPattern: "^.{8,}$",
		},
	}
	component := Index(conf)
	var buf bytes.Buffer
	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("failed to render component: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&buf)
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	if doc.Find("title").Text() != "LdapPasswordChange" {
		t.Errorf("expected title not found")
	}

	if doc.Find("form[hx-post='/change-password']").Length() == 0 {
		t.Errorf("expected form not found")
	}

	if doc.Find("input[name='username']").AttrOr("pattern", "") != conf.Validation.UsernamePattern {
		t.Errorf("expected username pattern not found")
	}
	if doc.Find("input[name='new-password']").AttrOr("pattern", "") != conf.Validation.PasswordPattern {
		t.Errorf("expected password pattern not found")
	}
}

func TestComponentsFailingWriter(t *testing.T) {
	for i := 0; i < 50; i++ {
		_ = Index(config.Config{}).Render(context.Background(), &failWriter{target: i})
		_ = SuccessfulPasswordChange().Render(context.Background(), &failWriter{target: i})
		_ = ErrorToastie("Test").Render(context.Background(), &failWriter{target: i})
	}
}
