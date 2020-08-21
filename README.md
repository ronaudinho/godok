[![Imgur](https://i.imgur.com/z7wBEwj.png)](https://i.imgur.com/z7wBEwj.png)
:trollface:

## what in the world is this
`godok` generates a boilerplate godoc comments

## how to use
`examples/main.go` uses `godok` go generate boilerplate godoc comments in a given dir

## tests dir
- `.go` for actual go code
- `.godok` for current expected output
- `.gowant` for roadmap
### why?
- creating different directories to match file make it more difficult to see (for me at least)
- working with `.go` file requires convoluted workaround to allow redeclaration of var/func (i stupid don't judge)

## milestones
- idiomatic Go comments
- allow flags usage
- allow template declaration
