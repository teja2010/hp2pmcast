import csv
import os
import sys

def main():
    '''
    Takes the packet count as an argument for and writes output to a file so that it
    can be run and have results collected later on.
    '''

    if len(sys.argv) < 2:
        print( "Must pass a packet count." ) 
        return 

    packet_count = int( sys.argv[1] )

    with open( 'endpoints.csv' ) as csv_file:
        csv_reader = csv.reader(csv_file, delimiter=',')
        line_count = 0
        for row in csv_reader:
            if line_count == 0:
                line_count += 1
            else:
                region_name         = row[ 0 ]
                region_endpoint_ip  = row[ 1 ]
                print( f"Pinging region:\t{ region_name } at \t{ region_endpoint_ip } and outputing to file 'to_{region_name}'." )
                
                os.system( f'ping -c {packet_count} {region_endpoint_ip} > to_{region_name}' )

                line_count += 1

        print(f'Pinged {line_count - 1} hosts.')

if __name__ == "__main__":
    main()