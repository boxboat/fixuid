# fixuid

![Build Status](https://github.com/boxboat/fixuid/workflows/Main/badge.svg)

`fixuid` is a Go binary that changes a Docker container's user/group and file permissions that were set at build time to the UID/GID that the container was started with at runtime.  Primary use case is in development Docker containers when working with host mounted volumes.

`fixuid` was born because there is currently no way to remap host volume UIDs/GIDs from the Docker Engine, [see moby issue 7198](https://github.com/moby/moby/issues/7198) for more details.

Check out [BoxBoat's blog post](https://boxboat.com/2017/07/25/fixuid-change-docker-container-uid-gid/) for a practical explanation of how `fixuid` benefits development teams consisting of multiple developers.

**fixuid should only be used in development Docker containers.  DO NOT INCLUDE in a production container image**

# Overview 

- build a Dockerfile with user/group `dockeruser:dockergroup` that has UID/GID `1000:1000`
- host is running as UID/GID 1001:1002, host mounted volume has permissions 1001:1002
- run the docker container with argument `-u 1001:1002` so that container is now running with same UID/GID as host
- `fixuid` can run as an entrypoint or in a startup script and performs the following:
  - changes `dockeruser` UID to 1001
  - changes `dockergroup` GID to 1002
  - changes all file permissions for old `dockeruser:dockergroup` to 1001:1002
  - updates $HOME inside container to `dockeruser` $HOME
- now container UID/GID matches host UID/GID and files created in the container on the host mount will be correct

## Motivation

Common Docker development workflows involve mounting source code into a container via a host volume.  Build tools such as `gradle`, `yarn`, `webpack`, etc. download dependencies and create files in the host mount.

Many times the UID/GID of the build tools running in the Docker container do not match the UID/GID of the mounted host volume, and files generated in the container do not match files in the host volume.  This can lead to problems, such as an IDE running on the host not able to modify a file that was created by the container due to a file ownership mismatch.

In large development teams, it is possible to have many developers running as different UIDs/GIDs on their host systems.  With `fixuid`, individual developers can run the same container using the appropriate UID/GID for their host environment.

## Install fixuid in Dockerfile

1. Create a non-root user and group inside of your docker container.  Use any UID/GID, 1000:1000 is a good choice.

    Note: some images already create UID/GID 1000:1000 for you, e.g. `nodejs` creates user/group `node:node` as UID/GID 1000:1000.  In this case you can skip this step and use the `node:node` user/group.

```
# sample command to create user/group on different base images
# creates user "docker" with UID 1000, home directory /home/docker, and shell /bin/sh
# creates group "docker" with GID 1000

# alpine
RUN addgroup -g 1000 docker && \
    adduser -u 1000 -G docker -h /home/docker -s /bin/sh -D docker
    
# debian / ubuntu
RUN addgroup --gid 1000 docker && \
    adduser --uid 1000 --ingroup docker --home /home/docker --shell /bin/sh --disabled-password --gecos "" docker

# centos / fedora
RUN groupadd -g 1000 docker && \
    useradd -u 1000 -g docker -d /home/docker -s /bin/sh docker
```

2. Install `fixuid` in the container, ensure that root owns the file, make it execuatble, and enable the [setuid bit](https://en.wikipedia.org/wiki/Setuid).  Create the file `/etc/fixuid/config.yml` with two lines, `user: <user>` and `group: <group>` using the user and group from step 1.

    Note: this command must be run as root and requires that `curl` is installed in the container

```
RUN USER=docker && \
    GROUP=docker && \
    curl -SsL https://github.com/boxboat/fixuid/releases/download/v0.6.0/fixuid-0.6.0-linux-amd64.tar.gz | tar -C /usr/local/bin -xzf - && \
    chown root:root /usr/local/bin/fixuid && \
    chmod 4755 /usr/local/bin/fixuid && \
    mkdir -p /etc/fixuid && \
    printf "user: $USER\ngroup: $GROUP\n" > /etc/fixuid/config.yml
```

3. Set the default user/group to `user:group` and set the entrypoint to `fixuid`.

```
USER docker:docker
ENTRYPOINT ["fixuid"]
```

4. Run the container using UID/GID of your host.  Replace `1000:1000` with your host's `UID/GID`:

```
docker run --rm -it -u 1000:1000 <image name> sh
```

## Set Default Values inside of Docker Compose

Set a default UID and GID for the container to run as inside of the `docker-compose.yml` file.  Developers who are running as a different UID/GID on their host can override the defaults using environment variables or a [.env file](https://docs.docker.com/compose/env-file/)

```
nginx:
  image: my-nginx
  user: ${FIXUID:-1000}:${FIXGID:-1000}
  volumes:
    - ./nginx:/etc/nginx
    - ./www:/var/www
```

## Specify Paths and Behavior across Devices

The default behavior of `fixuid` is to start at the root path `/` and recursively scan each file and directory on the same devices as `/`.  In the configuration file `/etc/fixuid/config.yml`, you can specify specify the directories that should be recursively scanned:

```yaml
user: docker
group: docker
paths:
  - /home/docker
  - /tmp
```

`fixuid` will only recurse into a directory as long as it is on the same initial device specified in `paths` and will not recurse into directories mounted on other devices.  This includes Docker volumes.  If you want `fixuid` to run on the root Docker filesystem and a Docker volume at `/home/docker/.cache`, your configuration should include:

```yaml
user: docker
group: docker
paths:
  - /
  - /home/docker/.cache
```

## Run in Startup Script instead of Entrypoint

You can run `fixuid` as part of your container's startup script.  `fixuid` will `export HOME=/path/to/home` if $HOME is the default value of `/`, so be sure to evaluate the output of `fixuid` when running as a script.  Supplementary groups will not be set in this mode.

```
#!/bin/sh

# UID/GID map to unknown user/group, $HOME=/ (the default when no home directory is defined)

eval $( fixuid )

# UID/GID now match user/group, $HOME has been set to user's home directory
```

## Command-Line Flags

`fixuid` has the following command-line flags:

```
Usage of ./fixuid:
  -q	quiet mode
```
