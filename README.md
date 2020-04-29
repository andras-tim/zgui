##### This is a [WIP](https://en.wikipedia.org/wiki/WIP) in alpha stage.

# zgui
Linux GUI viewer for ZFS pool, dataset and host storage.

zgui used libzfs directly and not ZFS command line tools.

## Screenshot

<img src="screenshot/dataset.png" width="400">
<img src="screenshot/pool.png" width="400">
<img src="screenshot/storage.png" width="400">

## Installation

```$ go get gitlab.com/beteras/zgui```

and to execute

```$ zgui```

## Prerequisite
- Debian Buster with packages
  - libgtk-3-dev
  - libzfslinux-dev
- Go 1.13
- ZFS ([OpenZFS](https://github.com/openzfs/zfs)) 0.8.3
- GTK 3.24.5

## Others Go project used
- [libzfs](https://github.com/bicomsystems/go-libzfs) - Go bindings for libzfs 
- [gotk3](https://github.com/gotk3/gotk3) - Go bindings for GTK+3
- [ghw](https://github.com/jaypipes/ghw) - Golang HardWare discovery/inspection library
- [humanize](https://github.com/dustin/go-humanize) - Humane Units

## License
- [GPL v3](LICENSE)

### Icons source
- https://www.flaticon.com/packs/servers-database-3
  - By [Those Icons](https://www.flaticon.com/authors/those-icons)
  - [Free for commercial use WITH ATTRIBUTION license](https://profile.flaticon.com/license/free)
- https://www.flaticon.com/packs/basic-ui-3
  - By [Pixel perfect](https://www.flaticon.com/authors/pixel-perfect)
  - [Free for commercial use WITH ATTRIBUTION license](https://profile.flaticon.com/license/free)
- https://www.flaticon.com/packs/computer-hardware-31
  - By [Smashicons](https://www.flaticon.com/authors/smashicons)
  - [Free for commercial use WITH ATTRIBUTION license](https://profile.flaticon.com/license/free)
- https://thenounproject.com/term/raid-drive/1464982/
  - By [Ben Davis](https://thenounproject.com/smashicons/)
  - [Creative Commons CCBY](https://creativecommons.org/licenses/by/3.0/us/legalcode)
