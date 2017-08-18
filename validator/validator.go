package validator

import (
	"errors"

	"github.com/JulzDiverse/aviator/cockpit"
)

//Error Types: Merge-Section
type MergeCombinationError struct{ error }
type MergeWithCombinationError struct{ error }
type MergeExceptCombinationError struct{ error }
type MergeRegexpCombinationError struct{ error }

//Error Types: ForEach-Section
type ForEachCombinationError error
type ForEachFilesCombinationError error

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateSpruce(cfg []cockpit.Spruce) error {
	for _, spruce := range cfg {
		if !isMergeArrayEmpty(spruce.Merge) {
			err := validateMergeSection(spruce.Merge)
			if err != nil {
				return err
			}
		}

		if !isForEachEmpty(spruce.ForEach) {
			err := validateForEachSection(spruce.ForEach)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

func validateForEachSection(forEach cockpit.ForEach) error {
	err := validateForEachCombination(forEach)
	if err != nil {
		return err
	}

	err = validateForEachFilesCombinations(forEach)
	if err != nil {
		return err
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
	if len(with.Files) == 0 && (with.InDir != "" || with.Skip == true) {
		withError.error = errors.New(
			"INVALID SYNTAX: 'with.in_dir' or 'with.skip_non_existing' can only be declared in combination with 'with.files'",
		)
	}
	return withError.error
}

func validateMergeExceptCombination(merge cockpit.Merge) error {
	var except MergeExceptCombinationError
	if (len(merge.Except) > 0) && (merge.WithIn == "" && merge.WithAllIn == "") {
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

func validateForEachCombination(forEach cockpit.ForEach) error {
	var err ForEachCombinationError
	if forEach.Files != nil && forEach.In != "" {
		err = errors.New(
			"INVALID SYNTAX: Mutually exclusive parameters declared 'for_each.in' and 'for_each.files'",
		)
	}
	return err
}

func validateForEachFilesCombinations(forEach cockpit.ForEach) error {
	var err ForEachFilesCombinationError
	if forEach.InDir != "" && forEach.Files == nil {
		err = errors.New(
			"INVALID SYNTAX: 'in_dir' can only be declared in combination with 'files'",
		)
	}
	return err
}

func isForEachEmpty(forEach cockpit.ForEach) bool {
	if (forEach.Files == nil || len(forEach.Files) == 0) &&
		forEach.InDir == "" &&
		(forEach.Except == nil || len(forEach.Except) == 0) &&
		forEach.In == "" &&
		forEach.Regexp == "" &&
		forEach.Skip == false &&
		forEach.SubDirs == false &&
		forEach.CopyParents == false &&
		forEach.EnableMatching == false &&
		forEach.ForAll == "" {
		return true
	}
	return false
}

func isMergeEmpty(merge cockpit.Merge) bool {
	if merge.With.InDir == "" &&
		merge.With.Files == nil &&
		merge.With.Skip == false &&
		merge.WithAllIn == "" &&
		merge.Except == nil &&
		merge.Regexp == "" {
		return true
	}
	return false
}

func isMergeArrayEmpty(merges []cockpit.Merge) bool {
	if merges == nil {
		return true
	}
	return false
}
