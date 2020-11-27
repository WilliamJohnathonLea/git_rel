# GitHub Release Creator
Create a release of your code with this handy CLI tool.

This tool uses [Device flow](https://docs.github.com/en/free-pro-team@latest/developers/apps/authorizing-oauth-apps#device-flow) to authenticate itself.

## Installation
- Download the correct binary for your OS. Look [here](https://github.com/WilliamJohnathonLea/git_rel/releases) to find the latest release.
- Copy the binary to your PATH as `gitrel`
- Check it's installed by running `gitrel --help`

## Usage
- Create a new patch version
```
gitrel --version=patch owner/repo_name
```
- Create a new minor version
```
gitrel --version=minor owner/repo_name
```
- Create a new major version
```
gitrel --version=major owner/repo_name
```
