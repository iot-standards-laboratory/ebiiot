# python runner.py -t config/topo/low-delay -x config/xp/low-tcp
import os
import sys
import template

# exp parameter
exp = 'simple'
proto = 'quic'
clients = '60'
messages = '100'
msgsize = '1000'

tmpl = template.simple_template.simpleXp
xpFile = './config/xp/{}_{}_{}_{}_{}'.format(
    exp, proto, clients, messages, msgsize)
f = open(
    xpFile,
    'w+',
)

f.write(tmpl.format(proto, clients, messages, msgsize))
f.close()

cmd = 'python runner.py -t config/topo/low-delay -x {}'.format(xpFile)

os.system(cmd)
