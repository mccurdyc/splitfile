# splitfile

[![CircleCI](https://circleci.com/gh/mccurdyc/splitfile.svg?style=svg)](https://circleci.com/gh/mccurdyc/splitfile) [![codecov](https://codecov.io/gh/mccurdyc/splitfile/branch/master/graph/badge.svg)](https://codecov.io/gh/mccurdyc/splitfile) [![Maintainability](https://api.codeclimate.com/v1/badges/7f656f940224c3fe4365/maintainability)](https://codeclimate.com/github/mccurdyc/splitfile/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/7f656f940224c3fe4365/test_coverage)](https://codeclimate.com/github/mccurdyc/splitfile/test_coverage) [![Go Report Card](https://goreportcard.com/badge/github.com/mccurdyc/splitfile)](https://goreportcard.com/report/github.com/mccurdyc/splitfile)

splitfile attempts to make Go code more readable by identifying clear splits in
your package files.

## How it Works

<p align="center">
  <img width="500" height="500" src="https://github.com/mccurdyc/splitfile/blob/master/docs/imgs/splitfile.png?raw=true">
</p>

splitfile builds a graph where nodes are `type`, `var`, `const,` or `func` declarations
and undirected edges define the relationships between declarations. splitifle identifies
"splits" or the set of dependent declarations. Where there is not a clear split,
declarations will be defined as "common". Method declarations are an exception and
will not be defined as common even if they are a common node in the graph. This is
because of the way methods are declared and tied to exactly one type.

## License
+ [GNU General Public License Version 3](./LICENSE)
