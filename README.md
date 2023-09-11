[programme.lv]: https://programme.lv
[database]: https://github.com/programme-lv/database 
[deploy action status badge]: https://github.com/programme-lv/backend/actions/workflows/deploy.yml/badge.svg
[test action status badge]: https://github.com/programme-lv/backend/actions/workflows/test.yml/badge.svg
[RabbitMQ]: https://www.rabbitmq.com/

[programme.lv] is a programming learning platform for scholars, students and others.

![deploy action status badge]
![test action status badge]

## Overview

This repository contains source code for the **backend** of the [programme.lv] system.
The **backend** is a `graphql` server that interacts with the `PostgreSQL` [database] and the [RabbitMQ] submission queue.


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

TODO: write how to contribute
