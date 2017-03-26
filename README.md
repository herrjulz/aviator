# Aviator

Aviator is a CLI tool based on `spruce` to run **Aviator** YAML files. An `aviator.yml` is  a configuration file to

- merge YAML files (e.g. Bosh Manifest Files)
- generate generic `Concourse` Pipelines

## Installation

comming soon ...

## Prereqs

- [Spruce](https://github.com/geofffranks/spruce) CLI Tool
- [Fly](https://github.com/concourse/fly) CLI Tool

## The Aviator CLI

The `aviator` CLI is a command line tool to execute `aviator.yml` files. To use `aviator`, simply navigate to a directory containing an `aviator.yml` and execute it with:

```
$ aviator
```

If you use it for `Concourse` pipelines, you need to specify a target and a pipeline name:

```
$ aviator -t <target> -p <pipeline-name>
```

## Configure an Aviator YAML

Aviator YAMLs provide plans on how yaml files should be merged. Such a plan is configured in up to two sections: `spruce` (required) and `fly` (optional). The `spruce` section specifies the files and the order they need to be merged. The `fly` section specifies the fly command which needs to be executed for a specific (`concourse`) YAML file. The following code snippet shows an example of an `aviator.yml` file:  

```
spruce:
- base: base.yml
  prune:
  - meta
  with:
    files:
    - another.yml
    - yet-another.yml
  to: result.yml
- base: result.yml
  for_each_in: path/to/dir/
  regexp: match-string
  to_dir: path/to/destination/
- base: another-base.yml
  walk_through: will/walt/through/subdirs
  to_dir: path/to/destiantion/

fly:
 config: pipeline.yml
 vars:
 - credentials.yml
 - personal.yml
```

### The `spruce` section

The `spruce` section is an Array of merge steps. Each merge step merges several files to an resulting file. The properties you can use to merge files are the following:

- **base [string] (required)** specifies the base YAML. All other specified YAMLs will be merged on top of this file.

- **prune [array] (optional):** lists all properties, that needs to be pruned from the merged files.

- **with [map] (optional)** specifies either specific files from different locations or from a specific location.

    - **files [array] (required)** lists specific files you want to spruce on top of the base YAML.
    - **in_dir [string] (optional)**  specifies the location you want to pick specific files.


- **with_in [string] (optional)** picks up each file in a given directory.

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

- **to** OR **to_dir** (one is required)

  - **to [string]** specifies the filename you want to save the result to.

  - **to_dir [string]** is the path you want to save the spruced files to. Use this property only in combination with `for_each`, `for_each_in`, and `walk_through`.

### The `fly` section (Optional)

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
