# fspop
Automate the creation of file structures using custom templates.

fspop => __f__(ile) __s__(tructure) __pop__(ulate)


## Defining a Structure File
Structure files are written in `.yaml` and require two things. (Use `fspop init` to create one for you)
1. name
2. structure

`fspop` will run though the entire structure and create a directory for each item

```yaml
name: media

structure:
    - games
    - music
    - photos:
        - personal
        - family
```


## Commands
| Command  	| Description                             	|
|----------	|-----------------------------------------	|
| `capture` | Capture existing directory structure      |
| `deploy` 	| Creates file structure from config file 	|
| `init`   	| Creates a new structure config file     	|
| `ls` 	    | List potential structure files            |


## Usage

### Create a new structure file
The fastest way to create a new structure is to use the `init` command.

```bash
$ fspop init STRUCTURE_NAME
```

### Deploy a structure
Deploying a structure is as simple as calling `deploy` and giving the file name of the structure

```bash
$ fspop deploy STRUCTURE_NAME
```


## Original `fspop`
Originally I wrote `fspop` in NodeJS but have since learned Go. This repo is a full rewrite of `fspop` in Go.

You can find the original version here (although I wouldn't recommend using it for anything other than testing / curiosity):
- [github.com/hmerritt/fspop-nodejs](https://github.com/hmerritt/fspop-nodejs)



<br>

## License
[Apache-2.0 License](LICENSE)
