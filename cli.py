# python runner.py -t config/topo/low-delay -x config/xp/low-tcp
import os
import sys
import template

if len(sys.argv) < 2 :
    print("enter cli type (run or clean)")
    exit(-1)

def experiment(numOfTrials, exp, proto, clients, messages, msgsize, delay, queueSize, bandwidth, loss):
    # make experimentation and topology config file at config/xp and config/topo respectively. 
    xpTmpl = template.xp_template.xp
    topoTmpl = template.topology_template.topo

    # create xpFile 
    xpFile = './config/xp/{}_{}_{}_{}_{}'.format(
        exp, proto, clients, messages, msgsize)
    if not os.path.exists(xpFile):
        f = open(
            xpFile,
            'w+',
        )
        f.write(xpTmpl.format(exp, proto, clients, messages, msgsize))
        f.close()
    else:
        print('file is already exist')


    # create topoFile
    topoFile = './config/topo/topo-{}_{}_{}_{}'.format(delay, queueSize, bandwidth, loss)

    if not os.path.exists(topoFile):
        f = open(topoFile, 'w+')
        f.write(topoTmpl.format(2, delay, queueSize, bandwidth, loss))
        f.close()
    else: 
        print('file is already exist')

    # run experimentation
    for i in range(numOfTrials):
        cmd = 'python runner.py -t {} -x {}'.format(topoFile, xpFile)
        os.system(cmd)
        with open('./netstat_router_after', 'r') as stream:
            for line in stream:
                if line.startswith('    InOctets: '):
                    l = line[len('    InOctets: '):].rstrip()
                    os.system("echo {} >> netstat".format(l))

if sys.argv[1] == 'run':
    numOfTrials = 5
    # exp parameter
    exps = ['http']
    protos = ['tcp', 'quic', 'hybrid']
    clients = [50]
    messages = [1]
    msgsizes = [100000]

    # topo parameter
    delay = [100]
    queueSize = [30]
    bandwidth = [4]
    loss = [0, 1.0, 2.0]

    if len(sys.argv) >= 3:
        numOfTrials = int(sys.argv[2])

    for e in exps:
        for p in protos:
            for c in clients:
                for m in messages:
                    for ms in msgsizes:
                        for d in delay:
                            for q in queueSize:
                                for b in bandwidth:
                                    for l in loss:
                                        experiment(numOfTrials, e, p, c, m, ms, d, q, b, l)
                                        
                                        dist = './dist/{}/{}'.format(e, p)
                                        if not os.path.exists(dist):
                                            os.makedirs(dist)
                                        os.system("mv atd.out {}/atd-c:{}_ms:{}_d:{}_b:{}_l:{}".format(dist, c, ms, d, b, l))
                                        os.system("mv netstat {}/volume-c:{}_ms:{}_d:{}_b:{}_l:{}".format(dist, c, ms, d, b, l))

    

elif sys.argv[1] == 'clean':
    def deleteRecursively(fname):
        for root, _, files in os.walk('.'):
            if root[:3] == "./.":
                continue
            for f in files:
                if f == fname:
                    os.remove(os.path.join(root, f))
                    # print(os.listdir('.'))
    
    os.system('rm netstat*')
    os.system('rm *.pcap')
    os.system('rm server-log.txt')
    os.system('rm client-log.txt')
    os.system('rm *.log')
    os.system('rm *.out')
    os.system('rm -r config/topo/*')
    os.system('rm -r config/xp/*')
    deleteRecursively('ssl-key.log')



