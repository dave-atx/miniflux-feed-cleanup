// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package rewrite

// PredefinedRewriteRules returns the built-in rewrite rules for the entry's
// domain, or "" if there are none. Custom feed rules replace these entirely.
func PredefinedRewriteRules(entryURL string) string {
	return getPredefinedRewriteRules(entryURL)
}
