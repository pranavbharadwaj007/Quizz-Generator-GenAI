package main

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
	"time"
    "github.com/fatih/color"
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

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

const jsonFile = "quiz_data.json"

// Initialize color printers
var (
    green = color.New(color.FgGreen, color.Bold)
    red   = color.New(color.FgRed, color.Bold)
    cyan  = color.New(color.FgCyan)
    yellow = color.New(color.FgYellow)
)

func main() {
    // Setup cleanup on interrupt
    cleanup := make(chan os.Signal, 1)
    signal.Notify(cleanup, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-cleanup
        removeJSONFile()
        os.Exit(0)
    }()

    // Get user details
    reader := bufio.NewReader(os.Stdin)
    cyan.Print("Enter your name: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    cyan.Print("Enter the topic for quiz: ")
    topic, _ := reader.ReadString('\n')
    topic = strings.TrimSpace(topic)

    // Generate questions using Gemini
    questions := generateQuestions(topic)
    
    // Save questions to JSON file
    saveQuestionsToJSON(questions)

    // Start the quiz
    score := conductQuiz(name, topic, questions)

    // Display final results
    yellow.Printf("\nFinal Results for %s:\n", name)
    yellow.Printf("Topic: %s\n", topic)
    if score.Score > 0 {
        green.Printf("Score: %d\n", score.Score)
    } else {
        red.Printf("Score: %d\n", score.Score)
    }
    yellow.Printf("Total Attempts: %d\n", score.Attempts)

    // Cleanup
    removeJSONFile()
}

func generateQuestions(topic string) []QuizQuestion {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(""))
    if err != nil {
        log.Fatalf("Error creating client: %v", err)
    }
    defer client.Close()

    model := client.GenerativeModel("gemini-1.5-flash-8b")
    model.SetTemperature(1)
    model.SetTopK(40)
    model.SetTopP(0.95)
    model.SetMaxOutputTokens(8192)

    prompt := fmt.Sprintf("Generate 5 latest quiz questions in JSON form for the topic %s with 4 options in this format:\n{\n    \"question\": \"Question text\",\n    \"options\": [\"option1\", \"option2\", \"option3\", \"option4\"],\n    \"answer\": correct_option_index\n}", topic)

    resp, err := model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        log.Fatalf("Error generating content: %v", err)
    }

    var questions []QuizQuestion
    var genaiText genai.Text
    if len(resp.Candidates) > 0 {
        for _, part := range resp.Candidates[0].Content.Parts {
            switch p := part.(type) {
            case genai.Text:
                genaiText = p
                genaiText = genai.Text(strings.ReplaceAll(strings.ReplaceAll(string(genaiText), "```", ""), "json", ""))
                err := json.Unmarshal([]byte(genaiText), &questions)
                if err != nil {
                    log.Fatalf("Error unmarshaling JSON: %v", err)
                }
            default:
                log.Printf("Unhandled part type: %T", p)
            }
        }
    } else {
        log.Fatalf("No candidates found in response")
    }
    return questions
}

func saveQuestionsToJSON(questions []QuizQuestion) {
    existingQuestions := make([]QuizQuestion, 0)
    if _, err := os.Stat(jsonFile); err == nil {
        // File exists, read and append
        data, err := os.ReadFile(jsonFile)
        if err != nil {
            log.Fatalf("Error reading existing file: %v", err)
        }
        err = json.Unmarshal(data, &existingQuestions)
        if err != nil {
            log.Fatalf("Error unmarshaling existing JSON: %v", err)
        }
        // Append new questions
        existingQuestions = append(existingQuestions, questions...)
    } else {
        // File doesn't exist, just write the new questions
        existingQuestions = questions
    }

    file, err := json.MarshalIndent(existingQuestions, "", "    ")
    if err != nil {
        log.Fatalf("Error marshaling JSON: %v", err)
    }

    err = os.WriteFile(jsonFile, file, 0644)
    if err != nil {
        log.Fatalf("Error writing file: %v", err)
    }
}

func conductQuiz(name string, topic string, questions []QuizQuestion) UserScore {
    score := UserScore{
        Name:  name,
        Topic: topic,
    }

    for i, q := range questions {
        // Display question and options first
        yellow.Printf("\nQuestion %d: %s\n", i+1, q.Question)
        for j, opt := range q.Options {
            cyan.Printf("%d. %s\n", j+1, opt)
        }
        cyan.Print("\nEnter your answer (1-4): ")

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        answerCh := make(chan int)
        errorCh := make(chan error)

        go func() {
            reader := bufio.NewReader(os.Stdin)
            input, err := reader.ReadString('\n')
            if err != nil {
                errorCh <- err
                return
            }
            answer, err := strconv.Atoi(strings.TrimSpace(input))
            if err != nil {
                errorCh <- err
                return
            }
            answerCh <- answer
        }()

        // Handle timeout and answer
        select {
        case <-ctx.Done():
            red.Println("\nTime's up! Moving to next question")
            score.Attempts++
            score.Score--
        case err := <-errorCh:
            red.Printf("\nError reading input: %v\n", err)
        case answer := <-answerCh:
            handleAnswer(answer, q, &score)
        }
    }
    return score
}
func handleAnswer(answer int, q QuizQuestion, score *UserScore) {
    if answer < 1 || answer > 4 {
        red.Println("Invalid input. Skipping question.")
        return
    }

    score.Attempts++
    if answer-1 == q.Answer {
        green.Println("Correct! +4 points")
        score.Score += 4
    } else {
        red.Println("Wrong! -1 point")
        score.Score--
    }
}
func removeJSONFile() {
    if _, err := os.Stat(jsonFile); err == nil {
        os.Remove(jsonFile)
    }
}
