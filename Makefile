BIN_DIR=`pwd`/bin
build:
	make -C json2csv
	make -C oneline
	make -C sqltable2csv
