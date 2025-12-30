# gittree

`git branch -l` in tree view. 

You can see which branch was created from which branch in a tree format, making the whole merging process a lot more straightforward.

Note that not all git branch relations will result in a tree due to cycles. I didn't handle this edge case just yet.

## Installation
``` bash
brew tap mucansever/gittree
brew install mucansever/gittree/gittree
```

## Usage

Run `gittree list` inside your repo.
```bash
.
└── main
    ├── fix/important-bug
    └── feat/feature-1
        └── chore/document-change*
```

If there is two of the same branch, all the children are duplicated. In below example, `feat/no-commit-branch` is a new branch from `main` without any commits.
```bash
.
├── main
│   ├── fix/important-bug
│   └── feat/feature-1
│       └── chore/document-change
└── feat/no-commit-branch*
    ├── fix/important-bug
    └── feat/feature-1
        └── chore/document-change
```

## Improvements

Please create an issue for any improvement that you might think of.
