import json
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd


# Specify the path to your JSON file
json_file = "metric/results-client-150-messages-150.json"


# Read JSON data from the file
with open(json_file, 'r') as f:
    json_data = json.load(f)

# Initialize lists to store data
clients = []
avg_times = []
min_times = []
max_times = []

# Process each client's data
for client_id, client_data in json_data.items():
    clients.append(client_id)
    avg_times.append(client_data['avg'])
    min_times.append(client_data['min'])
    max_times.append(client_data['max'])

    # Create a DataFrame for further analysis
df = pd.DataFrame({
    'Client': clients,
    'Average Time': avg_times,
    'Min Time': min_times,
    'Max Time': max_times
})

# Display the DataFrame (optional)
print(df)

### Perform Analytics and Visualizations ###

# Example 1: Bar Plot of Average Message Times
plt.figure(figsize=(10, 6))
plt.bar(df['Client'], df['Average Time'], color='skyblue')
plt.xlabel('Clients')
plt.ylabel('Average Message Time')
plt.title('Average Message Time per Client')
plt.xticks(rotation=45)
plt.tight_layout()
plt.show()

# Example 2: Box Plot of Message Times Distribution
plt.figure(figsize=(10, 6))
sns.boxplot(x='Client', y='Average Time', data=df, palette='pastel')
plt.xlabel('Clients')
plt.ylabel('Average Message Time')
plt.title('Distribution of Average Message Time per Client')
plt.xticks(rotation=45)
plt.tight_layout()
plt.show()