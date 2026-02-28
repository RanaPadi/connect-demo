import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import re
import csv

# Function to extract benchmark data from log files
def extract_benchmark_data(file_path, benchmark_pattern):
    values = []
    with open(file_path, 'r') as file:
        for line in file:
            match = re.search(benchmark_pattern, line)
            if match:
                values.append(float(match.group(1)))
    return values

# Function to calculate statistics
def calculate_stats(data):
    return {
        "min": pd.Series(data).min(),
        "max": pd.Series(data).max(),
        "mean": pd.Series(data).mean(),
        "median": pd.Series(data).median(),
        "std_dev": pd.Series(data).std()
    }

# Function to add statistics text below the graph within each subplot
def add_stats_below(ax, data):
    stats = calculate_stats(data)
    stats_text = (f"Min: {stats['min']:.3f}\n"
                  f"Max: {stats['max']:.3f}\n"
                  f"Mean: {stats['mean']:.3f}\n"
                  f"Median: {stats['median']:.3f}\n"
                  f"Std Dev: {stats['std_dev']:.3f}")
    ax.text(0.5, -0.1, stats_text, fontsize=10, ha='center', va='top', transform=ax.transAxes)

# Function to write statistics to a CSV file
def write_stats_to_csv(stats, benchmark_name):
    with open('benchmark_stats.csv', mode='a', newline='') as file:
        writer = csv.writer(file)
        writer.writerow([benchmark_name, stats['min'], stats['max'], stats['mean'], stats['median'], stats['std_dev']])

# Path to log file
log_file_path = "./benchmark.log"

# Paths to log files
file_paths = {
    "benchmark_3": log_file_path,
    "benchmark_4": "./aiv_eval.log",
    "benchmark_7": log_file_path,
    "benchmark_8": log_file_path,
}

# Extract data for each benchmark
benchmark_3_values = extract_benchmark_data(file_paths['benchmark_3'], r'Benchmark 3.*\):.* (\d+\.\d+)ms')
benchmark_4_values = extract_benchmark_data(file_paths['benchmark_4'], r'\t (\d+\.\d+)$')
benchmark_7_values = extract_benchmark_data(file_paths['benchmark_7'], r'Benchmark \d+: (\d+\.\d+)ms')
benchmark_8_vc1_vc2_values = extract_benchmark_data(file_paths['benchmark_8'], r'Benchmark 8.*VC1->VC2.* (\d+\.\d+)ms')
benchmark_8_vc2_vc1_values = extract_benchmark_data(file_paths['benchmark_8'], r'Benchmark 8.*VC2->VC1.* (\d+\.\d+)ms')
benchmark_8_combined_values = extract_benchmark_data(file_paths['benchmark_8'], r'Benchmark 8.* (\d+\.\d+)ms')
# Convert Benchmark 4 values from seconds to milliseconds
benchmark_4_values_ms = [value * 1000 for value in benchmark_4_values]

# Initialize CSV file with headers
with open('benchmark_stats.csv', mode='w', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(["Benchmark", "Min", "Max", "Mean", "Median", "Std Dev"])

# Benchmark 3: Plot with all values
plt.figure(figsize=(12, 10))

ax1 = plt.subplot(1, 1, 1)
ax1.boxplot(benchmark_3_values, vert=True)
ax1.set_title('Benchmark 3 - All Values')
ax1.set_ylabel('Time (ms)')
ax1.grid(True)
add_stats_below(ax1, benchmark_3_values)
write_stats_to_csv(calculate_stats(benchmark_3_values), "Benchmark 3 - All Values")

plt.tight_layout()
plt.savefig('benchmark_3.png')
plt.close()

# Benchmark 4: Plot with all values in milliseconds
plt.figure(figsize=(12, 10))

ax1 = plt.subplot(1, 1, 1)
ax1.boxplot(benchmark_4_values_ms, vert=True)
ax1.set_title('Benchmark 4 - All Values')
ax1.set_ylabel('Time (ms)')
ax1.grid(True)
add_stats_below(ax1, benchmark_4_values_ms)
write_stats_to_csv(calculate_stats(benchmark_4_values_ms), "Benchmark 4 - All Values (in ms)")

plt.tight_layout()
plt.savefig('benchmark_4.png')
plt.close()

# Benchmark 7: Plot with statistics
plt.figure(figsize=(6, 12))

ax = plt.subplot(1, 1, 1)
ax.boxplot(benchmark_7_values, vert=True)
ax.set_title('Benchmark 6 All Values')
ax.set_ylabel('Time (ms)')
ax.grid(True)
add_stats_below(ax, benchmark_7_values)
write_stats_to_csv(calculate_stats(benchmark_7_values), "Benchmark 6")

plt.tight_layout()
plt.savefig('benchmark_6.png')
plt.close()

# Benchmark 8: Plot with statistics (VC1->VC2, VC2->VC1, Combined)
plt.figure(figsize=(18, 12))

ax1 = plt.subplot(1, 3, 1)
ax1.boxplot(benchmark_8_vc1_vc2_values, vert=True)
ax1.set_title('Benchmark 7a VC1->VC2')
ax1.set_ylabel('Time (ms)')
ax1.grid(True)
add_stats_below(ax1, benchmark_8_vc1_vc2_values)
write_stats_to_csv(calculate_stats(benchmark_8_vc1_vc2_values), "Benchmark 7a VC1->VC2")

ax2 = plt.subplot(1, 3, 2)
ax2.boxplot(benchmark_8_vc2_vc1_values, vert=True)
ax2.set_title('Benchmark 7b VC2->VC1')
ax2.grid(True)
add_stats_below(ax2, benchmark_8_vc2_vc1_values)
write_stats_to_csv(calculate_stats(benchmark_8_vc2_vc1_values), "Benchmark 7b VC2->VC1")

ax3 = plt.subplot(1, 3, 3)
ax3.boxplot(benchmark_8_combined_values, vert=True)
ax3.set_title('Benchmark 7c Combined')
ax3.grid(True)
add_stats_below(ax3, benchmark_8_combined_values)
write_stats_to_csv(calculate_stats(benchmark_8_combined_values), "Benchmark 7 Combined")

plt.tight_layout()
plt.savefig('benchmark_7.png')
plt.close()
