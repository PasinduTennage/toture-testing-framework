import argparse
from benchmark.local import Local

# Create an ArgumentParser object
parser = argparse.ArgumentParser()


# server_name, threshold, test_time, view_time, benchmark_type, epoch_time

# Add command line arguments
parser.add_argument('--server_name', type=str, help='name of the server')
parser.add_argument('--threshold',   type=int, help='number of concurrently attacked nodes')
parser.add_argument('--test_time', type=int, help='test time in seconds')
parser.add_argument('--view_time', type=int, help='view_time in milliseconds')
parser.add_argument('--benchmark_type', type=str, help='type of benchmark')
parser.add_argument('--epoch_time', type=int, help='epoch time in milliseconds')

# Parse the command line arguments
args = parser.parse_args()

# Instantiate the LocalBenchmark class with the command line arguments
benchmark = Local(args.server_name, args.threshold, args.test_time, args.view_time, args.benchmark_type, args.epoch_time)
benchmark.attack()