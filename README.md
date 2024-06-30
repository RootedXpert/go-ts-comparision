# go-ts-comparison

This project compares performance metrics between a Go client and a TypeScript client based on various message sizes.

## Overview

This repository contains Python scripts to analyze performance metrics from JSON files generated by the Go and TypeScript clients. It uses Plotly for data visualization and Jinja2 for HTML templating to create interactive reports.

## Report Link

You can view the generated reports at: [Go vs TypeScript Performance Comparison](https://go-vs-typescript.web.app)

## Project Structure

- `metric/`: Directory containing JSON files with performance metrics from different client implementations.
- `templates/`: Directory containing HTML templates for report generation.
- `dist/`: Output directory where generated HTML reports are stored.
- `index.html`: Index page listing links to all generated reports.

## Usage

1. **Setup Environment:**

   - Clone the repository: `git clone <repository_url>`
   - Install dependencies: `pip install -r requirements.txt`

2. **Generate Reports:**

   - Update `json_files` list in `generate_reports.py` with paths to your JSON files.
   - Run the script: `python generate_reports.py`
   - This will generate HTML reports in the `dist/` directory.

3. **View Reports:**
   - Open `dist/index.html` in a web browser to navigate through the generated reports.
   - Each report provides insights into the performance metrics of Go and TypeScript clients.

## Dependencies

- Python 3.x
- Pandas
- Plotly
- Jinja2

## Contributing

Contributions are welcome! Please fork the repository, create a new branch, make your changes, and submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
