## solution_tester
This is a cli tool for testing your solutions against testcases fetched from various online judges(Codeforces, CodeChef, etc.). This tool was created for my use only, you can fork and work on it as per your needs.

### working
You start a *local judge* in your console, with the command `solution_tester judge`. This local judge will do the following, listens for problems from "Competative Companion" browser extension, at a time only one problem will be in his context. And listens for your submissions, you can submit your file using `solution_tester test <file_name> [--op-only]`, It compiles your submitted file using locally avialaible compiler and then executes it, also shows if there are any warnings/errors during compile time. The default behaviour is to judge your code's output against avilable output from testcases, the `--op-only` flag just shows the output of your code without any judgement useful for cases where multiple answers are true.

### usage
You can map some shortcut to the problem submission command in your preffered IDE. Whenever there is any contest or you want to practise locally just start the console, from browser send the problem to the console (TIP: just use the shortcut `Cmd + Shift + U`). And from your prefferef IDE send your submissions with the shortcuts you setup.

### note
* At a time only one problem will be in the context of the *local judge*, whenever you send new problems from browser, the context will be set to the new problem and further submissions will be tested against this new problem.
* you need to setup the solution_tester's port in the extension's options.
* a single local judge can be there at a single port

### env variables
* `PORT` : defaults to `12121`
* `COMPILE_COMMAND` : defaults to `g++ -DLOCAL -Wall -Wextra -Wconversion -Wshadow -Wfloat-equal -Wno-unused-includes -Wno-unused-const-variable -Wno-sign-conversion -O2 -std=c++23 \"%s\" -o \"%s\" -lstdc++exp`, the first %s will be replaced with path from your submission, and the second %s is handled by the judge, the compiled binaries are automatically cleaned up.

### building
`go build -o bin/solution_tester cmd`