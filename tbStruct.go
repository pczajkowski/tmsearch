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

	array = append(array, t.FriendlyName)
	array = append(array, strings.Join(t.Languages, ";"))
	array = append(array, t.Client)
	array = append(array, t.Domain)
	array = append(array, t.Subject)
	array = append(array, t.TBOwner)
	array = append(array, t.LastModified.String())
	array = append(array, fmt.Sprintf("%d", t.NumEntries))

	return array
}