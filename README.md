# Dron

Small shell command executor 

No more raw `bash` scripts!

## Example of config

File name `dron.yaml` or `dron.yml`

> All configs below is actual and work
> Do not use example file in repository - this is only for development 

### v1.0.0 [31.10.2020]

```yaml
commands:
  - name: up_www
    args:
      arg0: World
    commands:
      - echo Hello $arg0
```

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
