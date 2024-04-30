FROM ruby:3.3-slim

# Set the working directory within the container
WORKDIR /app

# Copy the entrypoint script into the container
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

# Make the entrypoint script executable
RUN chmod +x /usr/local/bin/entrypoint.sh

# Set the entrypoint for the container
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
