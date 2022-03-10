
try:
	snapcraft try
	snap try prime

sync:
	cp -r $$(ls | egrep -v '^prime') prime/

test:
	sudo edgex-snap-hooks.test \
		./ \
		./log \
		./snapctl