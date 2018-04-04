## CLI Guide

### katib-cli

#### Usage

```
./katib-cli [options] [arguments]
```

#### Arguments

- `Getstudies` Get list of studys and their status.
- `Createstudy` Send create new study request to katib api server.
- `Stopstudy` Delete specified study from API server.
But the results of trials in modelDB won't be deleted.

#### Options

- `-s string`
Set address of vizier-core. {IP Addr}:{Port}. default localhost:6789
Katib API is grpc.
Unfortunately, nginx ingress controller does not support grpc now ( next version it will support! )
So vizier-core expose port as NodePort(30678}.
- `-f` Specify the config file of your study, which is used in `Createstudy`
