import subprocess
import matplotlib.pyplot as plt

num_threads_list = [4, 8, 16, 32, 64]

results = []
serial_time = 0

for num_threads in num_threads_list:
    cmd_args = ["go", "run", "main.go", str(num_threads)]

    result = subprocess.run(cmd_args, stdout=subprocess.PIPE)
    output = result.stdout.decode().strip()
    print(output)
    [time_taken_p, time_taken_s] = (float(i[:-1]) for i in output.split("\n"))
    serial_time += time_taken_s
    results.append(time_taken_p)

serial_time /= len(num_threads_list)
plt.plot(num_threads_list, results, marker='o', label='Parallel')
plt.axhline(y=serial_time, color='r', linestyle='-', label='Serial')

plt.xlabel('Number of Threads')
plt.ylabel('Time Taken (Seconds)')
plt.title('Time Taken vs. Number of Threads')

plt.legend()
plt.show()
