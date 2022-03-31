package mocks

import "strings"

func MockAssetFunction(name string) ([]byte, error) {
	if strings.Contains(name, ".cy.toml") {
		return []byte("[ReleasedAfter]\none = \"Cyhoeddwyd ar ôl\"\n[ReleasedBefore]\none = \"Cyhoeddwyd cyn\"\n[DateFilterDescription]\none=\"Er enghraifft: 2006 neu 19/07/2010\""), nil
	}
	return []byte("[ReleasedAfter]\none = \"Released after\"\n[ReleasedBefore]\none = \"Released before\"\n[DateFilterDescription]\none=\"For example: 2006 or 19/07/2010\""), nil
}
