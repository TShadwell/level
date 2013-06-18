gofmt -d -e -s \
	-r 'options -> UnderlyingOptions' \
	-r 'closer -> Closer' \
	-r 'database -> UnderlyingDatabase'\
	-r 'writeOptions -> UnderlyingWriteOptions'\
	-r 'readOptions -> UnderlyingReadOptions'\
	-r 'writeBatch -> UnderlyingWriteBatch'\
	-r 'cache -> UnderlyingCache'\
	.
