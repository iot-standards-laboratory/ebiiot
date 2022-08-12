from core.experiment import ExperimentParameter, RandomFileExperiment, RandomFileParameter
import os
import threading


class PerformanceAnalysisExp(RandomFileExperiment):
    # APPNAME = "qcoap-uni"
    PROTO = "tcp"
    CLIENTS = "10"
    MESSAGES = "10"
    MSGSIZE = "100"

    EXPTYPE = "performance-analysis"
    NAME = "performance-analysis"
    SERVER_LOG = "server-log.txt"
    CLIENT_LOG = "client-log.txt"
    PING_OUTPUT = "ping.log"

    def __init__(self, experiment_parameter_filename, topo, topo_config):
        # Just rely on RandomFileExperiment
        super(PerformanceAnalysisExp, self).__init__(
            experiment_parameter_filename, topo, topo_config)

        self.PROTO = self.experiment_parameter.get("proto")
        self.CLIENTS = self.experiment_parameter.get("clients")
        self.MESSAGES = self.experiment_parameter.get("messages")
        self.MSGSIZE = self.experiment_parameter.get("msgsize")
        print("PROTO: ", self.PROTO)

    def load_parameters(self):
        # Just rely on RandomFileExperiment
        super(PerformanceAnalysisExp, self).load_parameters()

    def prepare(self):
        super(PerformanceAnalysisExp, self).prepare()
        self.topo.command_to(self.topo_config.client, "rm " +
                             PerformanceAnalysisExp.CLIENT_LOG)
        self.topo.command_to(self.topo_config.server, "rm " +
                             PerformanceAnalysisExp.SERVER_LOG)

    def getServerCmd(self):
        s = "{}/../utils/pa -server -proto {} >> server-log.txt &".format(
            os.path.dirname(os.path.abspath(__file__)),
            self.PROTO,
        )
        # s = "/home/mininet/pugit/sample/minitopo/utils/server & > {}".format(
        #     BASICQUIC.SERVER_LOG)

        print(s)
        return s

    def getClientCmd(self):
        s = "{}/../utils/pa -proto {} -clients {} -messages {} -size {} {}:8080 >> client-log.txt".format(
            os.path.dirname(os.path.abspath(__file__)),
            self.PROTO,
            self.CLIENTS,
            self.MESSAGES,
            self.MSGSIZE,
            self.topo_config.get_server_ip(0),
        )
        print(s)

        return s

    def clean(self):
        super(PerformanceAnalysisExp, self).clean()

    def checkNetwork(self, entity):
        self.topo.command_to(entity,
                             "route -n > route.out")
        self.topo.command_to(entity,
                             "ip addr > ipaddr.out")

    def run(self):

        self.topo.command_to(self.topo_config.servers[0],
                             "netstat -sn > netstat_server_before")
        self.topo.command_to(self.topo_config.router,
                             "netstat -sn > netstat_router_before")

        # self.topo.command_to(self.topo_config.clients[0],
        #                      "route -n > route.out")
        # self.topo.command_to(self.topo_config.clients[0],
        #                      "ip addr > ipaddr.out")

        cmd = self.getServerCmd()
        self.topo.command_to(self.topo_config.servers[0], cmd)
        print("Waiting for the server to run")

        self.topo.command_to(self.topo_config.clients[0], "sleep 1")
        cmd = self.getClientCmd()
        self.topo.command_to(self.topo_config.clients[0], cmd)

        self.topo.command_to(self.topo_config.servers[0],
                             "netstat -sn > netstat_server_after")
        self.topo.command_to(self.topo_config.router,
                             "netstat -sn > netstat_router_after")
        self.topo.command_to(self.topo_config.clients[0], "sleep 2")
        self.topo.command_to(self.topo_config.servers[0],
                             "pkill -f pa")
