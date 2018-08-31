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
  -loglevel int
    	0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug (default 4)
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
So what does pulli really do?

**fetch** and **pull**

So basically it just walks through the filesystem searching for git repositories.

When a repo is found (and passed the filters) a **git pull**.

So the only change that is made to the working tree might be a **fast-forward merge** during pull.


### fun fact
to reveal the facepalming introduction:

**pulli is the german colloquial for pullover**

