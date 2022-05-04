package mocks

import "strings"

var cyLocale = []string{
	"[ReleasedAfter]",
	"one = \"Cyhoeddwyd ar ôl\"",
	"[ReleasedBefore]",
	"one = \"Cyhoeddwyd cyn\"",
	"[DateFilterDescription]",
	"one=\"Er enghraifft: 2006 neu 19/07/2010\"",
	"[ReleaseCalendarPageTitle]",
	"one=\"Calendr datganiadau\"",
	"[BreadcrumbHome]",
	"one=\"Hafan\"",
	"[BreadcrumbReleaseCalendar]",
	"one = \"Calendr datganiadau\"",
	"[BreadcrumbUpcoming]",
	"one = \"Ar ddod\"",
	"[BreadcrumbCancelled]",
	"one = \"Canslwyd\"",
}

var enLocale = []string{
	"[ReleasedAfter]",
	"one = \"Released after\"",
	"[ReleasedBefore]",
	"one = \"Released before\"",
	"[DateFilterDescription]",
	"one=\"For example: 2006 or 19/07/2010\"",
	"[ReleaseCalendarPageTitle]",
	"one=\"Release Calendar\"",
	"[BreadcrumbHome]",
	"one=\"Home\"",
	"[BreadcrumbReleaseCalendar]",
	"one = \"Release calendar\"",
	"[BreadcrumbUpcoming]",
	"one = \"Upcoming\"",
	"[BreadcrumbCancelled]",
	"one = \"Cancelled\"",
}

func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte(strings.Join(cyLocale, "\n")), nil
	}
	return []byte(strings.Join(enLocale, "\n")), nil
}
