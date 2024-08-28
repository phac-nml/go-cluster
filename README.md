![coverage](https://raw.githubusercontent.com/phac-nml/go-cluster/badges/.badges/main/coverage.svg)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# go-cluster

- [Introduction](#introduction)
  * [Citation](#citation)
  * [Contact](#contact)
- [Install](#install)
    + [Compatibility](#compatibility)
- [Getting Started](#getting-started)
  * [Usage](#usage)
  * [Configuration and Settings](#configuration-and-settings)
  * [Data Input](#data-input)
- [Troubleshooting and FAQs](#troubleshooting-and-faqs)
- [Other information](#other-information)
- [Legal and Compliance Information](#legal-and-compliance-information)
- [Updates and Release Notes](#updates-and-release-notes)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

# Introduction

This program is created to supply a user with the ability to quickly create distance matrices of allelic profiles through the use of parallel processing. The program prints pairwise distances in a "molten" or flat format where distances are labeled as "sample1 sample2 distance". But a utility to quickly convert this molten format to a distance matrix is built into `go-cluster`. Created distance matrices can be used to cluster the data and output dendrograms. Querying a reference set of profiles against a larger set of profiles in parallel is supported as well.

## Contact

[Matthew Wells] : <matthew.wells@phac-aspc.gc.ca>

# Install

No install is currently provided, the program must be compiled using either the `go` toolchain or `gogcc` compiler. This program of go was built using version 1.1.8 as listed in the `go.mod` file. Required dependencies are listed in the `go.mod` file, however their installation is controlled by the `go` toolchain. Instructions for installing `go` can be found here [go](https://go.dev/doc/install), Go can also likely be installed through other packaging programs like conda.

If you have `go` installed, you can simply run `go build` in the source directory and a binary file called `go-cluster` will be generated.

### Compatibility

`go-cluster` has only been tested on linux.

# Getting Started

## Usage

### CLI Arguments for each utility

The main sub-commands for the program are listed below.
```
Parallel Distances - A program for getting distances between allelic profiles and creating distance matrices.

  Usage:
    Parallel Distances [distances|convert|tree|fast-match]

  Subcommands: 
    distances    Compute all pairwise distances between the specified input profile.
    convert      Convert the pairwise distance generated by the program into a distance matrix.
    tree         Create a dendrogram from a supplied distance matrix.
    fast-match   Tabulate distances between a query profile and reference profiles. Only distances exceeding a threshold will be kept.

  Flags: 
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.

```

#### Calculate pairwise distances - distance

```
distances - Compute all pairwise distances between the specified input profile.

  Flags: 
       --version                    Displays the program version string.
    -h --help                       Displays help with available flag, subcommand, and positional value parameters.
    -i --input                      File path to your alleles profiles.
    -l --load-factor                This value is used to compute how many profile calculations are assigned to thread, a larger value will result in fewer threads being used. Default: 2 (default: 2)
    -d --distance                   Enter an integer denoting the distance function you would like to use:
        Hamming Distance skipping missing values: 0
        Hamming distance missing values treated as alleles.: 1
        Scaled Distance skipping missing values: 2
        Scaled distance missing values treated as alleles.: 3 (default: 0)
    -o --output                     Name of output file. If nothing is specified results will be sent to stdout.
    -b --buffer-size                The default buffer size is: 16384. Larger buffers may increase performance. (default: 16384)
    -c --column-delimiter           Column delimiter, default value is a tab character (default:        )
    -m --missing-allele-character   String denoting missing alleles. (default: 0)
```

#### Convert output into a distance matrix - convert

```
convert - Convert the pairwise distance generated by the program into a distance matrix.

  Flags: 
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.
    -i --input     File path to a previously generated output for conversion into a distance matrix.
    -o --output    Name of output file. If nothing is specified results will be sent to stdout.

```

#### Create a dendrogram - tree

```
tree - Create a dendrogram from a supplied distance matrix.

  Flags: 
       --version          Displays the program version string.
    -h --help             Displays help with available flag, subcommand, and positional value parameters.
    -i --input            File path to previously generate distance matrix.
    -o --output           Name of output file. If nothing is specified results will be sent to stdout.
    -l --linkage-method   Please enter an integer corresponding to one of the linkage method of your choice: average: 0 centroid: 1 complete: 2 mcquitty: 3 median: 4 single: 5 ward: 6  (default: 0)

```

#### Query a profile(s) against another set - fast-match

```
fast-match - Tabulate distances between a query profile and reference profiles. Only distances exceeding a threshold will be kept.

  Flags: 
       --version                    Displays the program version string.
    -h --help                       Displays help with available flag, subcommand, and positional value parameters.
    -i --input                      File path to profiles for querying.
    -r --reference                  File path to reference profiles to query against.
    -c --column-delimiter           Column delimiter (default:  )
    -m --missing-allele-character   String denoting missing alleles. (default: 0)
    -d --distance                   Enter an integer denoting the distance function you would like to use:
        Hamming Distance: 0
        Hamming distance skipping missing values: 1
        Scaled Distance: 2
        Scaled distance skipping missing values: 3 (default: 0)
    -t --threshold                  Threshold for matching alleles. (default: 10.00)
    -o --output                     Name of output file. If nothing is specified results will be sent to stdout.
    -l --goroutine-limit            Limit the number of goroutines run at one time. Default: 100 (default: 100)

```

## Data Input

### File format for calculating pairwise-distances

The input file is a table with the first column containing a given samples name. While the subsequent columns are the names of the input loci. The program normalizes allele input first, so the actual content of the cells can vary and can be, integers, hashes or even nucleotides themselves (this would make the program run slower however). CSV files can be used, however you will have to specify the the delimiter to the command line arguments

|Sample|Loci_1|Loci_2|Loci_3|
|------|------|------|------|
|Sample1|Loci_x|5|1|
|Sample2|1|xxxxxxxxxxxx|1|
|Sample3|2|3|1|

### File format for creating a distance matrix

The pairwise distances output from the `distances` option is the input to the `distances` command. But if you are curious the output looks something like:

```
sample1 sample2 distance
sample1 sample3 distance
...
```

the output is a distance symmetric distance matrix.

### Clustering

The input is a symmetric distance matrix e.g.

```
    1   2   3
1   0  40  40
2   40  0   40
3   40  40  0
```

and the output is a newick file.

### Fast matching

The input for fast matching is two allelic profiles similar to what is passed when creating a distance matrix. There are some checks in the program to verify that you have the same number of columns but it is **IMPORTANT** That your alleles are in the same order for both of the compared profiles.

# Troubleshooting and FAQs

This program utilizes `go routines` for parallel processing which allows for more threads than logical CPUs. Currently you cannot set the maximum number of CPUs to use in this program as `go routines` will be distributed across CPUs at run time. 

This can be an issue on grid executors like `slurm` and it is best to run this program on an entire node of a cluster at a time. If you face this issue it is best to run the process on an entire node e.g. the `-w` flag in slurm to select a node and adding the `--exclusive` flag can solve this issue.

# Other information

This program is still in development however basic end-to-end tests have been added. More code examples, benchmarks and unit-tests are still to be added. A lot of further optimizations can still be made to this program as well.

# Legal and Compliance Information

Copyright Government of Canada [2024]

Written by: National Microbiology Laboratory, Public Health Agency of Canada

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this work except in compliance with the License. You may obtain a copy of the License at:

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

# Updates and Release Notes

Please see the `CHANGELOG.md`.
