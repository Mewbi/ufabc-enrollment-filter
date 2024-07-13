package server

import (
	"fmt"
	"net/url"
	"strings"

	"pdf-parser/types"
)

const (
	MIN_NAME_LENGTH = 5
	MAX_NAME_LENGTH = 50
	VALID_HOSTNAME  = "prograd.ufabc.edu.br"
)

func (ie *InputEnrollment) validate() (*types.InfoPDF, error) {
	if len(ie.Name) <= MIN_NAME_LENGTH {
		return nil, fmt.Errorf("informed name too short. Minimun size is %d characteres", MIN_NAME_LENGTH)
	}

	if len(ie.Name) >= MAX_NAME_LENGTH {
		return nil, fmt.Errorf("informed name too big. Maximum size is %d characteres", MAX_NAME_LENGTH)
	}

	pdfURL, err := url.Parse(strings.TrimSpace(ie.URL))
	if err != nil {
		return nil, fmt.Errorf("invalid URL informed: %s", err.Error())
	}

	if pdfURL.Hostname() != VALID_HOSTNAME {
		return nil, fmt.Errorf("invalid hostname informed, expected PDF from %s", VALID_HOSTNAME)
	}

	return &types.InfoPDF{
		URL:  pdfURL,
		Name: ie.Name,
	}, nil
}
