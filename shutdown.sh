#!/bin/bash

trap 'echo "Received shutdown signal"; exit 0' SIGTERM

/bin/gen_sh &
wait $!
