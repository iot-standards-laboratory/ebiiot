import os


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
