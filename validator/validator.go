package validator

import (
	"errors"

	"github.com/JulzDiverse/aviator/cockpit"
)

//Error Types
type MergeCombinationError struct{ error }
type MergeWithCombinationError struct{ error }
type MergeExceptCombinationError struct{ error }
type MergeRegexpCombinationError struct{ error }

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateSpruce(cfg []cockpit.Spruce) error {
	var err error
	for _, spruce := range cfg {
		err = validateMergeSection(spruce.Merge)
	}
	return err
}

func validateMergeSection(cfg []cockpit.Merge) error {
	for _, merge := range cfg {
		err := validateMergeCombinations(merge)
		if err != nil {
			return err
		}
		err = validateMergeWithCombinations(merge.With)
		if err != nil {
			return err
		}

		err = validateMergeExceptCombination(merge)
		if err != nil {
			return err
		}

		err = validateMergeRegexpCombination(merge)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateMergeCombinations(merge cockpit.Merge) error {
	var mergeError MergeCombinationError
	if (merge.With.Files != nil) && (merge.WithIn != "" || merge.WithAllIn != "") || (merge.WithIn != "" && merge.WithAllIn != "") {
		mergeError.error = errors.New(
			"INVALID SYNTAX: 'with', 'with_in', and 'with_all_in' are discrete parameters and cannot be defined together",
		)
	}
	return mergeError.error
}

func validateMergeWithCombinations(with cockpit.With) error {
	var withError MergeWithCombinationError
	if with.Files == nil && (with.InDir != "" || with.Skip == true) {
		withError.error = errors.New(
			"INVALID SYNTAX: 'with.in_dir' or 'with.skip_non_existing' can only be declared in combination with 'with.files'",
		)
	}
	return withError.error
}

func validateMergeExceptCombination(merge cockpit.Merge) error {
	var except MergeExceptCombinationError
	if (merge.Except != nil || len(merge.Except) > 0) && (merge.WithIn == "" && merge.WithAllIn == "") {
		except.error = errors.New(
			"INVALID SYNTAX: 'merge.except' is only allowed in combination with 'merge.with_in' or 'merge.with_all_in'",
		)
	}
	return except.error
}

func validateMergeRegexpCombination(merge cockpit.Merge) error {
	var regexpErr MergeRegexpCombinationError
	if (merge.Regexp != "") && ((merge.With.Files == nil || len(merge.With.Files) == 0) && merge.WithIn == "" && merge.WithAllIn == "") {
		regexpErr.error = errors.New(
			"INVALID SYNTAX: 'merge.regexp' is only allowed in combination with 'merge.with', 'merge.with_in' or 'merge.with_all_in'",
		)
	}
	return regexpErr.error
}
