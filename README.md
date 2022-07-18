CSV Runner v1.0.0
=================

Library and command line tool to use CSV input as data for a command line templating and running system


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


ToDo
----

* [x] Create `csvrun`
	* [x] Implement loading CSVs from command line arguments
	* [x] Implement reading CSV data from `stdin`
	* [x] Read template string from a file
	* [x] Read template string from an environment variable
	* [x] Read template string from command line
	* [ ] Storage of templates run prior for quick retrieval
	* [ ] Option to present most recent templates to choose from
* [ ] Create `CSV Runner` GUI
	* [ ] Choose CSVs to run by file picker
	* [ ] Choose CSVs to run by directory
	* [ ] Choose CSVs to run by directory with pattern
	* [ ] Choose template by file picker
	* [ ] Choose template by named environment variable
	* [ ] Input template via text field
	* [ ] Storage of templates run prior for quick retrieval
	* [ ] Template chooser
	* [ ] Template storage labeling and grouping


