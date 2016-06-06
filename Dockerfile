FROM golang:1.6.2-onbuild
MAINTAINER Stanley Shyiko <stanley.shyiko@gmail.com>

RUN ln -s /go/bin/app /go/bin/gitlab-ci-build-on-merge-request
