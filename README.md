# tapioca
a tool to summarize commits of today  

## Instllation
```
go install github.com/tomocy/tapioca/cmd/tapioca
```

## Example
- Summarize commits to tomocy/tapioca in 2019/08/12
```
$ date
Mon Aug 12 03:50:56 JST 2019

$ tapioca -r tomocy/tapioca -a tomocy
summary of commits to tomocy/tapioca in 2019/08/12
305 changes: 173 adds, 132 dels
```

## Usage
```
Usage of tapioca:
  -a string
        name of author
  -f string
        name of format (default "text") (available "text", "color")
  -m string
        name of mode (default "cli") (available "cli", "twitter")
  -r string
        name of owner/repo
```
