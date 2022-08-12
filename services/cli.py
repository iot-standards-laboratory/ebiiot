import os
import sys
import time

exp = 'simple'
proto = 'quic'

if len(sys.argv) < 2:
    print("enter runner type (run or build)")
    exit(-1)

if sys.argv[1] == 'run':
    cmd_server_run = 'go run main.go -server -exp {} -proto {} & > server.log'.format(
        exp,
        proto,
    )
    
    cmd_client_run = 'go run main.go -exp {} -proto {} {} > client.log'.format(
        exp, proto, 'localhost:8080')

    os.system(cmd_server_run)  # server run
    time.sleep(1)
    os.system(cmd_client_run)   # client run

    os.system("pkill main")
else:
    os.system('rm ../utils/pa')
    os.system('go build -o pa main.go')
    os.system('mv pa ../utils/')
