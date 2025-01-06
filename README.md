# QuizGenerator GenAI

An interactive CLI quiz application powered by Google's Gemini AI that generates dynamic questions based on any topic.

## Features

- Dynamic quiz generation using Gemini AI
- Topic-based customized questions
- Interactive command-line interface
- Colorful output for better user experience
- Timed questions (30-second limit per question)
- Score tracking with points system
- Automatic cleanup of temporary files

## Prerequisites

- Go 1.19 or higher
- Google Cloud API key with Gemini AI access

## Dependencies

```go
import (
    "github.com/fatih/color"          // Terminal color output
    "github.com/google/generative-ai-go/genai"  // Gemini AI client
    "google.golang.org/api/option"    // Google API options
)
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/quizgenerator-genai.git
cd quizgenerator-genai
```

2. Install dependencies:
```bash
go mod init quizgenerator
go get github.com/fatih/color
go get github.com/google/generative-ai-go/genai
go get google.golang.org/api/option
```

3. Set up your Google Cloud API key:
```bash
export GOOGLE_API_KEY="your-api-key-here"
```

## Usage

1. Run the application:
```bash
go run main.go
```

2. Follow the prompts:
   - Enter your name
   - Choose a quiz topic
   - Answer generated questions within the time limit

## Scoring System

- Correct answer: +4 points
- Wrong answer: -1 point
- Time expired: -1 point

## Technical Details

## Architecture Flow



## Project Structure

```
quizgenerator-genai/
├── main.go                # Main application code
├── quiz_data.json         # Temporary storage for quiz questions
└── README.md             # Project documentation
```

## Golang Implementation Details

### Data Structures
```go
type QuizQuestion struct {
    Question string   `json:"question"`
    Options  []string `json:"options"`
    Answer   int      `json:"answer"`
}

type UserScore struct {
    Name     string `json:"name"`
    Topic    string `json:"topic"`
    Score    int    `json:"score"`
    Attempts int    `json:"attempts"`
}
```

### Key Components

1. **Gemini AI Integration**
   - Uses `genai` package for question generation
   - Configurable model parameters: Temperature, TopK, TopP
   - Custom prompt engineering for consistent JSON output

2. **Concurrency Handling**
   - Context-based timeout management
   - Graceful shutdown with signal handling
   - Goroutines for async user input

3. **Error Management**
   ```go
   // Example error handling pattern
   if err := operation(); err != nil {
       log.Fatalf("Operation failed: %v", err)
   }
   ```

4. **File Operations**
   - JSON-based question storage
   - Atomic file operations
   - Automatic cleanup on exit

5. **User Interface**
   - Color-coded output using `fatih/color`
   - Interactive CLI prompts
   - Real-time feedback

## Functions

- `generateQuestions()`: Creates quiz questions using Gemini AI
- `conductQuiz()`: Manages quiz flow and user interaction
- `handleAnswer()`: Processes user answers and updates score
- `saveQuestionsToJSON()`: Handles temporary data storage

## Error Handling

- Input validation for user responses
- Timeout handling for questions
- API error management
- File operation error handling

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
