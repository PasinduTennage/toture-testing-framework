import sys

def getQuePaxaPerformance(root, initClient, numClients):
    throughputs = []
    medians = []
    ninety9s = []
    errors = []
    for cl in list(range(initClient, initClient + numClients, 1)):
        file_name = root + str(cl) + ".log"
        # print(file_name + "\n")
        try:
            f = open(file_name, 'r')
        except OSError:
            sys.exit("Error in " + file_name + "\n")

        with f:
            content = f.readlines()
        if len(content) < 10:
            sys.exit("Error in " + file_name + "\n")

        if content[0].strip().split(" ")[0] == "Warning:":
            content = content[1:]
        if not (content[6].strip().split(" ")[0] == "Throughput" and content[7].strip().split(" ")[0] == "Median"):
            sys.exit("Error in " + file_name + "\n")

        throughputs.append(float(content[6].strip().split(" ")[2]))
        medians.append(float(content[7].strip().split(" ")[3]))
        ninety9s.append(float(content[8].strip().split(" ")[4]))
        errors.append(float(content[9].strip().split(" ")[3]))

    return [sum(throughputs), sum(medians) / numClients, sum(ninety9s) / numClients, sum(errors)]


# main fuction calls the getQuePaxaPerformance function using the cmd line arguments
def main():
    if len(sys.argv) != 4:
        sys.exit("Usage: python3 quepaxa_summary.py <root> <initClient> <numClients>\n")

    root = sys.argv[1]
    initClient = int(sys.argv[2])
    numClients = int(sys.argv[3])

    [throughput, median, ninety9, error] = getQuePaxaPerformance(root, initClient, numClients)
    print("Throughput: " + str(throughput))
    print("Median: " + str(median))
    print("99th percentile: " + str(ninety9))
    print("Error rate: " + str(error))


# call the main function
if __name__ == "__main__":
    main()
