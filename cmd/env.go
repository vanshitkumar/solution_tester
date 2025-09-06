package main

import "os"

var(
	PORT string
	COMPILE_COMMAND string
)

func init(){
	PORT = os.Getenv("PORT")
	COMPILE_COMMAND = os.Getenv("COMPILE_COMMAND")
	if PORT == "" {PORT = "12121"}
	if COMPILE_COMMAND == "" {COMPILE_COMMAND = "g++ -DLOCAL -Wall -Wextra -Wconversion -Wshadow -Wfloat-equal -Wno-unused-includes -Wno-unused-const-variable -Wno-sign-conversion -O2 -std=c++23 %s -o %s -lstdc++exp"}
}