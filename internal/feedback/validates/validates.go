package feedback_validates

import "app_burse_backend/internal/feedback"

func RaitingValidate(rating int) bool {
	return rating >= feedback.MinAvailableRaitings && rating <= feedback.MaxAvailableRaitings
}

func CommentValidate(comment string) bool {
	return len(comment) >= feedback.MinCommentLength && len(comment) <= feedback.MaxCommentLength
}

func RecommendationsValidate(recommendations *string) bool {
	return recommendations == nil || len(*recommendations) <= feedback.MaxRecomendationsLength
}
