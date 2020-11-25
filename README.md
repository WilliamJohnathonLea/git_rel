# GitHub Release Creator
Create a release of your code with this handy CLI tool.

## Usage
- Set GITHUB_TOKEN environment variable with a GitHub personal access token
- Create a new patch version
```
git_rel -version=patch owner/repo_name
```
- Create a new minor version
```
git_rel -version=minor owner/repo_name
```
- Create a new major version
```
git_rel -version=major owner/repo_name
```