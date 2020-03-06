import random
import csv
import json
import networkx as nx
import matplotlib.pyplot as plt

class NetworkSim:
    def __init__( self, net_config:str ):

        with open( net_config ,'rt')as f:
            data = csv.reader(f)
            self.network = nx.Graph()
            
            regions_list        = []
            temp_lat_table      = {}
            for row in data:
                region = row[0].strip() 
                regions_list.append(region)
                temp_lat_table[ region ] = row[ 1: ]
            
            for region in regions_list:
                self.network.add_node( region )

            for base_region in regions_list:
                weighted_edges = []
                for i in range(0, len(regions_list)):
                    weighted_edges.append( (base_region, regions_list[i], float(temp_lat_table[region][i]) ) )

                self.network.add_weighted_edges_from( weighted_edges )

    def draw(self):
        nx.draw_networkx( self.network )
        plt.show()

    def minimum_span( self ):
        return nx.minimum_spanning_tree( self.network ) 

if __name__ == "__main__":
    sim                 = NetworkSim( 'latency_data/3.3.2020.csv' ) 
    optimal_multicast     = sim.minimum_span()

    for region in optimal_multicast:
        r_neighbors = []
        temp = optimal_multicast.neighbors( region )

    nx.draw_networkx( optimal_multicast )
    plt.show( )

    f = open("output/optimal_multicast.json", "w")

    json_output = {}
    narada      = [] 
    for region in optimal_multicast:
        narada.append( [region] + list( nx.neighbors(optimal_multicast, region) ) )
    json_output[ 'Narada' ] = narada

    f.write( json.dumps( json_output, indent = 4, sort_keys = True ) )

    