# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [1.0.1]

### Changed

- Add option for setting `http2-max-field-size` via values.yaml.
- Add optional http snippet parameter to extend nginx http configuration.

## [1.0.0]

### Changed

- Upgrade nginx controller to use the latest version [0.26.1](https://github.com/kubernetes/ingress-nginx/releases/tag/nginx-0.26.1).  
**This version includes the following breaking changes**:
   - The variable $the_real_ip variable was removed from template and default log_format.
   - The default value of configmap setting proxy-add-original-uri-header is now "false".
   - When the setting proxy-add-original-uri-header is "true", the ingress controller adds a new header X-Original-Uri with the value of NGINX variable $request_uri.
   - In most of the cases this is not an issue but with request with long URLs it could lead to unexpected errors in the application defined in the Ingress serviceName, like issue 4593 - 431 Request Header Fields Too Large

## [0.10.2]

### Removed

- Remove migration subcharts that are no longer required.

## [0.10.1]

### Changed

- Add prometheus annotations to expose nginx controller metrics.

## [0.10.0]

### Added

- Add `NET_BIND_SERVICE` capability for ingress-controller psp.

### Changed

- Reduce max replicas for HPA from 15 to 11 based on a recommendation from upstream.

## [0.9.1]

- Fix a security issue on nginx-ingress by using the latest version [0.25.1](https://github.com/kubernetes/ingress-nginx/releases/tag/nginx-0.25.1). 

    - Patch release to fix several security issues were identified in nginx HTTP/2 implementation, which might cause excessive memory consumption and CPU usage (CVE-2019-9511, CVE-2019-9513, CVE-2019-9516).


## [0.9.0]

### Added

- Upgrade nginx controller to use the latest version [0.25.0](https://github.com/kubernetes/ingress-nginx/releases/tag/nginx-0.25.0). Notable changes there:

    - Support new [networking.k8s.io/v1beta1 package](https://github.com/kubernetes/ingress-nginx/pull/4127).

    - Migration from Nginx to OpenResty 1.15.8.

    - This version adds a validation webhook controller for ingress resource validation but it has not been enabled yet. As we found some issues during the testing process.

### Changed

- Switch anti-affinity from `requiredDuringSchedulingIgnoredDuringExecution` to `preferredDuringSchedulingIgnoredDuringExecution`. With this we can run multiple nginx-ingress-controller pods on the same node when there are taints applied to other nodes.

## [0.8.1]

### Changed

- Allow all egress traffic via network policy.

## [0.8.0]

### Added

- Separate network policies for `nginx-ingress-controller` and `default-http-backend`.

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
