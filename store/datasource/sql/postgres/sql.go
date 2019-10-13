package postgres

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

// Not implemented, yet.
var index = map[string]string{

}