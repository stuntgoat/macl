`MACL`: *Make Arbitrarily long Contigious Lines*


You will need a working Golang installation to build.

`make`

Run tests.

`make test`

Run test coverage report.

`make cover`

Run HTML coverage report (will attempt to open a web browser).

`make coverhtml`


To start the server on port 80 (need sudo privileges to start on this port).

`$ sudo ./macl -port=80`

Help

    $ ./macl -h

    Usage of ./macl:
      -api_prefix string
            api URL prefix (default "game")
      -board_length int
            board length (default 4)
      -board_width int
            board width (default 4)
      -consecutive_length int
            consecutive line length required for a win (default 4)
      -log_path string
            logging path (default "macl.log")
      -num_players int
            required number of players (default 2)
      -port int
            server port (default 8080)
