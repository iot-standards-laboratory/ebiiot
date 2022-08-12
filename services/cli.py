import os
import sys
import time

proto = 'tcp'

if len(sys.argv) < 2:
    print("enter runner type (run or build)")
    exit(-1)

if sys.argv[1] == 'run':
    cmd_server_run = 'go run main.go -server -proto {} & > server.log'.format(
        proto)
    time.sleep(0.5)
    cmd_client_run = 'go run main.go -proto {} {} > client.log'.format(
        proto, 'localhost:8080')

    os.system(cmd_server_run)  # server run
    os.system(cmd_client_run)   # client run

    os.system("pkill main")
else:
    os.system('rm ../utils/pa')
    os.system('go build -o pa main.go')
    os.system('mv pa ../utils/')
