/*
Package client provides initialization code for connecting to an incus server,
using the same configuration as used by the "incus" command.

It uses the environment variables
- INCUS_SOCKET
- INCUS_DIR
- INCUS_CONF

and the paths:
- /var/lib/incus/
- ~/.config/incus/
*/
package client
