CSV Runner
==========

Library and command line tool to use CSV input as data for a command line templating system


Usage for `csvrun`
------------------


### Command Line Template

```sh
$ csvrun -t 'echo $Name' names1.csv path/to/names2.csv
```

-or-

```sh
$ cat names.csv | csvrun -t 'echo $Name'
```


### Environment Template

```sh
$ csvrun -e HIDDEN_TEMPLATE names1.csv path/to/names2.csv
```

-or-

```sh
$ cat names.csv | csvrun -e HIDDEN_TEMPLATE
```


### File Template

```sh
$ csvrun -f template.txt names1.csv path/to/names2.csv
```

-or-

```sh
$ cat names.csv | csvrun -f template.txt
```

If no options are given `csvrun` will check for the environment variable `CSV_RUNNER_TEMPLATE`.

