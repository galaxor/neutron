package imap

import (
	"github.com/mxk/go-imap/imap"
)

func wait(cmd *imap.Command, err error) (*imap.Command, *imap.Response, error) {
	if err != nil {
		return nil, nil, err
	}

	cmd, err = imap.Wait(cmd, err)
	if err != nil {
		return cmd, nil, err
	}

	res, err := cmd.Result(imap.OK)
	return cmd, res, err
}
