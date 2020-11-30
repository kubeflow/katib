## Katib Operators

### Overview
This bundle encompasses the Kubernetes python operators (a.k.a. charms) for Katib
(see [CharmHub](https://charmhub.io/?q=katib)). 

The Katib operators are python scripts that wrap the latest [Katib manifests](manifests),
providing lifecycle management for each application, handling events (install, upgrade,
integrate, remove).

[manifests]: https://github.com/kubeflow/manifests/tree/master/katib

## Install

### Install applications

To install Katib, run

    juju deploy katib

You can also install each application individually, like this:

    juju deploy <application>

where `<application>` is one of `katib-controller`, `katib-ui`, or `katib-manager`.
