Instructions to run on Linux console shell(bash):
=========================================

1) Install GO on Linux machine and set `GOROOT` variable  to install folder in `~/.bashrc` file
    
* Note: This step is not needed, if GO is installed in standard path `/usr/local/go` folder

2) Append below entries in `~/.bashrc` file, to set `GOPATH` & `PATH` variables:

        export PATH=$PATH:/usr/local/go/bin
        export GOPATH=/home/{user}/golib
        export PATH=$PATH:$GOPATH/bin
        # First segment of GOPATH is used by "go get" command
        # All segments of GOPATH are used for searching symbols of source code
        export GOPATH=$GOPATH:/home/{user}/code

3) Run command `source ~/.bashrc`

4) Run command `mkdir -p /home/{login-user}/code/pkg`

5) Run command `mkdir -p /home/{login-user}/code/bin/`

6) Run command `mkdir -p /home/{login-user}/code/src/github.com/shamhub`

7) Run command `cd /home/{login-user}/code/src/github.com/shamhub` and then run `git clone https://github.com/shamhub/threadapp` 
    
* Note: `../threadapp` is a command package

8) Run command `cd threadapp`

9) Run command: `go get "github.com/nicholasjackson/env"` to install a third party

10) Run command `export BATCH_SIZE=<positive_integer>`. Default batch size is `100` - Optional step(but recommended)

11) Run command `export MAX_OBJECTS_TO_PRINT=<positive_integer>` for any input value. Default value is `50000` - Optional step(but recommended)

12) Run command `make install` will generate 64 bit binary(as `~/code/bin/threadapp`), if GO compiler is 64 bit from Step 1

13) Run command  `~/code/bin/threadapp`


* Note: Output sequence is not more than first 50000 objects, for any input object



++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Problem Statement
=================

You are receiving n objects in a random order, and you need to print them to stdout correctly ordered by
sequence number.

The sequence numbers start from `0` (zero) and you have to wait until you get a complete, unbroken sequence
batch of `j` objects before you output them.
> You have to process all objects without loss.
> The program should exit once it completes outputting the first 50000 objects
> Batch size j = 100

The object is defined as such:

            {
            "id" : "object_id", // object ID (string)
            "seq" : 0, // object sequence number (int64, 0-49999)
            "data" : "" // []bytes
            }

Discussion Points
=================
- What can be added to the object structure to enable more robust architectures?
- What compromises can be made to alleviate memory pressure?
- What can be done to enable concurrent “processing” of the input data?


Example Output Statement
========================

                Step                Input Value                Output State j = 1                  Output state j = 3
                0                       6
                1                       0                           0
                2                       4                           0
                3                       2                           0
                4                       1                           0,1,2                               0,1,2
                5                       3                           0,1,2,3,4                           0,1,2
                6                       9                           0,1,2,3,4                           0,1,2
                7                       5                           0,1,2,3,4,5,6                       0,1,2,3,4,5
                                etc.