# Changelog

## [0.6.0](https://github.com/boxboat/fixuid/releases/tag/v0.6.0) - 2023-08-17

### Features

- Call `syscall.Setgroups` on groups from `/etc/passwd` and `/etc/group`: [#37](https://github.com/boxboat/fixuid/pull/37)

## [0.5.1](https://github.com/boxboat/fixuid/releases/tag/v0.5.1) - 2021-07-19

### Features

- Add linux architectures `mips64` `mips64le` `ppc64` `ppc64le` and `riscv64`: [#33](https://github.com/boxboat/fixuid/pull/33), [#34](https://github.com/boxboat/fixuid/pull/34)

## [0.5](https://github.com/boxboat/fixuid/releases/tag/v0.5) - 2020-06-12

### Fixes

- Use Lchown so that symbolic links are not followed: [#27](https://github.com/boxboat/fixuid/pull/27)

## [0.4.1](https://github.com/boxboat/fixuid/releases/tag/v0.4.1) - 2020-04-28

### Features

- Add linux arm64 release: [#23](https://github.com/boxboat/fixuid/pull/23)

## [0.4](https://github.com/boxboat/fixuid/releases/tag/v0.4) - 2018-05-24

### Features

- Add quiet mode command-line flag `-q`: [#11](https://github.com/boxboat/fixuid/issues/11)

## [0.3](https://github.com/boxboat/fixuid/releases/tag/v0.3) - 2018-01-15

### Features

- Allow specifying paths to search: [#5](https://github.com/boxboat/fixuid/issues/5)

### Fixes

- Change Mount Detection to read /proc/mounts: [#7](https://github.com/boxboat/fixuid/issues/7)
- Handle errors from `lstat` and `filepath.readDirNames`: [#4](https://github.com/boxboat/fixuid/issues/4)

## [0.2](https://github.com/boxboat/fixuid/releases/tag/v0.2) - 2017-11-08

### Fixes

- Properly skip mounted files: [#3](https://github.com/boxboat/fixuid/pull/3)

## [0.1](https://github.com/boxboat/fixuid/releases/tag/v0.1) - 2017-07-18

- Initial Release
