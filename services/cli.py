import os
import sys
import time

exp = 'simple'
proto = 'quic'
clients = '10'
trials = '5'
objs = '20'
msgsize = '1000'

if len(sys.argv) < 2:
    print("enter runner type (run or build)")
    exit(-1)

if sys.argv[1] == 'run':
    # run tcpdump
    os.system('sudo tcpdump -i lo -w dump.pcap &')
    time.sleep(2)
    cmd_server_run = 'go run main.go -server -exp {} -proto {} & '.format(
        exp,
        proto,
    )
    
    cmd_client_run = 'go run main.go -exp {} -proto {} -clients {} -trials {} -objs {} -size {} {} > client.log'.format(
        exp, proto, clients, trials, objs, msgsize, 'localhost:8080')

    os.system(cmd_server_run)  # server run
    time.sleep(1)
    os.system(cmd_client_run)   # client run
    time.sleep(2)
    os.system("pkill main")
    os.system('sudo pkill tcpdump')
elif sys.argv[1] == 'clean':
    print('clean')
    os.system('rm *.log')
    os.system('rm dump.pcap -y')
else:
    print('start build')
    os.system('rm ../utils/pa')
    os.system('go build -o pa main.go')
    os.system('mv pa ../utils/')
