# Gittree

`git branch --list` but prettier

```
baguette@~$ gittree
main
  b1
    b0*
      b2
      b3
```

```
baguette@~$ gittree --help
List branches of a git repository in a tree structure

Usage:
  gittree [flags]
  gittree [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List branches of a git repository in a tree structure

Flags:
  -h, --help          help for gittree
  -p, --path string   Path to the git repository (default ".")

Use "gittree [command] --help" for more information about a command.
```
