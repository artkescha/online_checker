package writer

import (
	try "github.com/artkescha/checker/online_checker/pkg/tries"
)

type Writer interface {
	Write(*try.Try) (int64, error)
}

type Updater interface {
	Update(*try.Try) error
}

type WriterUpdater interface {
	Writer
	Updater
}
