// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package scraper

// PredefinedScraperRules returns the built-in scraper rules for the website's
// domain, or "" if there are none. Custom feed rules replace these entirely.
func PredefinedScraperRules(websiteURL string) string {
	return getPredefinedScraperRules(websiteURL)
}
