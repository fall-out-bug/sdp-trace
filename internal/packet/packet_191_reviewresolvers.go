package packet

import (
	"strings"
)

func reviewResolvers(reviews []GitHubReview) string {
	// reviewResolvers keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	values := []string{}
	for _, review := range reviews {

		values = append(values, review.Reviewer+"="+review.Resolver)
	}
	return strings.Join(values, ", ")
}
