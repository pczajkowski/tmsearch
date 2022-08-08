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

func (t *TM) Header() []string {
	return []string{"Name", "Source language", "Target language", "Client", "Domain", "Subject", "Owner", "LastModified", "No. of segments"}
}

func (t *TM) ToArray() []string {
	var array []string

	array = append(array, t.FriendlyName)
	array = append(array, t.SourceLangCode)
	array = append(array, t.TargetLangCode)
	array = append(array, t.Client)
	array = append(array, t.Domain)
	array = append(array, t.Subject)
	array = append(array, t.TMOwner)
	array = append(array, t.LastModifiedDate().String())
	array = append(array, fmt.Sprintf("%d", t.NumEntries))

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