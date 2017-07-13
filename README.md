# fixuid

A Go binary to detect current running UID/GID and adapt

Primary use is in Docker so that container UID/GID matches host UID/GID when working with volumes mounted from host

Should only be used in development Docker containers.  DO NOT INCLUDE in a production container image.

More info coming soon
