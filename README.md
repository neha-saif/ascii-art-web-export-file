# ascii-art-web

## Overview

The ASCII Art Generator is a simple web application that converts multiline text input into ASCII art. Users can enter text, choose an art style, and view the generated ASCII art in their browser. The application is built with Go and provides a straightforward interface for creating and displaying ASCII art.

## Features

- **Text Input**: Enter multiline text to be converted into ASCII art.
- **Art Styles**: Choose from different ASCII art styles (e.g., Standard, Thinkertoy, Shadow).
- **Scrollable Output**: View the generated ASCII art with horizontal scrolling if the content exceeds the container width.

## Running the Project

1. **Clone the Repository**:

   ```sh
   git clone <repository-url>
   cd <repository-directory>
   ```

2. **Run the Application:**:

   ```sh
   go run .
   ```

3. **Access the Application:**:
   Open your web browser and navigate to `http://localhost:8080` to use the ASCII Art Generator.

## How It Works

- The application listens for HTTP requests on port 8080.
- Users submit their text and select an art style via a web form.
- The server processes the input, converts it to ASCII art using the specified style, and returns the result.
- The ASCII art is displayed in a scrollable container if it's too wide for the viewport.

## Dependencies

- Go (1.18+ recommended)

## Credits

- **Animated Background**: The animated background is inspired by [Manuel Pinto's CodePen](https://codepen.io/P1N2O/pen/pyBNzX).
