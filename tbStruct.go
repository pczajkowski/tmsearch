package main

import (
	"time"
	"strings"
	"fmt"
)

// TB stores information about TM.
type TB struct {
	NumEntries, AccessLevel                                                                         int
	Client, Domain, FriendlyName, Project, Subject, TBGuid, TBOwner string
	Languages []string
	LastModified time.Time
}

func (t *TB) Header() []string {
	return []string{"Name", "Languages", "Client", "Domain", "Subject", "Owner", "LastModified", "No. of entries"}
}

func (t *TB) ToArray() []string {
	var array []string

	array = append(array, t.FriendlyName,
		strings.Join(t.Languages, ";"),
		t.Client,
		t.Domain,
		t.Subject,
		t.TBOwner,
		t.LastModified.String(),
		fmt.Sprintf("%d", t.NumEntries))

	return array
}