package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LeXiaoWen/md2wechat-new/internal/publish"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	uploadHTMLTitle         string
	uploadHTMLAuthor        string
	uploadHTMLDigest        string
	uploadHTMLCoverImage    string
	uploadHTMLCoverMediaID  string
	uploadHTMLContentSource string
)

var uploadHTMLCmd = &cobra.Command{
	Use:   "upload_html <html_file>",
	Short: "Create a WeChat draft from an existing HTML file",
	Long: `Create a WeChat Official Account draft from an existing HTML file.

The HTML is used as-is. Use --cover to upload a local cover image, or
--cover-media-id to reuse an existing permanent WeChat cover material.`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		response, err := runUploadHTML(args[0])
		if err != nil {
			return err
		}
		responseSuccessWith(codeUploadHTMLCreated, "HTML draft created successfully", response)
		return nil
	},
}

func init() {
	uploadHTMLCmd.Flags().StringVar(&uploadHTMLTitle, "title", "", "Draft title (defaults to HTML filename)")
	uploadHTMLCmd.Flags().StringVar(&uploadHTMLAuthor, "author", "", "Draft author")
	uploadHTMLCmd.Flags().StringVar(&uploadHTMLDigest, "digest", "", "Draft digest, max 128 characters")
	uploadHTMLCmd.Flags().StringVar(&uploadHTMLCoverImage, "cover", "", "Cover image path for draft")
	uploadHTMLCmd.Flags().StringVar(&uploadHTMLCoverMediaID, "cover-media-id", "", "Existing WeChat cover media_id (mutually exclusive with --cover)")
	uploadHTMLCmd.Flags().StringVar(&uploadHTMLContentSource, "content-source-url", "", "Original content source URL")
}

func runUploadHTML(htmlFile string) (map[string]any, error) {
	if err := cfg.ValidateForWeChat(); err != nil {
		return nil, wrapCLIError(codeConfigInvalid, err, err.Error())
	}

	if strings.TrimSpace(uploadHTMLCoverImage) != "" && strings.TrimSpace(uploadHTMLCoverMediaID) != "" {
		return nil, newCLIError(codeUploadHTMLInvalid, "--cover and --cover-media-id are mutually exclusive")
	}
	if strings.TrimSpace(uploadHTMLCoverImage) == "" && strings.TrimSpace(uploadHTMLCoverMediaID) == "" {
		return nil, newCLIError(codeUploadHTMLInvalid, "--cover or --cover-media-id is required")
	}
	if strings.TrimSpace(uploadHTMLCoverMediaID) != "" && looksLikeURL(uploadHTMLCoverMediaID) {
		return nil, newCLIError(codeUploadHTMLInvalid, "--cover-media-id expects a WeChat media_id, not a URL")
	}

	html, err := os.ReadFile(htmlFile)
	if err != nil {
		return nil, wrapCLIError(codeUploadHTMLReadFailed, err, fmt.Sprintf("read HTML file: %v", err))
	}

	title := strings.TrimSpace(uploadHTMLTitle)
	if title == "" {
		title = titleFromHTMLPath(htmlFile)
	}
	author := strings.TrimSpace(uploadHTMLAuthor)
	digest := strings.TrimSpace(uploadHTMLDigest)
	if err := validateConvertMetadata(title, author, digest); err != nil {
		return nil, err
	}

	coverMediaID := strings.TrimSpace(uploadHTMLCoverMediaID)
	if coverMediaID == "" {
		log.Info("uploading HTML draft cover image", zap.String("path", uploadHTMLCoverImage))
		coverMediaID, err = uploadCoverImageFn(uploadHTMLCoverImage)
		if err != nil {
			return nil, wrapCLIError(codeUploadHTMLCoverFailed, err, fmt.Sprintf("upload cover: %v", err))
		}
	}

	svc := newDraftCreator()
	result, err := svc.CreateDraft(publish.Artifact{
		HTML: string(html),
		Metadata: publish.Metadata{
			Title:            title,
			Author:           author,
			Digest:           digest,
			ContentSourceURL: strings.TrimSpace(uploadHTMLContentSource),
		},
		CoverMediaID: coverMediaID,
	})
	if err != nil {
		return nil, wrapCLIError(codeUploadHTMLCreateFailed, err, fmt.Sprintf("create draft: %v", err))
	}

	response := map[string]any{
		"media_id":       result.MediaID,
		"cover_media_id": coverMediaID,
		"title":          title,
		"author":         author,
		"digest":         digest,
		"html_file":      htmlFile,
	}
	if result.DraftURL != "" {
		response["draft_url"] = result.DraftURL
	}
	return response, nil
}

func titleFromHTMLPath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	title := strings.TrimSuffix(base, ext)
	if strings.TrimSpace(title) == "" {
		return "未命名文章"
	}
	return title
}
