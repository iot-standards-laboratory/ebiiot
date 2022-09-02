# python runner.py -t config/topo/low-delay -x config/xp/low-tcp
import os
import sys
import template

# exp parameter
exp = 'simple'
proto = 'tcp'
clients = '60'
messages = '100'
msgsize = '1000'

tmpl = template.simple_template.simpleXp
xpFile = './config/xp/{}_{}_{}_{}_{}'.format(
    exp, proto, clients, messages, msgsize)

if not os.path.exists(xpFile):
    f = open(
        xpFile,
        'w+',
    )
    f.write(tmpl.format(proto, clients, messages, msgsize))
    f.close()
else:
    print('file is already exist')

cmd = 'python runner.py -t config/topo/low-delay -x {}'.format(xpFile)
os.system(cmd)
