# splitfile

[![CircleCI](https://circleci.com/gh/mccurdyc/splitfile.svg?style=svg)](https://circleci.com/gh/mccurdyc/splitfile) [![codecov](https://codecov.io/gh/mccurdyc/splitfile/branch/master/graph/badge.svg)](https://codecov.io/gh/mccurdyc/splitfile) [![Maintainability](https://api.codeclimate.com/v1/badges/8b473a645aab19597124/maintainability)](https://codeclimate.com/github/mccurdyc/splitfile/maintainability)

splitfile attempts to make Go code more readable by identifying clear splits in
your package files.

## How it Works

<p align="center">
  <img width="500" height="500" src="https://github.com/mccurdyc/splitfile/blob/master/docs/imgs/splitfile-graph.png?raw=true">
</p>

splitfile builds a graph where nodes are `type`, `var`, `const,` or `func` declarations
and undirected edges define the relationships between declarations. splitifle identifies
"splits" or the set of dependent declarations. Where there is not a clear split,
declarations will be defined as "common". Method declarations are an exception and
will not be defined as common even if they are a common node in the graph. This is
because of the way methods are declared and tied to exactly one type.
