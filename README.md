# Affirmation Generator

This project is a web-based Affirmation Generator that provides users with unique, uplifting affirmations each time they press a button. The project includes a frontend interface and a backend service that leverages the OpenAI API to generate affirmations. Both components are containerized and designed to run on a Kubernetes cluster.

## Features

- **Rate Limiting:** Limits each IP to 5 requests per hour.
- **Affirmation Variety:** Each request prompts a unique affirmation selected randomly from several positive themes.
- **OpenAI Integration:** The backend uses the OpenAI API to create personalized affirmations.
- **Frontend UI:** A simple button-driven user interface with a response display for the generated affirmations.
- **Kubernetes-Ready:** Both frontend and backend are designed to run in a Kubernetes environment with FluxCD.

## Project Structure

```
- /                         # Backend Golang server using Gin framework
- static/                   # Frontend HTML, CSS, and JS files
- kubernetes/               # Kubernetes FluxCD helmrelease and network policy
```

## Blog

This repository is used to store the code from my blog post: [Running a Button on Kubernetes](https://khenry.substack.com/p/button-on-kubernetes).