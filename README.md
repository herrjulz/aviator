# Aviator

Aviator is a small CLI tool to run genereic **Aviator** Concourse pipelines.

## Installation (Mac Only)

Download Aviator [here](https://ibm.box.com/s/hoz7v6x9tlx1yothmrwox90wa94ezhu8)

and run:

```
$ install ~/Downloads/aviator /usr/local/bin  
```

## Prereqs

- [Spruce](https://github.com/geofffranks/spruce) CLI Tool
- [Fly](https://github.com/concourse/fly) CLI Tool

## Usage

**aviator.yml**

```
spruce:
- base: base.yml
  prune:
  - meta
  with:
  - another.yml
  - yet-another.yml
  to: result.yml
- base: result.yml
  for_each_in: path/to/dir/
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

**Run Aviator**

```
$ aviator -t <target> -p <pipeline-name>
```
