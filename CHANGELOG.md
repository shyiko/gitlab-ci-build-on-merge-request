# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [0.4.0] - 2018-12-17

### Added 
- `--insecure` flag ([#7](https://github.com/shyiko/gitlab-ci-build-on-merge-request/pull/7)).

## [0.3.3] - 2018-03-26

### Fixed
- GitLab 9+ compatibility.

## [0.3.2] - 2018-02-27

### Fixed
- Error reported in case of invalid `private_token`.

## [0.3.1] - 2018-01-16

### Fixed
- Trigger ownership.  

## [0.3.0] - 2017-08-02

### Changed 
- **BREAKING**: GitLab API target to v4  
(v3 was deprecated in GitLab 9 and will be removed in GitLab 9.5 on August 22, 2017).

## [0.2.0] - 2017-06-16

### Added
- `private_token` hook query parameter support ([#2](https://github.com/shyiko/gitlab-ci-build-on-merge-request/pull/2)).

## [0.1.1] - 2017-04-26

### Fixed
- Triggering `when: manual`.

## 0.1.0 - 2016-06-05

[0.4.0]: https://github.com/shyiko/gitlab-ci-build-on-merge-request/compare/0.3.3...0.4.0
[0.3.3]: https://github.com/shyiko/gitlab-ci-build-on-merge-request/compare/0.3.2...0.3.3
[0.3.2]: https://github.com/shyiko/gitlab-ci-build-on-merge-request/compare/0.3.1...0.3.2
[0.3.1]: https://github.com/shyiko/gitlab-ci-build-on-merge-request/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/shyiko/gitlab-ci-build-on-merge-request/compare/0.2.0...0.3.0
[0.2.0]: https://github.com/shyiko/gitlab-ci-build-on-merge-request/compare/0.1.1...0.2.0
[0.1.1]: https://github.com/shyiko/gitlab-ci-build-on-merge-request/compare/0.1.0...0.1.1
