# Pretty-Parallel

Runs commands in parallel and display a pretty output

## Why?
It's easy enough to run commands in parallel, but the output will be messy and you will need to handle exit status properly to know what failed. Pretty-Parallel solve those problems prettily 

## How to install

1. Go the [Releases Page](https://github.com/tlopo-go/pretty-parallel/releases)
2. Download the desired version
3. Extract
4. Save the flat binary to a directory in your path

With one-liner: 
```
curl 'https://github.com/tlopo-go/pretty-parallel/releases/download/v0.0.3/pretty-parallel_0.0.3_darwin_arm64.tar.gz' -s -L | tar -C /usr/local/bin  -xzf -  pretty-parallel
```

One-liner on windows git-bash
```
curl 'https://github.com/tlopo-go/pretty-parallel/releases/download/v0.0.3/pretty-parallel_0.0.3_windows_amd64.tar.gz' -s -L | tar -C /usr/bin  -xzf -  pretty-parallel.exe
```

## USAGE:
```
Version: 0.0.3, Commit: 2372ca69122903aad35ea8630e00584a8f0d1a74, Date: 2025-12-03T18:43:56Z

USAGE:
pretty-parallel [OPTS] < input
  -c int
    	concurrency (default 10)

NOTES:
    Input must be either yaml or json following the schema below:
    [
        { name: string, cmd: [string|[]string] },
        ...
    ]
```
## Examples: 
Ping multiple hosts in parallel: 
```bash
for ip in 8.8.8.8 8.8.4.4; do echo - { name: $ip, cmd: "ping -c 2 $ip" }; done  | pretty-parallel 
```


