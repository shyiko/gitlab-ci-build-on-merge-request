# gitlab-ci-build-on-merge-request ![Latest Version](https://img.shields.io/badge/latest-0.3.1-blue.svg)

Missing "build on Merge Request" feature of Gitlab CI.

Build is triggered when merge request is first opened and when commits are added to the source branch later on
(up until the merge/close event).

> Tested on Gitlab CI 8.7.2.

## Setup

* Start `./gitlab-ci-build-on-merge-request`. If you are using docker:

```sh
docker run -it --rm --name gitlab-ci-build-on-merge-request \
  shyiko/gitlab-ci-build-on-merge-request:0.3.1 \
  gitlab-ci-build-on-merge-request --url=http://gitlab.example.com --private_token=<authentication token>
```

> Authentication token can be found on `Profile Settings` -> `Account` -> `Private Token`.

* Now that you have service running:

  - Make sure job definitions in `.gitlab-ci.yml` file have "triggers" policy enabled (see below).

    ```yaml
    build:
      script:
        - ...
      only:
        - master
        - develop
        - /^release-.*$/
        # the line below is required by shyiko/gitlab-ci-build-on-merge-request
        - triggers
    ```

  - Create a webhook (go to `Project Settings` -> `Webhooks`) that points to `gitlab-ci-build-on-merge-request`
  server (e.g. "http://gitlab-ci-build-on-merge-request.example.com/hook") and has "Merge Request events" as a trigger.

  > Starting from 0.2.0 private token can also be specified using hook's query parameter "private_token" (see https://github.com/shyiko/gitlab-ci-build-on-merge-request/pull/2).

That's it.

## Building your own Docker image

```sh
docker build -t custom-image-name .
```

## Development

> PREREQUISITE: [go1.6](https://github.com/moovweb/gvm)

```sh
git clone https://github.com/shyiko/gitlab-ci-build-on-merge-request \
  $GOPATH/src/github.com/shyiko/gitlab-ci-build-on-merge-request
cd $GOPATH/src/github.com/shyiko/gitlab-ci-build-on-merge-request

go build
```

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
