package main

import (
	"time"
	"strings"
	"log"
	"fmt"
)

// TM stores information about TM.
type TM struct {
	NumEntries, AccessLevel                                                                         int
	Client, Domain, FriendlyName, Project, SourceLangCode, Subject, TMGuid, TMOwner, TargetLangCode string
	LastModified string
}

func (_ *TM) Header() []string {
	return []string{"Name", "Source language", "Target language", "Client", "Domain", "Subject", "Owner", "LastModified", "No. of segments"}
}

func (t *TM) ToArray() []string {
	var array []string

	array = append(array, t.FriendlyName,
		t.SourceLangCode,
		t.TargetLangCode,
		t.Client,
		t.Domain,
		t.Subject,
		t.TMOwner,
		t.LastModifiedDate().String(),
		fmt.Sprintf("%d", t.NumEntries))

	return array
}

func (t *TM) LastModifiedDate() *time.Time {
	if !strings.HasSuffix(t.LastModified, "Z") {
		t.LastModified += ".000Z"
	}

	modified, err := time.Parse(time.RFC3339, t.LastModified)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &modified
}