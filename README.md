# Aviator

Aviator is a tool to merge YAML files in a convenient fashion based on a configuration file called `aviator.yml`. The merge itself utilizes Spruce for the merge and therefore enables you to use all the Spruce operators in your YAML files.

If you have to handle rather complex YAML files (for BOSH or Concourse), you just provide the flight plan (`aviator.yml`), the Aviator flies you there.

# Table of Content

<!-- TOC depthFrom:1 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Aviator](#aviator)
- [Table of Content](#table-of-content)
	- [Installation](#installation)
		- [OS X](#os-x)
		- [Linux](#linux)
	- [Prereqs](#prereqs)
	- [Usage](#usage)
	- [Configure an `aviator.yml`](#configure-an-aviatoryml)
		- [Spruce Section (required)](#spruce-section-required)
			- [Base (`string`)](#base-string)
			- [Prune (`[string]Array`)](#prune-stringarray)
			- [Merge (`Array`)](#merge-array)
			- [To (`string`)](#to-string)
			- [Read From & Write To Variables](#read-from-write-to-variables)
			- [Environment Variables](#environment-variables)
			- [ForEach, ForEachIn & WalkThrough](#foreach-foreachin-walkthrough)
		- [The `fly` section (Optional)](#the-fly-section-optional)
- [Development](#development)

<!-- /TOC -->

## Installation

### OS X

```
$ wget -O /usr/local/bin/aviator https://github.com/JulzDiverse/aviator/releases/download/v0.2.0/aviator-darwin-amd64 && chmod +x /usr/local/bin/aviator
```

**Via Homebrew**

```
$ brew tap julzdiverse/tools  
$ brew install aviator
```

### Linux

```
$ wget -O /usr/bin/aviator https://github.com/JulzDiverse/aviator/releases/download/v0.2.0/aviator-linux-amd64 && chmod +x /usr/bin/aviator
```

## Prereqs

Aviator does not require any further prereqs, except you want to use `aviator` to _automagically_ set your concourse pipeline (for more information see [Concourse Section](#Concourse_Section)).

- For more information about [CLICK HERE](https://github.com/concourse/fly)

## Usage

To run Aviator navigate to a directory that contains an `aviator.yml` and run:

```
$ aviator
```

OR 

Specify an AVIATOR YAML FILE `.vtr` with the [--file|-f] option:

```
$ aviator -f myAviatorFile.yml
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

#### Base (`string`)

The `base` property specifies the path to the base YAML file. All other YAML files will be merged on top of the YAML files specified in this property.

---

#### Prune (`[string]Array`)

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

---

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

**with_in** (`string`)

`with_in` specifies a path (do not forget the trailing "/") to a directory. All files  within this directory (but not subdirectories) will be included in the merge.

Example:

```
spruce:
- base: path/to/base.yml
  merge:
  - with_in: path/to/dir/
  to: result.yml
```

**regexp** (`string`(quoted))

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

---

#### To (`string`)

`to` specifies the target file, where the merged files should be saved to.

---

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

#### Environment Variables

Aviator supports to read from Environment Variables. Environment variables can be set with `$VAR` or `${VAR}` at an arbitrary place in the `aviator.yml`.

Example:

```
spruce:
- base: $BASE_PATH/app-${NUMBER}.yml
  merge:
  - with_in: path/to/dir/
  to: {{result}}

- base: {{result}}
  merge:
  - with_in: $TARGET_PATH
  to: $RESULT_YAML
```

Executing `aviator` as follows:

```
$ BASE_PATH=/tmp/ NUMBER=1 RESULT_YAML=result.yml aviator
```

will resolve:

```
spruce:
- base: /tmp/app-1.yml
  merge:
  - with_in: path/to/dir/
  to: {{result}}

- base: {{result}}
  merge:
  - with_in: $TARGET_PATH
  to: result.yml
```

---

#### ForEach, ForEachIn & WalkThrough

On top of the basic `merge` the user can do more complex merges with `for_each`, `for_each_in` and `walk_through`. Note that only one of these properties can be specified per merge. For example, it is not allowed to combine `for_each` and `walk_through` in one merge step. Moreover, it requires to specify `to_dir` isntead `to` to save the merged files.  

**for_each**

`for_each` specifies a list of files that will be included in your merge seperately.

Example:

```
spruce:
- base: path/to/base.yml
  merge:
  - with:
      files:
      - top.yml
  regexp: ".*.(yml)"
  for_each:
  - env.yml
  - env2.yml
  to_dir: results/
```

This merge step will execute two merges and generate two files. It will merge `base.yml` and `top.yml` with `env.yml`, write it to `results/` and do the same with `env2.yml`.

**for_each_in**

`for_each_in` is basically the same as `for_each` with the difference that it will merge all files for a given path sperately

Example:

```
spruce:
- base: path/to/base.yml
  merge:
  - with:
      files:
      - top.yml
  regexp: ".*.(yml)"
  for_each_in: path/to/dir/
  to_dir: results/
```

**walk_through**

`walk_through` includes all files in a directory into a merge seperately. Including all subdirectories.

```
spruce:
- base: path/to/base.yml
  merge:
  - with:
      files:
      - top.yml
  regexp: ".*.(yml)"
  walk_through: path/to/dir/
  enable_matching: true
  copy_parents: true
  to_dir: results/
```

In combination with `walk_through` there are another two proprties you can define:

- `enable_matching`: this will only include files in the merge, that contains the same substring as the parent directory.

- `copy_parents`: setting this property to `true` (default `false`) will copy the parent folder of a file to the target directory (in the above example `results/`)

**regexp**

The `regexp` property can also be set in combination with `for_each`, `for_each_in`, and `walk_through` to only include files matching the regular expression.

---

### The `fly` section (Optional)

If you want to merge Concourse pipeline YAML files and set them on the fly you can specify additionally the `fly` section. If Aviator find this section it will _automagically_ execute fly for you if the following configurations are set:

- **name**: Name of the pipeline
- **target**: Target short name (`fly` target)
- **config (string):** the pipeline config file (yml)
- **vars (array):** List of all property files (-l)

Example:

```
spruce:
- base: path/to/stub.yml
  merge:
  - with_in: path/to/dir/
  to: pipeline.yml

fly:
	name: myPipelineName
	target: myFlyTarget
	config: pipeline.yml
	vars:
	- credentials.yml
```

Note, that the generated `pipeline.yml` is used in the `fly` section as `config`.

_NOTE: You will need to fly login first, before executing `aviator`_

# Development

```
$ go get github.com/JulzDiverse/aviator
```

Navigate to `aviator` directory

```
$ godep save
```
