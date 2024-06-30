import json
import pandas as pd
import plotly.graph_objs as go
from jinja2 import Environment, FileSystemLoader
import os

# List of JSON files
json_files = [
    "metric/results-client-go-500-messages-500.json",
    "metric/results-client-go-1500-messages-500.json",
    "metric/results-client-go-10000-messages-500.json",
    "metric/results-client-typescript-150-messages-150.json",
    "metric/results-client-typescript-100-messages-100.json",
]

# Dictionary to store report file paths and details
reports = {}

for json_file in json_files:
    # Extract technology name from file path
    pathname = os.path.basename(json_file)
    technology = pathname.split('-')[2]  # Assumes file structure is consistent
    active_client = pathname.split("-")[3]
    message_per_client = pathname.split("-")[5].split(".")[0]

    # Read JSON data from the file
    with open(json_file, 'r', encoding='utf-8') as f:
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

    # Generate HTML report file path
    output_dir = 'dist'
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    report_file = os.path.join(output_dir, f'report-{technology}-{active_client}-{message_per_client}.html')

    # Generate HTML using Jinja2 templating
    env = Environment(loader=FileSystemLoader('templates'))
    template = env.get_template('report_template.html')

    # Prepare data for Plotly graphs
    bar_trace = go.Bar(
        x=df['Client'],
        y=df['Average Time'],
        name='Average Time',
        marker=dict(color='skyblue')
    )

    box_trace = go.Box(
        x=df['Client'],
        y=df['Average Time'],
        name='Distribution',
        marker=dict(color='lightseagreen')
    )

    # Scatter Plot for Min Time vs. Max Time
    scatter_trace = go.Scatter(
        x=df['Min Time'],
        y=df['Max Time'],
        mode='markers',
        marker=dict(color='salmon'),
        name='Min Time vs. Max Time'
    )

    # Violin Plot for Distribution of Min Time and Max Time
    violin_trace = go.Violin(
        y=df['Min Time'],
        box_visible=True,
        line_color='black',
        name='Min Time Distribution',
        marker=dict(color='lightseagreen')
    )

    # Create layout for the graphs
    scatter_layout = go.Layout(
        title=f'Scatter Plot: Min Time vs. Max Time - {technology.capitalize()}',
        xaxis=dict(title='Min Time'),
        yaxis=dict(title='Max Time'),
        margin=dict(l=50, r=50, b=100, t=100, pad=4),
        paper_bgcolor='rgba(0,0,0,0)',
        plot_bgcolor='rgba(0,0,0,0)'
    )

    violin_layout = go.Layout(
        title=f'Violin Plot: Distribution of Min Time - {technology.capitalize()}',
        yaxis=dict(title='Min Time'),
        margin=dict(l=50, r=50, b=100, t=100, pad=4),
        paper_bgcolor='rgba(0,0,0,0)',
        plot_bgcolor='rgba(0,0,0,0)'
    )

    # Create figures
    bar_fig = go.Figure(data=[bar_trace], layout=go.Layout(title='Bar Plot'))
    box_fig = go.Figure(data=[box_trace], layout=go.Layout(title='Box Plot'))
    scatter_fig = go.Figure(data=[scatter_trace], layout=scatter_layout)
    violin_fig = go.Figure(data=[violin_trace], layout=violin_layout)

    # Convert figures to HTML strings
    bar_div = bar_fig.to_html(full_html=False)
    box_div = box_fig.to_html(full_html=False)
    scatter_div = scatter_fig.to_html(full_html=False)
    violin_div = violin_fig.to_html(full_html=False)

    # Render the HTML template for individual report
    html_output = template.render(
        technology=technology.capitalize(),
        active_client=active_client,
        message_per_client=message_per_client,
        bar_div=bar_div,
        box_div=box_div,
        scatter_div=scatter_div,
        violin_div=violin_div,
        index_path='../index.html'
    )

    # Write individual report HTML to file
    with open(report_file, 'w', encoding='utf-8') as f:
        f.write(html_output)

    # Store report file path and details in reports dictionary
    report_key = f'{technology.capitalize()}-{active_client}-{message_per_client}'
    reports[report_key] = {
        'path': f'./{os.path.basename(report_file)}',
        'technology': technology.capitalize(),
        'active_client': active_client,
        'message_per_client': message_per_client
    }

# Generate index.html using index_template.html
index_template = env.get_template('index_template.html')
index_html_output = index_template.render(reports=reports)

# Write index.html to output directory
index_file = os.path.join(output_dir, 'index.html')
with open(index_file, 'w', encoding='utf-8') as f:
    f.write(index_html_output)

print(f"index.html generated successfully: {index_file}")
