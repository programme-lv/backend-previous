[programme.lv]: https://programme.lv
[database]: https://github.com/programme-lv/database 

[Deploy Action Status Badge]: https://github.com/programme-lv/backend/actions/workflows/deploy.yml/badge.svg
[Test Action Status Badge]: https://github.com/programme-lv/backend/actions/workflows/test.yml/badge.svg
[Go Report Card]: https://goreportcard.com/badge/gojp/goreportcard
[License]: https://img.shields.io/badge/license-GPLv3-blue

[programme.lv] is a programming learning platform for scholars, students and others.

[![Test Action Status Badge]](https://github.com/programme-lv/backend/actions/workflows/test.yml)
[![Deploy Action Status Badge]](https://github.com/programme-lv/backend/actions/workflows/deploy.yml)
[![Go Report Card]](https://goreportcard.com/report/github.com/programme-lv/backend)
[![License]](https://github.com/programme-lv/backend/blob/main/LICENSE)

## Overview

This repository contains source code for the **backend** of the [programme.lv] system.
The **backend** is a GraphQL server that interacts with the PostgreSQL [database] and the RabbitMQ submission queue.

## Git workflow & CI/CD

The two core branches of this repo are: `dev` & `main`:

- `dev` is the branch that contains the latest development version of the code.
Unit & integration tests are not automatically run against this branch.
You can push to this branch directly, but it is recommended to
create a feature branch and then create a pull request to `dev`.

- `main` is a staging branch, where the code is tested before being deployed to production.
After each push on `main`, unit and integration tests are run.
Integration tests are run against the newest version of the [database].
If a commit is tagged it is then deployed to production.

For more information look at actions defined in `./.github/workflows/`.

TODO: trigger the testing workflow also on database changes.

## Contributing

When contributing to this repository, please first discuss the change you wish
to make via issue, email, or any other method with the owners of this repository
before making a change.

Pull requests are the best way to propose changes to the codebase. We actively
welcome your pull requests:

1. Fork the repo and create your branch from `master`.
2. If you've added code that should be tested, add tests.
3. If you've added code that need documentation, update the documentation.
4. Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
5. Issue that pull request!

Join our community on [Discord](https://discord.gg/7c8GwpGt)!

