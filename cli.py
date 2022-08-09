# python runner.py -t config/topo/low-delay -x config/xp/low-tcp
import os
import sys


exp = 'tcp'

cmd = 'python runner.py -t config/topo/low-delay -x config/xp/echo_tcp_100'
os.system(cmd)
