package draft

import (
	"github.com/lexiaowenn/md2wechat-new/internal/config"
	"github.com/lexiaowenn/md2wechat-new/internal/publish"
	"go.uber.org/zap"
)

// ArtifactDraftCreator adapts the draft service to the publish-layer contract.
type ArtifactDraftCreator struct {
	service *Service
}

// NewArtifactDraftCreator creates a publish-layer draft adapter.
func NewArtifactDraftCreator(cfg *config.Config, log *zap.Logger) *ArtifactDraftCreator {
	return &ArtifactDraftCreator{
		service: NewService(cfg, log),
	}
}

// CreateDraft publishes a canonical artifact as a WeChat draft.
func (c *ArtifactDraftCreator) CreateDraft(artifact publish.Artifact) (*publish.DraftResult, error) {
	result, err := c.service.CreateDraft([]Article{
		{
			Title:        artifact.Metadata.Title,
			Author:       artifact.Metadata.Author,
			Digest:       firstNonEmptyDraft(artifact.Metadata.Digest, GenerateDigestFromContent(artifact.HTML, 120)),
			Content:      artifact.HTML,
			ThumbMediaID: artifact.CoverMediaID,
			ShowCoverPic: showCoverPic(artifact.CoverMediaID),
		},
	})
	if err != nil {
		return nil, err
	}
	return &publish.DraftResult{
		MediaID:  result.MediaID,
		DraftURL: result.DraftURL,
	}, nil
}

func showCoverPic(coverMediaID string) int {
	if coverMediaID == "" {
		return 0
	}
	return 1
}

func firstNonEmptyDraft(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
