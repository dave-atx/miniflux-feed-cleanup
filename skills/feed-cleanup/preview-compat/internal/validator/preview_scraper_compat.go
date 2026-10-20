// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"github.com/andybalholm/cascadia"

	"miniflux.app/v2/internal/locale"
)

// IsValidScraperRules backfills the scraper selector check for Miniflux
// releases that predate it (2.3.4 and earlier). Those servers save an invalid
// selector without complaint, so the preview is the only place it's caught.
// The message is a format string because these releases have no translation
// for it.
func IsValidScraperRules(rules string) *locale.LocalizedError {
	if rules == "" {
		return nil
	}
	if _, err := cascadia.Compile(rules); err != nil {
		return locale.NewLocalizedError("Invalid CSS selector in scraper rules: %v", err.Error())
	}
	return nil
}
