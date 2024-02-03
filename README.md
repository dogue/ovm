# Odin Version Manager

---

Odin Version Manager (ovm) is a tool for managing your Odin installs. With Odin
being pre-1.0 the closest thing to a "regular release" are the monthly releases,
and those are currently bugged on Linux. While OVM does require that you have
the dependencies required to build Odin from source, it handles the downloads,
extraction, and actual build for you.

At present, OVM is barely more than an automation for downloading and building
the Odin compiler (and Odin Language Server), if I'm being honest. However, the
plan is to adapt this tool as the language grows, eventually being a fully fledged
version manager.

This tool is \*ahem\* *heavily inspired* by [ZVM](https://github.com/tristanisham/zvm).
In fact, Tristan gave me his blessing to rip off ZVM so that's what I did. If you
find OVM useful, please do go check out ZVM and give Tristan some love. He's done
some great work over there (and over here, technically).

# Installing OVM

OVM lives entirely in `$HOME/.ovm` on all platforms it supports. Inside of the
directory, OVM will download new Odin versions and symlink whichever version you
specify with `ovm use` to `$HOME/.ovm/bin`. You should add this folder to your
path. OVM's installer will add OVM to `$HOME/.ovm/self`. You should also add this
directory as the environment variable `OVM_INSTALL`. The installer should handle
this for you automatically if you're on *nix systems, but you'll have to manually
do this on Windows. You can then add `OVM_INSTALL to your path.`

If you don't want to use OVM_INSTALL (like you already have OVM in a place you
like), then OVM will update the exact executable you've called `upgrade` from.

# Linux, BSD, MacOS, *nix

```sh
curl https://raw.githubusercontent.com/dogue/ovm/master/install.sh | bash
```

Then add OVM's directories to your `$PATH`

```sh
echo "# OVM" >> $HOME/.profile
echo export OVM_INSTALL="$HOME/.ovm/self" >> $HOME/.profile
echo export PATH="$PATH:$HOME/.ovm/bin" >> $HOME/.profile
echo export PATH="$PATH:$OVM_INSTALL/" >> $HOME/.profile
```

# Windows

If you're on Windows, please grab the
[latest release](https://github.com/dogue/ovm/releases/latest).

## Putting OVM on your Path

OVM requires a few directories to be on your `$PATH`. If you don't know how to
update your environment variables permanently on Windows, you can follow
[this guide](https://www.computerhope.com/issues/ch000549.htm). Once you're in
the appropriate menu, add or append to the following environment variables:

Add

- OVM_INSTALL: `%USERPROFILE%\.ovm\self`

Append

- PATH: `%USERPROFILE%\.ovm\bin`
- PATH: `%OVM_INSTALL%`

## Community Package

### AUR

TODO!

~~`ovm` on the [Arch AUR](https://aur.archlinux.org/packages/ovm) is a community
maintained package, and may be out of date.~~

# Why should I use OVM?

While Odin is still pre-1.0 if you're going to stay up-to-date with the master
branch, you're going to be downloading Odin quite often. You could do it
manually, having to scoll around to find your appropriate version, decompress
it, and install it on your `$PATH`. Or, you could install OVM and run
`ovm i master` every time you want to update. `ovm` is a static binary under a
permissive license. Whether you're on Windows, MacOS, Linux, a flavor of BSD,
or Plan 9 `zvm` will let you install, switch between, and run multiple versions of Zig.

# Contributing and Notice

`ovm` is alpha software. Pre-v1.0.0 any breaking changes will be clearly
labeled, and any commands potentially on the chopping block will print notice.
The program is under constant development, and the author is very willing to
work with contributors. **If you have any issues, ideas, or contributions you'd
like to suggest
[create a GitHub issue](https://github.com/dogue/ovm/issues/new/choose)**.

# How to use OVM

## Install

```sh
ovm install <version> 
# Or
ovm i <version>
```

Use `install` or `i` to download a specific version of Odin. To install the
latest monthly release, use "latest". To install from the master branch, use
"master". 

```sh
# Example
ovm i master
```

### Install OLS with OVM
 You can install OLS with your Odin download! To install OLS with OVM, simply pass the `-l/--lsp` flag with `ovm i`. For example:
```sh
ovm i master -l
```

## Switch between installed Odin versions

```sh
ovm use <version>
```

Use `use` to switch between versions of Odin.
Also available as `switch`.

```sh
# Example
ovm use master
```

## List installed Odin versions

```sh
# Example
ovm ls
```

Use `ls` to list all installed version of Odin.
Also available as `list`.

### List all versions of Odin available
```sh
ovm ls --remote
```
The `-r/--remote` flag will list the versions of Odin available for download rather than those locally installed.

## Uninstall a Odin version

```sh
# Example
ovm rm dev-2023-12
```

Use `remove` or `rm` to remove a locally installed version from your system.

## Upgrade your OVM installation

You can upgrade your OVM installation from ovm.
Just run:

```sh
ovm upgrade
```

The latest version of OVM should install on your machine, regardless of where
your binary lives (though if you have your binary in a privileged folder, you
may have to run this command with `sudo`).

## Print program version

```sh
ovm version
```

Prints the version of OVM you have installed.

## Toggle output colors

```sh
ovm colors
```

Shows whether colors are currently enabled and asks if you'd like to toggle them.

## Print program help

```sh
ovm help

```

<hr>

## Option flags

```sh
-v / --verbose | Enable more informational output from OVM
```
