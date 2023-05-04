# Munich Animal Shelter Adoption Data Analysis
This is a Go project that uses the Gin web framework and Colly web scraping library to crawl the Munich Animal Shelter home page and gather adoption data of dogs. The collected data is then analyzed to gain insights about the adoption trends and patterns.

**Mind that this project is still in development and currently not working.**

## Requirements
To run this project, you need the following dependencies:

* Docker v20.10 or higher

## Installation

1. Clone this repository to your local machine.
2. Make sure you have Docker installed and set up on your machine.
3. Build the Docker image by running the following command:
```bash
docker build -t munich-animal-shelter-adoption-data-analysis .
```

## Usage

To use this project, follow these steps:

1. Start the Docker container by running the following command:
```bash
docker run -p 8080:8080 munich-animal-shelter-adoption-data-analysis
```

2. Open your web browser and go to [http://localhost:8080](http://localhost:8080).
3. The server will crawl the Munich Animal Shelter home page and collect adoption data of dogs. 
Once the data is collected, it will be displayed on the web page (when it is implemented).
You can use the data to analyze the adoption trends and patterns.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
Thank you to the Munich Animal Shelter for their compassionate work.