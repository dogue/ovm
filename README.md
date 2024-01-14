# Odin Version Manager

A simple tool for installing and updating [Odin](https://github.com/odin-lang/Odin) and [OLS](https://github.com/DanielGavin/ols).

## What It Do

OVM can install Odin (and optionally OLS, the Odin language server) for you. It can also remove or update those installations.

## What It Don't

OVM is not *really* a "version mananger". Not yet. Odin is still pre-1.0 and the current recommended installation method is to clone the repo and build it yourself.

This is what OVM does for you.

As such, it has all the same system requirements that the manual process does.

>[!NOTE]
> While Odin supports Windows and Mac, OVM does not at this time. I don't own a Mac and I very rarely boot into Windows. If you have these systems and would like to work on support for them for OVM, pull requests are open.

## Features

You can install a specific commit for either Odin or OLS, though updating will always update to the latest changes.

It's assumed that if you're installing a specific commit, you have a good reason (so you will probably stay on that commit).

## Usage

OVM is incredibly simple. It has three commands:

* `install`
* `update`
* `remove`

### Install

Install has the most options:

* `-f, --force`: forces OVM to overwrite an existing installation
* `-l, --lsp`: if set, OLS will be installed too
* `-c, --odin-verson string`: `string` is a commit hash for a specific version of Odin 
* `-s, --ols-version string`: same as above, but for OLS
* `-p, --path`: sets a base path - Odin (and OLS) will be installed *inside* of this directory (defaults to `$XDG_DATA_HOME`)

OVM stores a small text file in `$XDG_CONFIG_HOME/ovm` containing the base path used for the installation. This is used for the update/remove command.

### Update

Update has no options. It will pull the latest changes for Odin (and OLS if present) and rebuild.

### Remove

Remove also has no options. It will remove the Odin and OLS directories, as well as a file that OVM creates containing the base path that was given during installation.
