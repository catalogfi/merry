package merry

import "fmt"

func (m *Merry) FormatVersion() string {
	return fmt.Sprintf(
		"\nVersion: %s\nCommit: %s\nDate: %s",
		m.Version, m.Commit, m.Date,
	)
}
