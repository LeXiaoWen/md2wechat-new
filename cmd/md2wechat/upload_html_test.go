package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/lexiaowenn/md2wechat-new/internal/config"
	"github.com/lexiaowenn/md2wechat-new/internal/publish"
	"go.uber.org/zap"
)

func TestRunUploadHTMLUsesMetadataAndCoverUpload(t *testing.T) {
	oldCfg, oldLog := cfg, log
	oldNewDraftCreator, oldUploadCoverImageFn := newDraftCreator, uploadCoverImageFn
	oldTitle, oldAuthor, oldDigest := uploadHTMLTitle, uploadHTMLAuthor, uploadHTMLDigest
	oldCover, oldCoverMediaID, oldSource := uploadHTMLCoverImage, uploadHTMLCoverMediaID, uploadHTMLContentSource
	t.Cleanup(func() {
		cfg, log = oldCfg, oldLog
		newDraftCreator, uploadCoverImageFn = oldNewDraftCreator, oldUploadCoverImageFn
		uploadHTMLTitle, uploadHTMLAuthor, uploadHTMLDigest = oldTitle, oldAuthor, oldDigest
		uploadHTMLCoverImage, uploadHTMLCoverMediaID, uploadHTMLContentSource = oldCover, oldCoverMediaID, oldSource
	})

	cfg = &config.Config{WechatAppID: "appid", WechatSecret: "secret"}
	log = zap.NewNop()
	uploadHTMLTitle = "标题"
	uploadHTMLAuthor = "作者"
	uploadHTMLDigest = "摘要"
	uploadHTMLCoverImage = "/tmp/cover.jpg"
	uploadHTMLCoverMediaID = ""
	uploadHTMLContentSource = "https://example.com/source"

	htmlFile := filepath.Join(t.TempDir(), "article.html")
	if err := os.WriteFile(htmlFile, []byte("<p>Hello</p>"), 0600); err != nil {
		t.Fatalf("write html: %v", err)
	}

	drafter := &fakeDraftCreator{result: &publish.DraftResult{MediaID: "draft-1", DraftURL: "https://example.com/draft"}}
	newDraftCreator = func() publish.DraftCreator { return drafter }
	uploadCoverImageFn = func(imagePath string) (string, error) {
		if imagePath != "/tmp/cover.jpg" {
			t.Fatalf("cover image path = %q", imagePath)
		}
		return "cover-media-id", nil
	}

	response, err := runUploadHTML(htmlFile)
	if err != nil {
		t.Fatalf("runUploadHTML() error = %v", err)
	}
	if response["media_id"] != "draft-1" || response["cover_media_id"] != "cover-media-id" {
		t.Fatalf("response = %#v", response)
	}
	if len(drafter.artifacts) != 1 {
		t.Fatalf("drafter artifacts = %#v", drafter.artifacts)
	}
	artifact := drafter.artifacts[0]
	if artifact.HTML != "<p>Hello</p>" || artifact.CoverMediaID != "cover-media-id" {
		t.Fatalf("artifact = %#v", artifact)
	}
	if artifact.Metadata.Title != "标题" || artifact.Metadata.Author != "作者" || artifact.Metadata.Digest != "摘要" || artifact.Metadata.ContentSourceURL != "https://example.com/source" {
		t.Fatalf("metadata = %#v", artifact.Metadata)
	}
}

