# Dron

Small shell command executor

No more raw `bash` scripts!

## How to use

Create config in `dron.yaml` and run!  
Example of config see below

```bash
dron <command_name>
```

> Development only!  
> `dron -debug <command_name>`

## Installation

### ArchLinux

You can install it from [AUR](https://aur.archlinux.org/packages/dron/)

### Manual

Requirements:

* [Go](https://golang.org/)

```bash
git clone https://github.com/DizoftTeam/dron.git
cd dron
go build -o dron main.go
# For all users
sudo ln -s /home/<user>/path/to/dron /usr/local/bin/dron
```

## Example of config

File name `dron.yaml` or `dron.yml`

> All configs below is actual and work  
> Do not use example file in a repository - this is only for development

### v1.1.0 [01.11.2020]

```yaml
commands:
  - name: up_www
    args:
      arg0: World
      arg1: $env(APP_ENV)
    commands:
      - echo Hello $arg0
      - echo env_param_APP_ENV $arg1
      - echo "arg0 $arg0 with quotes on end $arg1"
```

* Add suppor for `$env` command
  * If env param not find in system - error will be provided
* Bug fix — double quotes not removing after argument name
* If `.env` exist — it will be load automatically

### v1.0.0 [31.10.2020]

```yaml
commands:
  - name: up_www
    args:
      arg0: World
    commands:
      - echo Hello $arg0
```
