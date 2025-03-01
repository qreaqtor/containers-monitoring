package containersinfo

import "errors"

var (
	errDestinationUnreachable = errors.New("destination unreachable")
	errRedirect               = errors.New("should sent to another address")
	errCantReply              = errors.New("cant reply")
	errTimeout                = errors.New("ping timeout")
)
