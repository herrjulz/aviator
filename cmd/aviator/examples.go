package main

var mergeCombination = `...
  merge:
  - with:
     files:
     - file.yml
  - with_in: path/to/dir/
  - with_all_in: another/dir/
`

var forEachCombination = `...
  for_each:           for_each:
	  files:              in: path/to/dir/
		- file.yml
`

var withCombination = `...
  merge:
  - with:
      files:
      - file.yml
      skip_non_existing: true
      in_dir: path/to/dir/
`

var withInCombination = `...
  merge:
  - with_in: path/to/dir/
    except:
    - file.yml
    regexp: ".*.(yml)"
`

var withAllInCombination = `...
  merge:
  - with_all_in: path/to/dir/
    except:
    - file.yml
    regexp: ".*.(yml)"
`

var forEachFilesCombination = `...
  for_each:
    files:
    - file.yml
    in_dir: path/to/dir/
    skip_non_existing: true
`

var forEachWalkCombination = `...
  for_each:
    in: path/to/dir/
    include_sub_dirs: true
    copy_parents: true
		enable_matching: true
		except:
		- filetoexcept.yml
`
var mergeRegexpCombination = `...
  merge:
  - with_in: path/to/dir/
    regexp: ".*.(yml)"
  - with_all_in: path/to/dir/
    regexp: ".*.(yml)"
`

var mergeExceptCombination = `...
  merge:
  - with_in: path/to/dir/
    except:
		- file.yml
  - with_all_in: path/to/dir/
    except:
    - file.yml
`

var forEachRegexpCombination = `...
  for_each:
    in: path/to/dir/
		regexp: ".*.(yml)"
`
