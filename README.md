# Regit
## Description
Regit is for git commands that support regular expressions.

This allows you to work with specific files or branches in bulk.

## Setup
```
$go install github.com/YutaSeki36/regit
```
## Usage

### checkout
Checkout is used for a file restore

```
Usage:
  regit checkout [flags]

Flags:
  -h, --help            help for checkout
      --ours
  -t, --target string   Set the target file name to check out with a regular expression.
      --theirs

Global Flags:
  -d, --dryRun   dryRun enable flag
```

### del_branch
Del_branch is used to delete a branch.

```
Usage:
  regit del_branch [flags]

Flags:
  -h, --help            help for del_branch
  -t, --target string   Set the branch name to be deleted with a regular expression.

Global Flags:
  -d, --dryRun   dryRun enable flag
```

example.
```
$regit del_branch -t feature.*
branch delete targets: feature/hoge1 feature/hoge10 feature/hoge2 feature/hoge3 feature/hoge4 feature/hoge5 feature/hoge6 feature/hoge7 feature/hoge8 feature/hoge9
Finished
```

## Author
Yuta Seki