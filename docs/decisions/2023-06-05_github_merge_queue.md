# 2023-06-05 GitHub Merge Queue

* status: accepted
* date: 2023-06-05
* authors: Brian Samek

## Context and Problem Statement

* [Epic](https://go/github-merge-queue)
* [Scope](https://go/github-merge-queue-scope)
* [Design](https://go/github-merge-queue-design)

GitHub merge queue integrates natively with GitHub Actions. Integration with
other CI systems exists but is lightly documented. The
[docs](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-a-merge-queue#triggering-merge-group-checks-with-other-ci-providers)
state that a CI system should listen for branches matching a pattern and that
the PR will be merged once required checks pass.

The purpose of this ADR is to keep a record of the details of the integration
between the GitHub merge queue and Evergreen for later reference.

## Architectural Overview

Evergreen listens for new GitHub merge queue items, creates and runs versions
for these items, and posts the results of these versions to GitHub.

### Listen

Evergreen listens for new branches created by the merge queue by listening for
the merge_group event.  These are generated by a user clicking the add to queue
or (eventually) the merge when ready buttons in the GitHub UI.

### Create version

Evergreen creates a new version with a new requester type "github-merge-queue".
The variants and tasks will be those defined in the existing commit queue
settings on the project configuration page.

There are two alternative possibilities for the requester type: First, use the
ad hoc version type.  However, the ad hoc version type cannot independently be
given increased priority in the scheduler. Using ad hoc overloads the meaning of
this type.  Second, use the commit queue type. However, this type triggers a
substantial amount of Evergreen commit queue logic, so it's more straightforward
to make a new type. A disadvantage is that anyone downstream who has written
queries that assumes the number of requesters is fixed will need to take into
account the new requester.

### Post results

Post a check to the checks API. Note that GitHub "status checks" are of two
types, "statuses" and "checks". Branch protection rules are fulfilled by
statuses, whereas the GitHub merge queue listens for a check. The merge when
ready button waits for the branch protection rules to be fulfilled before adding
the PR to the merge queue.
