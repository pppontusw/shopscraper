# ShopScraper

ShopScraper is a web scraping tool that allows you to scrape product information from various online shops. It includes an API, a mailer service for notifications, and a scraper. It provides a flexible configuration system to define the scraping rules for each shop.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Steps](#steps)
  - [Next Steps](#next-steps)
- [Configuration](#configuration)
  - [Example YAML Configuration](#example-yaml-configuration)
  - [YAML Configuration Options](#yaml-configuration-options)
  - [Command Line Flags](#command-line-flags)
    - [Scraper](#scraper)
    - [Mailer](#mailer)
    - [API](#api)
  - [Environment Variables](#environment-variables)
- [Usage](#usage)
  - [Running with Docker](#running-with-docker)
  - [Building and Running Binaries](#building-and-running-binaries)
  - [Running During Development](#running-during-development)
  - [Running the Frontend](#running-the-frontend)
- [Deployment](#deployment)
  - [Frontend](#frontend)
  - [Backend](#backend)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Web Scraping**: Scrapes multiple web shops for product listings, including support for both simple HTML and JavaScript-driven websites.
- **Customizable Scraper Configurations**: Easily configurable for different shop layouts and pagination.
- **Automated Email Notifications**: Sends email notifications with newly found products, ensuring you're always up to date with the latest listings.
- **API**: Provides a RESTful API to access the scraped product data.
- **Scheduled Scraping Runs**: Configurable intervals for scraping operations, allowing for regular updates without manual intervention.
- **Docker Support**: Includes Docker and Docker Compose configurations for easy deployment and isolated environments.
- **Frontend Deployment**: Easily build and deploy the frontend application to AWS S3 and CloudFront.

## Installation

### Prerequisites

Before you begin the installation, ensure that you have the following prerequisites installed on your system:

- Go (version 1.22 or later)
- Node.js (version 21 or later)
- Docker and Docker Compose (for Docker deployment)
- AWS CLI (for deploying the frontend to AWS)

Make sure you have an internet connection for downloading the necessary packages and dependencies.

### Steps

1. **Clone the repository:**

   Start by cloning the ShopScraper repository to your local machine:

   ```shell
   git clone https://github.com/yourusername/ShopScraper.git
   cd ShopScraper
   ```

2. **Install Go dependencies:**

   Run the following command in the root directory of the project to install all required Go dependencies:

   ```shell
   go mod tidy
   ```

3. **Install Node.js dependencies:**

   Navigate to the frontend directory and install the Node.js dependencies:

   ```shell
   cd frontend
   npm install
   ```

4. **Build the Docker images (Optional):**

   If you prefer to run ShopScraper using Docker, build the Docker images:

   ```shell
   docker-compose build
   ```

You have now installed all the necessary components to run ShopScraper. Proceed to the [Configuration](#configuration) section to configure ShopScraper according to your needs.

### Next Steps

After installation, you can:

- Run ShopScraper in development mode by following the instructions in the [Running During Development](#running-during-development) section.
- Build and run the production binaries as described in the [Building and Running Binaries](#building-and-running-binaries) section.
- Deploy ShopScraper using Docker by following the steps in the [Running with Docker](#running-with-docker) section.
- Deploy the frontend application to AWS S3 and CloudFront by referring to the [Deployment](#deployment) section.

For more information on using and configuring ShopScraper, refer to the subsequent sections of this README.

## Configuration

ShopScraper uses a YAML configuration file to define the scraping rules for each shop. The configuration file is located at `config/config.yaml`.

### Example YAML Configuration

Here's an example configuration file:

```yaml
email:
  server: smtp.mailgun.org
  recipient: john@example.com
  sender: shopscraper@example.com
  subject: New items found
  port: 587

scrapers:
  - shopName: ExampleShop
    type: WebShopScraper
    urls:
      - https://example.com/products
    itemSelector: div.product
    nameSelector: h2.product-name
    priceSelector: 
      - span.price
    linkSelector: a.product-link
    nextPageSelector: a.next-page
    priceFormat: reverse
```

### YAML Configuration Options

- `email`: Email configuration for sending notifications.
  - `server`: SMTP server address.
  - `recipient`: Email address of the recipient.
  - `sender`: Email address of the sender.
  - `subject`: Subject of the email notification.
  - `port`: SMTP server port.
- `scrapers`: An array of scraper configurations.
  - `shopName`: Name of the shop being scraped.
  - `type`: Type of the scraper ("WebShopScraper" for regular web shops, "JavaScriptWebShopScraper" for JavaScript-rendered web shops).
  - `urls`: List of URLs to scrape.
  - `itemSelector`: CSS selector for identifying individual product items.
  - `nameSelector`: CSS selector for extracting the product name.
  - `priceSelector`: List of CSS selector(s) for extracting the product price.
  - `linkSelector`: CSS selector for extracting the product link.
  - `nextPageSelector`: (optional) CSS selector for identifying the next page link.
  - `priceFormat`: (optional) Format of the price string ("reverse" for prices in the format "1.499,00€", "double_eur" for prices in the format "1 499,00EUR 2 500,00EUR").
  - `retryString`: (optional) String to search for in the HTML content to determine if the page needs to be retried (used for JavaScript-rendered web shops), i.e. if this string is found the scraper will reload the page.

### Command Line Flags

#### Scraper

- `--daemon`: Enable daemon mode to run the scraper continuously at the interval specified in the yaml configuration.
- `--debug`: Enable debug mode to print additional information during scraping.
- `--config-path`: Specify the path to the configuration YAML file (default: `./config/config.yaml`).
- `--interval`: Interval between scraper runs. (only applicable in daemon mode)
- `--max-workers`: Maximum numbers of workers per scraper.
- `--keep-duration`: Duration of time to keep items in database (ex: 12h, 24h, 72h) (default: 72h)

#### Mailer

- `--daemon`: Enable daemon mode to run the mailer continuously at the interval specified in the yaml configuration.
- `--config-path`: Specify the path to the configuration YAML file (default: `./config/config.mailer.yaml`).
- `--interval`: Interval between mailer runs. Emails will still only be sent if there are new products to notify about. (only applicable in daemon mode)

#### API

- No command line flags for the API component.

### Environment Variables

The following environment variables are used by ShopScraper:

- `SHOPSCRAPER_DB_CONNECTION_STRING`: The connection string to your PostgreSQL database. This is required for all components that interact with the database (Scraper, Mailer, API).
  
  Example: `postgresql://user:password@localhost:5432/database_name`

- `SHOPSCRAPER_SMTP_PASSWORD`: The password for the SMTP server used by the Mailer component to send email notifications.
  
  Example: `your_smtp_password`

- `SHOPSCRAPER_API_KEY`: A custom API key for securing access to your API. This should be kept secret and used by clients to authenticate requests.

  Example: `your_secure_api_key`

- `S3_BUCKET_NAME`: The name of the S3 bucket used for deploying the frontend application. This is required by the `make deploy-frontend` command.

  Example: `your_s3_bucket_name`

- `DISTRIBUTION_ID`: The ID of the CloudFront distribution used for serving the frontend application. This is required by the `make deploy-frontend` command.

  Example: `your_cloudfront_distribution_id`

- `REACT_APP_API_URL`: The URL of the API endpoint. This is used to configure the frontend application to communicate with your API.

  Example: `https://api.yourdomain.com`

## Usage

ShopScraper provides multiple ways to run the scraper, mailer, and API components.

### Running with Docker

First make sure you have the configuration in place in `config/config.yaml`.

1. Build the Docker image:

   ```shell
   make build-docker
   ```

2. Start the services using Docker Compose:

   ```shell
   docker-compose up
   ```

   This will start the scraper, mailer, API, and PostgreSQL database containers.

### Building and Running Binaries

1. Build the binaries:

   ```shell
   make build
   ```

   This will generate the `scraper`, `mailer`, and `api` binaries.

2. Run the scraper:

   ```shell
   export SHOPSCRAPER_DB_CONNECTION_STRING="postgresql://YOUR-PRODUCTION-POSTGRES-CONNECTION-STRING"
   ./scraper --config-path ./config/config.yaml --daemon
   ```

3. Run the mailer:

   ```shell
    export SHOPSCRAPER_SMTP_PASSWORD="YOUR-SMTP-PASSWORD"
    export SHOPSCRAPER_DB_CONNECTION_STRING="postgresql://YOUR-PRODUCTION-POSTGRES-CONNECTION-STRING"
   ./mailer --config-path ./config/config.yaml --daemon
   ```

4. Run the API:

   ```shell
   export SHOPSCRAPER_DB_CONNECTION_STRING="postgresql://YOUR-PRODUCTION-POSTGRES-CONNECTION-STRING"
   export SHOPSCRAPER_API_KEY="YOUR-PRODUCTION-API-KEY"
   ./api
   ```

### Running During Development

1. Start the PostgreSQL test database:

   ```shell
   make start-dependencies
   ```

2. Run the scraper:

   ```shell
   export SHOPSCRAPER_DB_CONNECTION_STRING="postgresql://test:Test1234@localhost:5432/test?sslmode=disable"
   make run-scraper
   ```

   This will run the scraper (once) using `go run cmd/scraper/main.go`.

3. Run the mailer:

   ```shell
    export SHOPSCRAPER_SMTP_PASSWORD="YOUR-SMTP-PASSWORD"
    export SHOPSCRAPER_DB_CONNECTION_STRING="postgresql://test:Test1234@localhost:5432/test?sslmode=disable"
   make run-mailer
   ```

   This will run the mailer (once) using `go run cmd/mailer/main.go`.

4. Run the API:

   ```shell
   export SHOPSCRAPER_DB_CONNECTION_STRING="postgresql://test:Test1234@localhost:5432/test?sslmode=disable"
   export SHOPSCRAPER_API_KEY="test1234"
   make run-api
   ```

   This will run the API using `go run cmd/api/main.go`.

### Running the Frontend

For development:

1. Install the frontend dependencies:

   ```shell
   cd frontend
   npm install
   ```

2. Start the frontend development server:

   ```shell
   npm start
   ```

   This will start the frontend development server and open the application in your default browser, use the API key you specified when starting the API server.

## Deployment

### Frontend

To deploy the frontend to AWS S3 and CloudFront:

Create the necessary resources in AWS:

- S3 bucket
- ACM Certificate (to enable https on cloudfront)
- Cloudfront Distribution

Once all the resources are in place, make sure that you have AWS CLI configured and set the following environment variables:

- `S3_BUCKET_NAME`: The name of the S3 bucket to deploy to.
- `DISTRIBUTION_ID`: The ID of the CloudFront distribution.
- `REACT_APP_API_URL`: The URL of the API.

Then run: `make deploy-frontend`

### Backend

The backend components (scrapers, API, mailer) can be easily deployed using either Docker Compose or by building and running the binaries. A `docker-compose.example.yml` example is provided that you can adjust according to your deployment requirements.

## Cleaning Up

To clean up resources created during the build or runtime, you can use:

- `make clean`: Cleans up local build artifacts.
- `make clean-all`: Cleans up all local artifacts including node modules in the frontend.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## Releasing New Versions

To tag a new release based on the current Git tags:

```bash
make release
```

This command increments the patch version, creates a new Git tag, and pushes it to the repository.

## License

This project is licensed under the [MIT License](LICENSE).
