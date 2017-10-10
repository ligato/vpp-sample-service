#!/usr/bin/env bash
(
echo "open 0 5001"
sleep 1
echo "show ip fib"
sleep 1
echo "quit"
sleep 1
echo "exit"
) | telnet