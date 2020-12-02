## Katib Operators

### Overview
This bundle encompasses the Kubernetes python operators (a.k.a. charms) for Katib
(see [CharmHub](https://charmhub.io/?q=katib)). 

The Katib operators are python scripts that wrap the latest released [Katib manifests][manifests],
providing lifecycle management for each application, handling events (install, upgrade,
integrate, remove).

[manifests]: https://github.com/kubeflow/katib/tree/master/manifests

## Install

### Install applications

To install Katib, run:

    juju deploy katib

You can also install each application individually, like this:

    juju deploy <application>

where `<application>` is one of `katib-controller`, `katib-ui`, or `katib-db-manager`.

> Note: As a default, when you `juju deploy` an application or the full Katib bundle, you will deploy the latest released version of Katib, even if unreleased updates are already available. If you would like to try the latest available charm run `juju deploy foo --channel=stable/edge`.
