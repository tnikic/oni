FROM fedora:40

# Install go and nodejs
RUN dnf install -y \
    git-2.45.2-2.fc40 \
    golang-1.22.5-1.fc40 \
    nodejs-20.12.2-1.fc40 \
    chromium-126.0.6478.182-1.fc40 \
    && dnf clean all

# Preparing folders
WORKDIR /app
COPY ./frontend ./frontend
COPY ./backend ./backend
COPY ./docker/start.sh .

# Install Go dependencies
WORKDIR /app/backend
RUN go mod tidy

# Install NodeJS dependencies
WORKDIR /app/frontend
RUN npm install

# Build the backend
WORKDIR /app/backend
RUN go build -o oni cmd/oni/main.go

# Build the frontend
WORKDIR /app/frontend
RUN npm run build

# Set environment variables
ENV ONI_HOST 0.0.0.0
ENV ONI_PORT 8080

# Expose port
EXPOSE 8080 3000

# Run the backend
WORKDIR /app
CMD ["./start.sh"]
