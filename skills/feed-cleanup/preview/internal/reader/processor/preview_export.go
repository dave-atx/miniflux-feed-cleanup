// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package processor

// MinifyContent exposes the minification applied to scraped content, so
// rewrite-preview can reproduce the processor pipeline exactly.
func MinifyContent(content string) string {
	return minifyContent(content)
}
