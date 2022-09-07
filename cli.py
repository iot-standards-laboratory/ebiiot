# python runner.py -t config/topo/low-delay -x config/xp/low-tcp
import os
import sys
import template

if len(sys.argv) < 2 :
    print("enter cli type (run or clean)")
    exit(-1)

if sys.argv[1] == 'run':
    # exp parameter
    exp = 'http'
    proto = 'hybrid'
    clients = 50
    messages = 1
    msgsize = 100000

    # topo parameter
    delay = 100
    queueSize = 30
    bandwidth = 4
    loss = 2.0

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
    cmd = 'python runner.py -t {} -x {}'.format(topoFile, xpFile)
    os.system(cmd)

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
    deleteRecursively('ssl-key.log')



