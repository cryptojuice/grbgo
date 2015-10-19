### grb - Simple Git Repository Helper

`grb` is a simple command line utility for displaying and deleting multiple local/remote branches. 

#### Usage Examples

List all branches:
$ grb

Return branches with contents matching search term:
$ grb "my-branch"

Delete branches with contents matching search term:
$ grb -d "my-branch"

Delete delete both local and remote branches with contents matching search term:
$ grb -d -l "my-branch"
