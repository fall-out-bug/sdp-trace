package prreview

import (
	"errors"
)

func promptTemplateCannotVerify(err error) bool {
	return errors.Is(err, errPromptTemplateCannotVerify)
}

func promptEvidenceCannotVerify(err error) bool {
	return errors.Is(err, errPromptEvidenceCannotVerify)
}
