package filters

import "github.com/pkg/errors"

var (
	errFilterNotFound         = errors.New("filter not found")
	errInvalidBlockRange      = errors.New("invalid block range params")
	errPendingLogsUnsupported = errors.New("pending logs are not supported")
	errExceedMaxTopics        = errors.New("exceed max topics")
)
