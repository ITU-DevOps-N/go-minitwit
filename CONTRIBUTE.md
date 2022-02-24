# Contribute

## Repository setup
  We use a [GitHub Organization](https://github.com/ITU-DevOps-N) to manage our repositories.
  For this repository we have set up branch protection rules where `main` and `develop` branch requires approval from 2-3 members to merge.
## Branching model
- `main`: Entirely stable code should reside here, possibly only code that has been or will be released.
- `develop`: A parallel branch that is worked from or used to test stability — it isn’t necessarily always stable, but whenever it gets to a stable state, it can be merged into Master. Used to pull in topic branches, i.e. hotfixes, features etc. Tested on develop and merged into Main.
- `feature/` branches: A short-lived branch that you create and use for a single particular feature or related work

## Distributed development workflow
A contributor will create a `feature/` branch in its local environment where it's possible to work on some changes. Once the feature is ready, they will create a pull request from the feature branch to the `develop` branch. All the main functionalities of the system should work in the `feature/` branch before approving the pull request. At least two reviewers will be required to approve the pull request. The `feature/` branch will be deleted. Once a stable release is made, the `develop` branch will be merged into the `main` branch. At least three reviewers will be required to approve the pull request. 

A contribution is rapresented by a pull request.
