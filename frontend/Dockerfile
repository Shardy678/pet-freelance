# Stage 1: Build the project
FROM node:18 AS build

# Set working directory in the container
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package.json and pnpm-lock.yaml files first to leverage Docker cache
COPY package.json pnpm-lock.yaml ./

# Install dependencies using pnpm
RUN pnpm install

# Copy the rest of the application code
COPY . .

# Build the Vite app
RUN pnpm run build

# Stage 2: Serve the project
FROM nginx:alpine

# Copy the built project from the build stage to the Nginx directory
COPY --from=build /app/dist /usr/share/nginx/html

# Expose the port that Nginx will run on
EXPOSE 80

# Command to run Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
