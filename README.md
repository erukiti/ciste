# ciste

ciste is tiny CI and PaaS integration repository server.

## make

```
$ git clone git@github.com:erukiti/ciste.git
$ cd ciste
$ git submodule update --init --recursive
$ make get
$ make
```

## development with vagrant

```
$ vagrant up
$ vagrant ssh
```

```
$ sudo cp ciste /usr/local/bin/
$ sudo su - git
$ ciste setup
$ ciste pubkey <username> <ssh public key>
$ ciste server &
```

## development with docker-machine

* create user 'git'

```
$ sudo cp ciste /usr/local/bin/
$ sudo su - git
$ docker-machine create --driver virtualbox dev
$ ciste setup
$ ciste pubkey <username> <ssh public key>
$ ciste server &
$ open http://$(docker-machine ip):3000/
```

