/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package metrics

import (
	"io/fs"

	"github.com/pkg/errors"

	composeerrdefs "github.com/compose-spec/compose-go/errdefs"
)

// ComposeError error to categorize failures and extract metrics info
type ComposeError struct {
	Err      error
	Category *FailureCategory
}

// WrapComposeError wraps the error if not nil, otherwise returns nil
func WrapComposeError(err error) error {
	if err == nil {
		return nil
	}
	return ComposeError{
		Err: err,
	}
}

// WrapCategorisedComposeError wraps the error if not nil, otherwise returns nil
func WrapCategorisedComposeError(err error, failure FailureCategory) error {
	if err == nil {
		return nil
	}
	return ComposeError{
		Err:      err,
		Category: &failure,
	}
}

// Unwrap get underlying error
func (e ComposeError) Unwrap() error { return e.Err }

func (e ComposeError) Error() string { return e.Err.Error() }

// GetMetricsFailureCategory get metrics status and error code corresponding to this error
func (e ComposeError) GetMetricsFailureCategory() FailureCategory {
	if e.Category != nil {
		return *e.Category
	}
	var pathError *fs.PathError
	if errors.As(e.Err, &pathError) {
		return FileNotFoundFailure
	}
	if composeerrdefs.IsNotFoundError(e.Err) {
		return FileNotFoundFailure
	}
	return ComposeParseFailure
}
