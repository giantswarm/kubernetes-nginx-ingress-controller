# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [0.7.0]

### Added

- Custom serviceaccount for `default-http-backend`.
- Separate podsecuritypolicies for `nginx-ingress-controller` and `default-http-backend`.
- Run nginx-ingress-controller as `www-data` user instead of `root`.

## [0.6.0]

### Changed

- Upgrade nginx controller to use the latest version [0.24.1](https://github.com/kubernetes/ingress-nginx/releases/tag/nginx-0.24.1). Notable changes there:

    - Nginx version has changed to `1.15.10`, read more in the [project changelog](https://nginx.org/en/CHANGES).

    - Fix regression with dynamic SSL certificates when a custom default certificate is defined.

    - Nginx kubectl plugin has a new command to `lint` your ingress resources. That way you can observe deprecated annotations or properties. To install the plugin go [here](https://github.com/kubernetes/ingress-nginx/blob/29f7d9a77ade24a7366ef4a6f258b8aeef50678c/docs/kubectl-plugin.md).

[0.6.0]: https://github.com/giantswarm/kubernetes-nginx-ingress-controller/pull/90
