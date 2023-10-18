cloudinitincus configures Incus instances (containers, VMs), by
applying cloud-init config files, using the Incus API.

It does not require that the guest OS has cloud-init installed.

It can be used as a standalone executable, or as a Go library.

# CLI Usage
```
cloudconfig-incus apply [-ostype <ostype>] -i <instance> <cloud-config-file>...
```

You can find example cloud-config files here:
https://github.com/melato/cloudconfig/tree/main/examples

for example:

```
cd cloudconfig/examples
cloudconfig-incus apply -i <instance> files.yaml
```

ostype is needed for packages and users, but not for write_files or runcmd.

Supported ostypes are: alpine, debian.
  

# Supported cloud-init features
The following cloud-init modules (sections) are supported and applied in this order:

- packages
- write_files (only plain text, without any encoding)
- users
- runcmd

# Standalone executable
cloudconfig-incus will connect to the Incus server
using either the UNIX socket or HTTPS.
For HTTPS, it uses configuration information found in the incus user configuration files,
in ~/.config/incus

# compile

```
cd main
date > version
go install cloudconfig-incus.go
```
