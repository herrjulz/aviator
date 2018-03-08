package validator

import (
	"errors"

	"github.com/JulzDiverse/aviator"
	"github.com/starkandwayne/goutils/ansi"
)

//Error Types: Merge-Section
type MergeCombinationError struct{ error }
type MergeWithCombinationError struct{ error }
type MergeExceptCombinationError struct{ error }
type MergeRegexpCombinationError struct{ error }

//Error Types: ForEach-Section
type ForEachCombinationError struct{ error }
type ForEachFilesCombinationError struct{ error }
type ForEachInCombinationError struct{ error }
type ForEachRegexpCombinationError struct{ error }
type ForEachWalkCombinationError struct{ error }

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateSpruce(cfg []aviator.Spruce) error {
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

func validateMergeSection(cfg []aviator.Merge) error {
	for _, merge := range cfg {
		if !isMergeEmpty(merge) {
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
	}
	return nil
}

func validateForEachSection(forEach aviator.ForEach) error {
	err := validateForEachCombination(forEach)
	if err != nil {
		return err
	}

	err = validateForEachFilesCombinations(forEach)
	if err != nil {
		return err
	}

	err = validateForEachInCombinations(forEach)
	if err != nil {
		return err
	}

	err = validateForEachRegexpCombination(forEach)
	if err != nil {
		return err
	}

	err = validateForEachWalkCombinations(forEach)
	if err != nil {
		return err
	}

	return nil
}

func validateMergeCombinations(merge aviator.Merge) error {
	if (merge.With.Files != nil) && (merge.WithIn != "" || merge.WithAllIn != "") || (merge.WithIn != "" && merge.WithAllIn != "") {
		err := errors.New(
			ansi.Sprintf("@R{INVALID SYNTAX}: 'with', 'with_in', and 'with_all_in' are discrete parameters and cannot be defined together"),
		)
		return MergeCombinationError{err}
	}
	return nil
}

func validateMergeWithCombinations(with aviator.With) error {
	if (len(with.Files) == 0 || with.Files == nil) && (with.InDir != "" || with.Skip == true) {
		err := errors.New(
			ansi.Sprintf("@R{INVALID SYNTAX}: 'with.in_dir' or 'with.skip_non_existing' can only be declared in combination with 'with.files'"),
		)
		return MergeWithCombinationError{err}

	}
	return nil
}

func validateMergeExceptCombination(merge aviator.Merge) error {
	if (len(merge.Except) > 0) && (merge.WithIn == "" && merge.WithAllIn == "") {
		err := errors.New(
			ansi.Sprintf("@R{INVALID SYNTAX}: 'merge.except' is only allowed in combination with 'merge.with_in' or 'merge.with_all_in'"),
		)
		return MergeExceptCombinationError{err}
	}
	return nil
}

func validateMergeRegexpCombination(merge aviator.Merge) error {
	if (merge.Regexp != "") && (merge.WithIn == "" && merge.WithAllIn == "") {
		err := errors.New(
			ansi.Sprintf("@R{INVALID SYNTAX}: 'merge.regexp' is only allowed in combination with 'merge.with_in' or 'merge.with_all_in'"),
		)
		return MergeRegexpCombinationError{err}
	}
	return nil
}

func validateForEachCombination(forEach aviator.ForEach) error {
	if forEach.Files != nil && forEach.In != "" {
		err := errors.New(
			ansi.Sprintf("@R{INVALID SYNTAX}: Mutually exclusive parameters declared 'for_each.in' and 'for_each.files'"),
		)
		return ForEachCombinationError{err}
	}
	return nil
}

func validateForEachFilesCombinations(forEach aviator.ForEach) error {
	if (forEach.InDir != "" || forEach.Skip == true) && forEach.Files == nil {
		err := errors.New(
			ansi.Sprintf("@R{INVALID SYNTAX}: 'for_each.in_dir' and 'for_each.skip_non_existing' can only be declared in combination with 'for_each.files'"),
		)
		return ForEachFilesCombinationError{err}
	}
	return nil
}

func validateForEachInCombinations(forEach aviator.ForEach) error {
	if ((forEach.Except != nil || len(forEach.Except) > 0) || forEach.SubDirs == true) && forEach.In == "" {
		err := errors.New(
			ansi.Sprintf("@R{INVALID SYNTAX}: 'for_each.except' and 'for_each.include_sub_dirs' can only be declared in combination with 'for_each.in'"),
		)
		return ForEachInCombinationError{err}
	}
	return nil
}

func validateForEachRegexpCombination(forEach aviator.ForEach) error {
	if (forEach.Regexp != "") && (forEach.In == "") {
		err := errors.New(
			ansi.Sprintf("@R{INVALID SYNTAX}: 'for_each.regexp' is only allowed in combination with 'for_each.in'"),
		)
		return ForEachRegexpCombinationError{err}
	}
	return nil
}

func validateForEachWalkCombinations(forEach aviator.ForEach) error {
	if (forEach.SubDirs == false) && (forEach.CopyParents == true || forEach.EnableMatching == true || forEach.ForAll != "") {
		err := errors.New(
			ansi.Sprintf("INVALID SYNTAX: 'for_each.copy_parents', 'for_each.enable_matching', 'for_each.for_all' can only be declared in combination with 'for_each.inlcude_sub_dirs'"),
		)
		return ForEachWalkCombinationError{err}
	}
	return nil
}

func isForEachEmpty(forEach aviator.ForEach) bool {
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

func isMergeEmpty(merge aviator.Merge) bool {
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

func isMergeArrayEmpty(merges []aviator.Merge) bool {
	if merges == nil {
		return true
	}
	return false
}
