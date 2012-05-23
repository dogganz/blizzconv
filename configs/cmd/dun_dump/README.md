dun_dump
========

dun_dump is a tool for constructing dungeons, based on the information retrieved
from a given DUN file, and storing these dungeons as png images.

Installation
------------

    $ go get github.com/mewkiz/blizzconv/configs/cmd/dun_dump
    $ go install github.com/mewkiz/blizzconv/configs/cmd/dun_dump

Usage
-----

    $ mkdir blizzdump/
    $ cd blizzdump/
    $ ln -s /path/to/extracted/diabdat_mpq mpqdump
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/configs/dunconf/dun.ini
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/images/imgconf/cel.ini
    $ ln -s $GOPATH/src/github.com/mewkiz/blizzconv/mpq/mpq.ini
    $ dun_dump -a