# Project root

![example workflow](https://github.com/magdyamr542/project-root/actions/workflows/doBuild.yaml/badge.svg)

[![ezgif.com-gif-maker578618f2773b60af.gif](https://s10.gifyu.com/images/ezgif.com-gif-maker578618f2773b60af.gif)](https://gifyu.com/image/SSISi)

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)

## Introduction

#### Stop wasting your time with `cd ../`

- This enables you to return to the root directory of your project faster.
- It's meant for people who find it annoying to type `cd ../` a lot of times `:D`
- If you are in `./src/api/routes` you can simply use `pr` and you will automatically **cd** to the root of the project if its path was saved was `pr add` command.

```
.
├── dist
├── node_modules
├── package.json
├── src
│   ├── api
│   │   └── routes
│   │       └── getProduct.ts
│   └── utils
│       └── array.ts
└── tsconfig.json

6 directories, 4 files
```

## Installation

1. Run `mkdir ~/.proot`
1. Run `touch entryPoint.sh`
1. Put this in the file  `entryPoint.sh`
  ```bash
  function pr {
      output=$(./proot $@)
      retCode=$?
      if [[ ( $@ == "go" || $@ == "" ) && $retCode -eq 0 ]]; then
          # cd when go command. hide output
          cd $output
      else
          echo $output
      fi
      if [ $retCode -ne 0 ]; then
          return $retCode
      fi
  }
  ```
1. Install the executable in www.google.com and put it in the directory `~/.proot`
1. Add `source ~/.proot/entryPoint.sh` to your `.bashrc` or `.zshrc`. (This step is important)
1. Source your shell again or open a new terminal session.
1. Use `pr help`

## Usage

1. Use `pr add [path]` to save the current path as a project root.
1. Use `pr` or `pr go` from any child directory of your project and you will jump to the root of the project as long as you saved its path like in step 1.
1. Use `pr list` to list all saved project root paths.
1. Use `pr purge` to delete any paths that were saved before but don't exist in the file system any more.
1. Use `pr clear` to delete all saved paths. (use with **CAUTION**)
