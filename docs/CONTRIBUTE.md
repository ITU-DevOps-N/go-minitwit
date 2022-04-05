# Contribute

## Repository setup
  We use a [GitHub Organization](https://github.com/ITU-DevOps-N) to manage our repositories.
  For this repository we have set up branch protection rules where the `main` branch requires two approvals and `develop` branch requires one approval from other members to merge.
## Branching model
- `main`: Entirely stable code should reside here, possibly only code that has been or will be released. Code in this branch will go through the CI/CD pipeline.
- `develop`: A parallel branch that is worked from or used to test stability — it is not necessarily always stable, but whenever it gets to a stable state, it can be merged into main. Used to pull in topic branches, i.e. hotfixes, features etc. Tested on develop and merged into main.
- `feature/` branches: A short-lived branch that you create and use for a single particular feature or related work.

## Distributed development workflow
A contributor will create a `feature/` branch in its local environment where it is possible to work on changes. Once the feature is ready, the contributor will create a pull request from the feature branch to the `develop` branch. All the main functionalities of the system should work in the `feature/` branch before approving the pull request. At least one reviewer will be required to approve the pull request. Once the pull request has been merged, the `feature/` branch should be deleted locally and remotely. Once a stable release is made, the `develop` branch will be merged into the `main` branch. At least two reviewers will be required to approve the pull request. 

A contribution is represented by a pull request.