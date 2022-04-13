package main

type options struct {
	AbsScoreLimit       int      `json:"AbsScoreLimit"`
	Ascending           bool     `json:"Ascending"`
	CaseSensitive       bool     `json:"CaseSensitive"`
	Column              int      `json:"Column"`
	FirstLast           int      `json:"FirstLast"`
	IgnorePunctuation   bool     `json:"IgnorePunctuation"`
	Meta                int      `json:"Meta"`
	NumericEquivalence  bool     `json:"NumericEquivalence"`
	RelScroreLimit      int      `json:"RelScroreLimit"`
	ResultsLimit        int      `json:"ResultsLimit"`
	ReverseLookup       bool     `json:"ReverseLookup"`
	TargetFilterStrings []string `json:"TargetFilterStrings"`
	TopScoreCount       int      `json:"TopScoreCount"`
}

type searchQuery struct {
	SearchExpression []string `json:"SearchExpression"`
	Options          options `json:"Options"`
}
