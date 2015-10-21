### grb - Simple Git Repository Helper

`grb` is a simple command line utility for displaying and deleting multiple local/remote branches. 

#### Installation
If you already have go available
```shell
$ go get github.com/cryptojuice/grbgo
```

#### Usage Examples

List all branches:
```shell
$ grb
```
List branches with contents matching search term:
```shell
$ grb "my-branch"
```

Delete branches with contents matching search term:
```shell
$ grb -d "my-branch"
```

Delete delete both local and remote branches with contents matching search term:
```shell
$ grb -d -l "my-branch"
```

Usage/Help
```shell
$ grb help
```
