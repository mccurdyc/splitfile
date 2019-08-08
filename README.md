# splitfile

splitfile identifies clear splits in your files to make your Go package more readable.

In Go, packages can be comprised of multiple files.

## How it Works

<p align="center">
  <img width="460" height="300" src="https://github.com/mccurdyc/splitfile/blob/master/docs/imgs/splitfile-graph.png?raw=true">
</p>

splitfile builds a graph where nodes are `type`, `var`, `const,` or `func` declarations
and undirected edges define the relationships between declarations. splitifle identifies
clear "splits" or where a set of declarations are only used by a small set of other
declarations.