func TestRunUploadHTMLUsesExistingCoverMediaID(t *testing.T) {
	oldCfg, oldLog := cfg, log
	oldNewDraftCreator, oldUploadCoverImageFn := newDraftCreator, uploadCoverImageFn
	oldTitle, oldCover, oldCoverMediaID := uploadHTMLTitle, uploadHTMLCoverImage, uploadHTMLCoverMediaID
	t.Cleanup(func() {
		cfg, log = oldCfg, oldLog
		newDraftCreator, uploadCoverImageFn = oldNewDraftCreator, oldUploadCoverImageFn
		uploadHTMLTitle, uploadHTMLCoverImage, uploadHTMLCoverMediaID = oldTitle, oldCover, oldCoverMediaID
	})

	cfg = &config.Config{WechatAppID: "appid", WechatSecret: "secret"}
	log = zap.NewNop()
	uploadHTMLTitle = ""
	uploadHTMLCoverImage = ""
	uploadHTMLCoverMediaID = "existing-cover"

	htmlFile := filepath.Join(t.TempDir(), "article.html")
	if err := os.WriteFile(htmlFile, []byte("<p>Hello</p>"), 0600); err != nil {
		t.Fatalf("write html: %v", err)
	}

	drafter := &fakeDraftCreator{result: &publish.DraftResult{MediaID: "draft-2"}}
	newDraftCreator = func() publish.DraftCreator { return drafter }
	uploadCoverImageFn = func(imagePath string) (string, error) {
		t.Fatalf("uploadCoverImageFn should not be called")
		return "", nil
	}

	response, err := runUploadHTML(htmlFile)
	if err != nil {
		t.Fatalf("runUploadHTML() error = %v", err)
	}
	if response["title"] != "article" {
		t.Fatalf("expected filename title fallback, got %#v", response)
	}
	if drafter.artifacts[0].CoverMediaID != "existing-cover" {
		t.Fatalf("artifact = %#v", drafter.artifacts[0])
	}
}

func TestRunUploadHTMLRejectsInvalidCoverInputs(t *testing.T) {
	oldCfg := cfg
	oldCover, oldCoverMediaID := uploadHTMLCoverImage, uploadHTMLCoverMediaID
	t.Cleanup(func() {
		cfg = oldCfg
		uploadHTMLCoverImage, uploadHTMLCoverMediaID = oldCover, oldCoverMediaID
	})

	cfg = &config.Config{WechatAppID: "appid", WechatSecret: "secret"}
	uploadHTMLCoverImage = ""
	uploadHTMLCoverMediaID = ""

	if _, err := runUploadHTML("article.html"); err == nil {
		t.Fatal("expected missing cover error")
	}

	uploadHTMLCoverImage = "/tmp/cover.jpg"
	uploadHTMLCoverMediaID = "existing-cover"
	if _, err := runUploadHTML("article.html"); err == nil {
		t.Fatal("expected conflicting cover error")
	}
}

func TestUploadHTMLCmdOutputsStableEnvelope(t *testing.T) {
	oldCfg, oldLog := cfg, log
	oldJSON := jsonOutput
	oldNewDraftCreator, oldUploadCoverImageFn := newDraftCreator, uploadCoverImageFn
	oldTitle, oldCover, oldCoverMediaID := uploadHTMLTitle, uploadHTMLCoverImage, uploadHTMLCoverMediaID
	t.Cleanup(func() {
		cfg, log = oldCfg, oldLog
		jsonOutput = oldJSON
		newDraftCreator, uploadCoverImageFn = oldNewDraftCreator, oldUploadCoverImageFn
		uploadHTMLTitle, uploadHTMLCoverImage, uploadHTMLCoverMediaID = oldTitle, oldCover, oldCoverMediaID
		uploadHTMLCmd.SetArgs(nil)
	})

	cfg = &config.Config{WechatAppID: "appid", WechatSecret: "secret"}
	log = zap.NewNop()
	jsonOutput = true

	htmlFile := filepath.Join(t.TempDir(), "article.html")
	if err := os.WriteFile(htmlFile, []byte("<p>Hello</p>"), 0600); err != nil {
		t.Fatalf("write html: %v", err)
	}

	drafter := &fakeDraftCreator{result: &publish.DraftResult{MediaID: "draft-3"}}
	newDraftCreator = func() publish.DraftCreator { return drafter }
	uploadCoverImageFn = func(imagePath string) (string, error) {
		return "cover-media-id", nil
	}

	uploadHTMLCmd.SetArgs([]string{htmlFile, "--title", "标题", "--cover", "/tmp/cover.jpg"})
	stdout := captureStdout(t, func() {
		if err := uploadHTMLCmd.Execute(); err != nil {
			t.Fatalf("uploadHTMLCmd.Execute() error = %v", err)
		}
	})

	var response map[string]any
	if err := json.Unmarshal(stdout, &response); err != nil {
		t.Fatalf("unmarshal response: %v\n%s", err, stdout)
	}
	if response["success"] != true || response["code"] != codeUploadHTMLCreated {
		t.Fatalf("unexpected response: %#v", response)
	}
	data, _ := response["data"].(map[string]any)
	if data["media_id"] != "draft-3" || data["title"] != "标题" {
		t.Fatalf("unexpected data: %#v", data)
	}
}
