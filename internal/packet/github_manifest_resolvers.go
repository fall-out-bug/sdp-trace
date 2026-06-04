package packet

import "strings"

func checkResolvers(checks []GitHubCheck) string {
	values := []string{}
	for _, check := range checks {
		values = append(values, check.Name+"="+check.URL)
	}
	return strings.Join(values, ", ")
}

func reviewResolvers(reviews []GitHubReview) string {
	values := []string{}
	for _, review := range reviews {
		values = append(values, review.Reviewer+"="+review.Resolver)
	}
	return strings.Join(values, ", ")
}
