/*
Copyright © 2025 Daniel Rivas <danielrivasmd@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

////////////////////////////////////////////////////////////////////////////////////////////////////

/*
Package horus is an error-handling library for Go that offers advanced capabilities
for generating informative, categorizable, and trackable error messages. The library
provides functions to create custom errors with detailed context—including the
operation being performed, user-friendly messages, underlying error causes, and
additional details—to facilitate debugging and logging.

Key features include:
  - Wrapping underlying errors with additional context.
  - Categorizing errors (e.g., IO, Validation, etc.) for better organization.
  - Enriching error messages with custom details, making them more descriptive.
  - Providing utilities for formatting and inspecting errors.

By using horus, developers can ensure that runtime errors carry enough information
for effective diagnosis and tracking, streamlining maintenance and troubleshooting
across complex applications.
*/

////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////
