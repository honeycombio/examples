#!/bin/sh

while true;
do
	for element in earth fire wind water;
	do
		curl "localhost:8000/echo/$element"
	done
done
