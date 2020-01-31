[![Go Report Card](https://goreportcard.com/badge/github.com/Oppodelldog/pulli)](https://goreportcard.com/report/github.com/Oppodelldog/pulli)
[![Linux build](http://nulldog.de:12080/api/badges/Oppodelldog/pulli/status.svg)](http://nulldog.de:12080/Oppodelldog/pulli)
# pulli
> pulli pulls git repos recursively

**Put your pulli on, it's getting cool.**

It's always nice to have a pulli.
Especially when
- you are a cool developer
- your are working in a team of even cooler developers
- all of you produce a lot of dependent git repositories
- it is annoying to keep all repos up to date
- you want to be the coolest guy

So put **pulli** on and pull your git repositories recursively.

### Usage
```bash
Usage of pulli:
  -dir string
    	defines the folder where to find git repos (default ".")
  -filter value
    	filters the given folder. (can be absolute path or regex)
  -filtermode string
    	whitelist or blacklist
```

### Examples
```bash
# execute in your projects folder
pulli

# exclude a folder by filtering a blacklist
pulli -filtermode blacklist -filter /tmp/some-directory/big-bad-repo

# filters might be regex (golang regex)
pulli -filtermode blacklist -filter /tmp/some-directory/test.*

# include a folder by using filtermode whitelist
pulli -filtermode whitelist -filter /tmp/some-directory/big-bad-repo

# define multiple filters
pulli -filtermode whitelist -filter /tmp/a -filter /tmp/b

# define directory for discovery
pulli -dir /tmp/some-directory

```

### Behavior
**Rambo-like explanation:**

Q: "And this, what is this?"

A: "It's pulli"

Q: "What does it do?"

A: "It pulls"


**Seriously:**

Pulli walks through the filesystem searching for git repositories.

When a repo is found and it passed the filters a **git pull** is executed.

So the only change that might be made to the working tree could be a **fast-forward merge** from ```git pull```.


### fun fact
to reveal the facepalming introduction:

**pulli is the german colloquial for pullover**

