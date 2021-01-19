package doc

// Tag represents start & end html tag
// (e.g) start: <span class="highlight">
// end: </span>
type Tag struct {
	Start string
	End   string
}

// WordPos represents word location
type WordPos struct {
	StartIdx int
	EndIdx   int
}
