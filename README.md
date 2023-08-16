# programme.lv backend

Programme.lv is a programming learning platform for scholars, students and others. This repository contains source code for the backend of the programme.lv system. CI / CD practices are implemented.

## Git workflow & CI/CD

The three important branches of this repo are: `main`, `dev`, `prod`.

- `dev` is the branch that contains the latest development version of the code. Unit & integration tests are not automatically run against this branch. You can push to this branch directly, but it is recommended to create a feature branch and then create a pull request to `dev`.

- `main` is a staging branch, where the code is tested before being deployed to production. After each push on `main`, unit and integration tests are run. Note that integration tests are run against the newest version of the database (main branch).

- `prod` is the branch that contains the code that is deployed to production. It is a protected branch, so one cannot push directly to it. To make changes, one can first create a pull request from `main` to `prod`, and then merge the pull request. The changes will be automatically built and deployed via GitHub Actions.

The backend should be the only part of the programme.lv system that initiates writes to the Postgres database.

During deployment:
- docker image is built;
- backend pods are first scaled to zero;
- database is migrated to the newest version;
- backend deployment manifest is updated;

TODO: trigger the testing workflow also on database changes.
