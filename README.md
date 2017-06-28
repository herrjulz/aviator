# Aviator

Aviator is a tool to merge YAML files in a convenient fashion based on a configuration file called `aviator.yml`. The merge itself utilizes Spruce for the merge and therefore enables you to use all the Spruce operators in your YAML files.

If you have to handle rather complex YAML files (for BOSH or Concourse), you just provide the flight plan (`aviator.yml`), the Aviator flies you there.

# Table of Content
<!-- TOC depthFrom:1 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Aviator](#aviator)
- [Table of Content](#table-of-content)
	- [Installation](#installation)
	- [Prereqs](#prereqs)
	- [Usage](#usage)
	- [Configure an `aviator.yml`](#configure-an-aviatoryml)
		- [Spruce Section (required)](#spruce-section-required)
			- [Properties:](#properties)
				- [Base (`string`)](#base-string)
				- [Prune (`[string]Array`)](#prune-stringarray)
			- [Merge (`Array`)](#merge-array)
			- [To (`string`)](#to-string)
			- [Read From & Write To Variables](#read-from-write-to-variables)
			- [Other properties](#other-properties)
		- [The `fly` section (Optional)](#the-fly-section-optional)
- [Development](#development)

<!-- /TOC -->

## Installation

comming soon ...

## Prereqs

Aviator does not require any further prereqs, except you want to use `aviator` to _automagically_ set your concourse pipeline (for more information see [Concourse Section](#Concourse_Section)).

- For more information about [CLICK HERE](https://github.com/concourse/fly)

## Usage

To run Aviator navigate to a directory that contains an `aviator.yml` and run:

```
$ aviator
```

That's it! :)

## Configure an `aviator.yml`

Aviator provides a verbose style of configuration. You can read it like a sentence. For example the sentence _"Take the `base.yml`, merge it with `top.yml` and save it to `result.yml`"_ looks as follows in an `aviator.yml`:

```
spruce:
- base: path/to/base.yml
  merge:
  - with:
    files:
    - top.yml
  to: result.yml
```

### Spruce Section (required)

The `spruce` section is an array that defines the "plan" how YAML files should be merged. You can defines an arbitrary amount of spruce steps in this section.

#### Properties:

##### Base (`string`)

The `base` property specifies the path to the base YAML file. All other YAML files will be merged on top of the YAML files specified in this property.


##### Prune (`[string]Array`)

`prune` defines YAML properties which will be pruned during the merge.

Example:

```
spruce:
- base: base.yml
  prune:
  - meta
  - properties
  merge:
  - with:
      files:
      - top.yml
  to: result.yml
```

In this case `meta` and `properties` will be pruned during merge.

#### Merge (`Array`)

You can configure two types of objects in a `merge`: `with` and `with_in`.

**with**

With `with` you can specify specific files you want to include into the merge.

- `files` (required): List of paths to YAML files

- `in_dir` (optional): If all of the files you want to include into the merge are in one specific directory, you can specify this directory in this property and only list the file names in the `files` list. _Note: Do not forget to add the trailing "/" when specifing a path_

- `skip_non_existing` (optional): Setting this property to `true` (default `false`) will skip non existing files that are specified in the `files` list rather then returning an error.

Example:

```
spruce:
- base: path/to/base.yml
  merge:
  - with:
    files:
    - top.yml
    - top2.yml
    - top3.yml
    in_dir: path/to/
    skip_non_existing: true
  to: result.yml
```

**with_in**

`with_in` specifies a path (do not forget the trailing "/") to a directory. All files  within this directory (but not subdirectories) will be included in the merge.

Example:

```
spruce:
- base: path/to/base.yml
  merge:
  - with_in: path/to/dir/
  to: result.yml
```

**regexp**

Only files matching the regular expression will be included in the merge. It affects both `with` and `with_in`. This could be required if the target directory contains other then only YAML files.

Example:

```
spruce:
- base: path/to/base.yml
  merge:
  - with_in: path/to/dir/
  - with:
      files:
      - top.yml
  regexp: ".*.(yml)"
  to: result.yml
```

#### To (`string`)

`to` specifies the target file, where the merged files should be saved to.

#### Read From & Write To Variables

Sometimes it is required to do more than one merge step, which creates intermediate YAML files. In this case you can save merge results to variables which are defined in double courly braces `{{var}}`. You can read from & write to such a variable.

Example:

```
spruce:
- base: path/to/base.yml
  merge:
  - with_in: path/to/dir/
  to: {{result}}

- base: {{result}}
  merge:
  - with_in: another/path/
  to: final.yml
```

#### Other properties

- **for_each** OR **for_each_in** OR **walk_through**:

  - **for_each [array] (optional)** lists files, which will be merged with the `base` YAML file seperately.

  - **for_each_in [string] (optional)** (similar to `for_each`) specifies a direcotry, where each file within this direcotry will be merged with the `base` YAML.

  - **walk_through [string] (optional):** Same principle as `for_each_in`, with the difference that it walks through all subdirectories.

    - **for_all [string] (optional)** will merge all files, which was merged with `walk_through` with all files within the directory speciefied with `for_all`.

    Example (pseudo-code):

    ```
    for $element_1 in $for_all do
      for $element_2 in $walk_through do
        merge $base with $element_2, $element_1 to $to
      od
    od
    ```

    - **enable_matching [bool] (optional)** this will spruce only those files, which contain the same substring.
    - **copy_parents [bool] (optional)** copies parent directories of each file to the destination specified with `to`.


- **regexp [string] (optional):** will include only files matching the regexp.


- **to_dir [string]** is the path you want to save the spruced files to. Use this property only in combination with `for_each`, `for_each_in`, and `walk_through`.


### The `fly` section (Optional)

- **name**: Name of the pipeline
- **target**: Target short name (`fly` target)
- **config (string):** the pipeline config file (yml)
- **vars (array):** List of all property files (-l)

# Development

```
$ go get github.com/JulzDiverse/aviator
```

Navigate to `aviator` directory

```
$ godep save
```
