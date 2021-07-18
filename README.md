# bugmenot-cli

*bugmenot-cli* is a command line tools to request logins that have been shared on [bugmenot.com](http://bugmenot.com)

## Install
```
go get -u github.com/guervild/bugmenot-cli
```

## Usage
The basic command require the *-domain* argument:
```
./bugmenot-cli -domain <your domain>
Results for <TRUNKED>:
+-----------------------+------------------------+------------------------+------------------+
|       USERNAME        |        PASSWORD        |         OTHER          |      RATING      |
+-----------------------+------------------------+------------------------+------------------+
| 2cgc5h@<TRUNKED>.jp   | b<TRUNKED>t            |                        | 79% success rate |
+-----------------------+------------------------+------------------------+------------------+
[...]
+-----------------------+------------------------+------------------------+------------------+
| hcp8hf@v<TRUNKED>.jp  | A<TRUNKED>a            |                        | 50% success rate |
+-----------------------+------------------------+------------------------+------------------+
```
By default the results are printed into a table.

You can use the *-json* argument to get the results in the JSON format:
```
./bugmenot-cli -domain <your domain> -json
```
Do not hesiste to combine your json output with the great tool [gron](https://github.com/tomnomnom/gron).

To filter the result and keep a minimum of success rate, you can use the *-filter <value>* option (default value is 0):
```
./bugmenot-cli -domain <your domain> -filter 51
Results for <TRUNKED>:
+---------------------+-------------+-------+------------------+
|      USERNAME       |  PASSWORD   | OTHER |      RATING      |
+---------------------+-------------+-------+------------------+
| 2cgc5h@<TRUNKED>.jp | b<TRUNKED>o |       | 79% success rate |
+---------------------+------------+-------+-------------------+
| pb5h@<TRUNKED>.jp   | A<TRUNKED>a |       | 64% success rate |
+---------------------+-------------+-------+------------------+
```