# splitfile

[![CircleCI](https://circleci.com/gh/mccurdyc/splitfile.svg?style=svg)](https://circleci.com/gh/mccurdyc/splitfile) [![codecov](https://codecov.io/gh/mccurdyc/splitfile/branch/master/graph/badge.svg)](https://codecov.io/gh/mccurdyc/splitfile) [![Maintainability](https://api.codeclimate.com/v1/badges/7f656f940224c3fe4365/maintainability)](https://codeclimate.com/github/mccurdyc/splitfile/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/7f656f940224c3fe4365/test_coverage)](https://codeclimate.com/github/mccurdyc/splitfile/test_coverage) [![Go Report Card](https://goreportcard.com/badge/github.com/mccurdyc/splitfile)](https://goreportcard.com/report/github.com/mccurdyc/splitfile)

splitfile aims to improve Go code readability by reducing the lines of code per file and ultimately improving the focus of files in a package.

## How it Works

splitfile identifies partitions of packages via the following method. First, a weighted directed graph is constructed, where nodes are `type`, `var`, `const,` or `func` declarations and weighted directed edges define the relationships between declarations and their uses. 
splitfile identifies partitions by assigning _configurable_ weights to edges (i.e., type -> method edges should probably be of a higher weight than type -> usage edges, but this is ultimately up to the user). Next, these weights are used by the distance function. After calculating the distance of all paths, partitions are identified by comparing the distance to a _configurable_ epsilon value. Paths with a distance greater than epsilon will be partitioned, leaving path with a distance less than epsilon "clustered".

_configurable_: edge weights and epsilon can be configured via cli flags or a `.splitfile.yml` file.

<p align="center">
  <img width="500" height="500" src="https://github.com/mccurdyc/splitfile/blob/master/docs/imgs/splitfile.png?raw=true">
</p>

## License
+ [GNU General Public License Version 3](./LICENSE)
