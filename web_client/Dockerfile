# Stage 1: Build the Next.js application
FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

# Stage 2: Serve the Next.js application
FROM node:18-alpine

# Set the working directory
WORKDIR /app

# Copy the built application from the previous stage
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/.env ./.env

# Install production dependencies
RUN npm install --only=production

# Expose port 3000
EXPOSE 3000

# Start the Next.js application
CMD ["npm", "start"]
