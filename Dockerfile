# Use an official Golang runtime as a parent image
FROM golang:1.22 AS build

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the current directory contents into the container
COPY . .

# Build the Go app
RUN make build

# Start a new stage from scratch
FROM debian:bookworm-slim

# Install PostgreSQL client and dependencies for Chrome
RUN apt-get update && \
    apt-get install -y \
    wget \
    gnupg \
    ca-certificates \
    fonts-liberation \
    libappindicator3-1 \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcups2 \
    libdbus-1-3 \
    libgdk-pixbuf2.0-0 \
    libnspr4 \
    libnss3 \
    libx11-xcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    xdg-utils \
    libxss1 \
    libxtst6 \
    unzip \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Use ARG to define the architecture dynamically at build time
ARG TARGETARCH

# Install Google Chrome or Chromium based on architecture
RUN if [ "$TARGETARCH" = "amd64" ]; then \
        wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
        echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list && \
        apt-get update && apt-get install -y google-chrome-stable; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
        apt-get update && apt-get install -y chromium; \
    else \
        echo "Unsupported architecture"; \
        exit 1; \
    fi

# Set Chrome/Chromium as the headless browser for ChromeDP based on architecture
ENV CHROME_PATH=/usr/bin/google-chrome
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        export CHROME_PATH=/usr/bin/chromium; \
    fi

# Set the working directory in the container
WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=build /go/src/app/scraper .
COPY --from=build /go/src/app/api .
COPY --from=build /go/src/app/mailer .

COPY scripts/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh


# Run the Go app
CMD ["./scraper"]
